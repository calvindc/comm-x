package orm_db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var db *gorm.DB

/*
ORM 允许通过一个现有的数据库连接来初始化 *gorm.DB
import (

	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

)

	sqlDB, err := sql.Open("mysql", "mydb_dsn")
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
	  Conn: sqlDB,
	}), &gorm.Config{})

import (

	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

)

	sqlDB, err := sql.Open("pgx", "mydb_dsn")
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
	  Conn: sqlDB,
	}), &gorm.Config{})
*/
func SetupDB(dbtype, dsn string, firstcreate interface{}, inittables ...interface{}) (err error) {
	if dsn == "" {
		return fmt.Errorf("you should given a the driver's datasource to connect")
	}
	switch dbtype {
	case "mysql":
		dsn = "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		//db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{})
	case "postgres":
		//dsn = "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		//db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		//使用 pgx 作为 postgres 的 database/sql 驱动，默认情况下，它会启用 prepared statement 缓存，可以这样禁用它：
		dsn = "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})

	case "sqlite3":
		//dsn = "your/path/database"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "sqlserver":
		dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
		db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	}

	//database/sql维护连接池
	sqlDB, err := db.DB()
	if err != nil {
		err = fmt.Errorf("failed to connect database, got error %v", err)
		return
	}
	sqlDB.SetMaxIdleConns(10)           //空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)          //数据库的最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //连接可以重用的最大时间量
	sqlDB.SetConnMaxIdleTime(0)         //连接可能空闲的最长时间

	//loger
	db.Logger.LogMode(logger.Warn)

	//migrate table(create table first)
	for _, table := range inittables {
		err = db.AutoMigrate(&table)
		if err != nil {
			return
		}
	}

	//some data create first
	db.FirstOrCreate(firstcreate)

	return
}

func CloseDB() error {
	sqldb, err := db.DB()
	if err != nil {
		return fmt.Errorf("cannot get a impl for db-conn, close err=%v", err)
	}
	if err = sqldb.Close(); err != nil {
		return fmt.Errorf("close db err=%v", err)
	}
	return nil
}

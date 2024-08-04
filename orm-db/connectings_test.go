package orm_db

import (
	"fmt"
	"os"
	"path"
	"testing"
)

type Product struct {
	ID    uint
	Code  string
	Price uint
	Sum   string
}

type Apple struct {
	ID    uint
	Code  string
	Price uint
	Sole  bool
}
type Aa struct {
	A1 int
	A2 int
}

func TestSetupDB(t *testing.T) {
	dbSource := path.Join(os.TempDir(), fmt.Sprintf("test%s.db", ""))
	t.Log(dbSource)
	err := os.Remove(dbSource)
	if err != nil {
		t.Logf("remove err %s", err)
	}
	dsn := dbSource
	db, err := SetupDB("sqlite3", dsn, &Product{}, &Apple{}, &Aa{})
	if err != nil {
		t.Fatal(err)
	}
	var latestid = &Product{Code: "code01", Price: 11, Sum: "123000000000001"}
	db.Create(latestid)

	read := &Product{}
	db.Last(&read, "code = ?", "code01")
	if read.Price != 11 {
		t.Error("not passed")
	}

	aa := &Aa{111, 222}

	FirstOrCreate(db, aa)

	err = CloseDB(db)
	if err != nil {
		t.Error(err)
	}
}

func TestSetupDBPostgres(t *testing.T) {
	dsn := "host=16.163.154.147 user=dbuser1 password=123456 dbname=testdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := SetupDB("postgres", dsn, &Product{}, &Apple{}, &Aa{})
	if err != nil {
		t.Fatal(err)
	}
	var latestid = &Product{Code: "code01", Price: 10, Sum: "www"}
	db.Create(latestid)

	read := &Product{}
	db.First(&read, "code = ?", "code01")
	if read.Price != 10 {
		t.Error("not passed")
	}

	err = CloseDB(db)
	if err != nil {
		t.Error(err)
	}
}

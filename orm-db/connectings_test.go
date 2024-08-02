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
}

func TestSetupDB(t *testing.T) {
	dbSource := path.Join(os.TempDir(), fmt.Sprintf("test%s.db", ""))
	t.Log(dbSource)
	err := os.Remove(dbSource)
	if err != nil {
		t.Logf("remove err %s", err)
	}
	dsn := dbSource
	err = SetupDB("sqlite3", dsn, nil, &Product{})
	if err != nil {
		t.Fatal(err)
	}
	var latestid = &Product{Code: "code01", Price: 10}
	db.Create(latestid)

	read := &Product{}
	db.Last(&read, "code = ?", "code01")
	if read.Price != 10 {
		t.Error("not passed")
	}

	err = CloseDB()
	if err != nil {
		t.Error(err)
	}
}

func TestCloseDB(t *testing.T) {
	err := CloseDB()
	if err != nil {
		t.Error(err)
	}
}

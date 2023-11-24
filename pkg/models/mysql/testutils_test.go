package mysql

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"testing"
)

func newTestDB(t *testing.T) (*sql.DB, func()) {
	dsnTest := flag.String("dsn", "root:password@/test_snippetbox?parseTime=true", "MYSQL Database Connection string")
	db, err := sql.Open("mysql", *dsnTest)
	if err != nil {
		t.Fatal(err)
	}
	script, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(script))
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
	return db, func() {
		script, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	}
}

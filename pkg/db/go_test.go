package db

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func Test_Tables(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error creat db and mock: '%v'", err)
	}
	defer db.Close()

	instAct, err := RepoDB(db)
	if err != nil {
		t.Fatalf("an error creat actions instance: '%v'", err)
	}

	// Ожидаем вызов Exec на INSERT
	mock.ExpectExec(regexp.QuoteMeta(`
	CREATE TABLE IF NOT EXISTS main (
	id INTEGER PRIMARY KEY,
	nameTableI string UNIQUE,
	nameTableW string UNIQUE,
	nameTableE string UNIQUE,
	timestamp TEXT DEFAULT CURRENT_TIMESTAMP);
	`))

	_ = instAct

}

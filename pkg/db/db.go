package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

type MessageT struct {
	TypeMessage   string // type of message - I, W, E
	NameProject   string
	LocationEvent string
	BodyMessage   string
}

type ObjectDB struct {
	DB *sql.DB
}

type ActionsDB interface {
	Tables() error
	StoreMessage(msg MessageT) error
}

// Create the db object. Return interface.
func RepoDB(db *sql.DB) ActionsDB {
	return &ObjectDB{DB: db}
}

// Create tables
func (o ObjectDB) Tables() error {

	// main table
	_, err := o.DB.Exec(`
	CREATE TABLE IF NOT EXISTS main (
	id INTEGER PRIMARY KEY,
	nameTableI string UNIQUE,
	nameTableW string UNIQUE,
	nameTableE string UNIQUE,
	timestamp TEXT DEFAULT CURRENT_TIMESTAMP);
	`)
	if err != nil {
		return fmt.Errorf("main table is not created: %v", err)
	}

	// log tables
	row := o.DB.QueryRow("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1")

	var nameI, nameW, nameE string

	err = row.Scan(&nameI, &nameW, &nameE)
	if err != nil && errors.Is(err, sql.ErrNoRows) {

		nameI = "logI_1"
		nameW = "logW_1"
		nameE = "logE_1"
		_, err := o.DB.Exec("INSERT INTO main (nameTableI, nameTableW, nameTableE) VALUES (?, ?, ?)", nameI, nameW, nameE)
		if err != nil {
			return fmt.Errorf("fault initialisation the main table: %v", err)
		}
	}

	err = tableCheckCreate(o.DB, nameI)
	if err != nil {
		return fmt.Errorf("fault: {%v}", err)
	}
	err = tableCheckCreate(o.DB, nameW)
	if err != nil {
		return fmt.Errorf("fault: {%v}", err)
	}
	err = tableCheckCreate(o.DB, nameE)
	if err != nil {
		return fmt.Errorf("fault: {%v}", err)
	}

	return nil
}

// Saving the received message in the database. Return error
func (o ObjectDB) StoreMessage(msg MessageT) error {

	row := o.DB.QueryRow("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1")

	var nameI, nameW, nameE string

	err := row.Scan(&nameI, &nameW, &nameE)
	if err != nil {
		return fmt.Errorf("store a information -> flt read main: %v", err)
	}

	// Saving
	switch msg.TypeMessage {
	case "I":
		q := fmt.Sprintf("INSERT INTO %s (nameProject, locationEvent, bodyMessage) VALUES (:project, :location, :body)", nameI)
		_, err := o.DB.Exec(q,
			sql.Named("project", msg.NameProject),
			sql.Named("location", msg.LocationEvent),
			sql.Named("body", msg.BodyMessage))
		if err != nil {
			return fmt.Errorf("store a information -> flt store I: %v", err)
		}
	case "W":
		q := fmt.Sprintf("INSERT INTO %s (nameProject, locationEvent, bodyMessage) VALUES (:project, :location, :body)", nameW)
		_, err := o.DB.Exec(q,
			sql.Named("project", msg.NameProject),
			sql.Named("location", msg.LocationEvent),
			sql.Named("body", msg.BodyMessage))
		if err != nil {
			return fmt.Errorf("store a information -> flt store W: %v", err)
		}
	case "E":
		q := fmt.Sprintf("INSERT INTO %s (nameProject, locationEvent, bodyMessage) VALUES (:project, :location, :body)", nameE)
		_, err := o.DB.Exec(q,
			sql.Named("project", msg.NameProject),
			sql.Named("location", msg.LocationEvent),
			sql.Named("body", msg.BodyMessage))
		if err != nil {
			return fmt.Errorf("store a information -> flt store E: %v", err)
		}
	default:
		return errors.New("not allowed type of message")

	}

	return nil
}

// Connect DB
func ConDb(typeDB, nameDB string) (*sql.DB, func() error, error) {

	db, err := sql.Open("sqlite", "iwe.db")
	if err != nil {
		return nil, nil, fmt.Errorf("error connect DB: %v", err)
	}

	closeDB := func() error {
		err := db.Close()
		if err != nil {
			return fmt.Errorf("fault close connect DB: %v", err)
		}
		return nil
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, fmt.Errorf("fault ping DB: %v", err)
	}

	return db, closeDB, nil
}

// Check create table by name
func tableCheckCreate(db *sql.DB, name string) error {

	if db == nil {
		return fmt.Errorf("fault check create table {%s} -> not pointer db", name)
	}
	if name == "" {
		return errors.New("fault check create table -> no name table")
	}

	q := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
	id INTEGER PRIMARY KEY,
	nameProject string NOT NULL,
	locationEvent string NOT NULL,
	bodyMessage string NOT NULL,
	timestamp TEXT DEFAULT CURRENT_TIMESTAMP);
	`, name)

	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("table {%s} is not created: %v", name, err)
	}

	return nil
}

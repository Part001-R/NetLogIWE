package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type MessageT struct {
	TypeMessage   string // type of message - I, W, E, T(test connect)
	NameProject   string
	LocationEvent string
	BodyMessage   string
}

type ObjectDB struct {
	DB *sql.DB
}

type ActionsDB interface {
	Tables() error
	SavingMessage(msg MessageT) error
}

// Create the db object. Return interface.
func RepoDB(db *sql.DB) (ActionsDB, error) {
	if db == nil {
		return nil, errors.New("empty pinter db")
	}
	return &ObjectDB{DB: db}, nil
}

// Working with database tables
func (o ObjectDB) Tables() error {

	err := checkCreateMainTable(o.DB)
	if err != nil {
		log.Fatal(err)
	}

	nI, nW, nE, err := readLogTablesName(o.DB)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		nI = "logI_1"
		nW = "logW_1"
		nE = "logE_1"
		_, err := initMainTable(o.DB, nI, nW, nE)
		if err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	err = checkCreateLogTable(o.DB, nI)
	if err != nil {
		log.Fatal(err)
	}

	err = checkCreateLogTable(o.DB, nW)
	if err != nil {
		log.Fatal(err)
	}

	err = checkCreateLogTable(o.DB, nE)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// Saving the received message in the database. Return error
func (o ObjectDB) SavingMessage(msg MessageT) error {

	nameI, nameW, nameE, err := readLogTablesName(o.DB)
	if err != nil {
		return fmt.Errorf("fault read name of tables: {%v}", err)
	}

	// Saving the message
	switch msg.TypeMessage {
	case "I":
		err := savingMessageCheckResult(o.DB, nameI, msg)
		if err != nil {
			return fmt.Errorf("fault save I: {%v}", err)
		}
	case "W":
		err := savingMessageCheckResult(o.DB, nameW, msg)
		if err != nil {
			return fmt.Errorf("fault save W: {%v}", err)
		}
	case "E":
		err := savingMessageCheckResult(o.DB, nameE, msg)
		if err != nil {
			return fmt.Errorf("fault save E: {%v}", err)
		}
	default:
		return errors.New("not allowed type of message when saving")
	}

	return nil
}

// Check overload the log table
func checkOverloadLogTable(typeTable string, lastId int64) (bool, error) {

	// read env constant
	var maxId int64
	switch typeTable {
	case "I":
		valENV := os.Getenv("MAX_IDNUMB_LOGI")
		maxId_t, err := strconv.ParseInt(valENV, 10, 64)
		if err != nil {
			return false, fmt.Errorf("fault parse string I: {%v}", err)
		}
		maxId = maxId_t
	case "W":
		valENV := os.Getenv("MAX_IDNUMB_LOGW")
		maxId_t, err := strconv.ParseInt(valENV, 10, 64)
		if err != nil {
			return false, fmt.Errorf("fault parse string I: {%v}", err)
		}
		maxId = maxId_t
	case "E":
		valENV := os.Getenv("MAX_IDNUMB_LOGE")
		maxId_t, err := strconv.ParseInt(valENV, 10, 64)
		if err != nil {
			return false, fmt.Errorf("fault parse string I: {%v}", err)
		}
		maxId = maxId_t
	default:
		return false, fmt.Errorf("not supported type of table: {%s}", typeTable)
	}

	if lastId > maxId {
		return true, nil
	}

	return false, nil
}

// Saving rx message
func doSaving(db *sql.DB, tableName string, msg MessageT) (int64, error) {

	q := fmt.Sprintf("INSERT INTO %s (nameProject, locationEvent, bodyMessage) VALUES (:project, :location, :body)", tableName)

	result, err := db.Exec(q,
		sql.Named("project", msg.NameProject),
		sql.Named("location", msg.LocationEvent),
		sql.Named("body", msg.BodyMessage))
	if err != nil {
		return 0, fmt.Errorf("store an information -> flt store %s message: %v", msg.TypeMessage, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fault get information about id %s table", msg.TypeMessage)
	}

	return id, nil
}

// Save message + check overload log table + update name log table + create new table
func savingMessageCheckResult(db *sql.DB, nameTable string, msg MessageT) error {

	id, err := doSaving(db, nameTable, msg)
	if err != nil {
		return fmt.Errorf("fault saving {%s} message: {%v}", msg.TypeMessage, err)
	}

	over, err := checkOverloadLogTable(msg.TypeMessage, id)
	if err != nil {
		return fmt.Errorf("fault check overload {%s} table: {%v}", msg.TypeMessage, err)
	}

	if over {
		err := changeLogTableNameCreate(db, msg.TypeMessage)
		if err != nil {
			return fmt.Errorf("fault update name of {%s} table: {%v}", msg.TypeMessage, err)
		}
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
func checkCreateLogTable(db *sql.DB, name string) error {
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

// Initialisation the main table
func initMainTable(db *sql.DB, nameI, nameW, nameE string) (int64, error) {
	if db == nil {
		return 0, errors.New("missed db pointer")
	}
	if nameI == "" {
		return 0, errors.New("empty content of nameI")
	}
	if nameW == "" {
		return 0, errors.New("empty content of nameW")
	}
	if nameE == "" {
		return 0, errors.New("empty content of nameE")
	}

	result, err := db.Exec("INSERT INTO main (nameTableI, nameTableW, nameTableE) VALUES (?, ?, ?)", nameI, nameW, nameE)
	if err != nil {
		return 0, fmt.Errorf("fault initialisation the main table: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("fault get id by LastInsertId: {%v}", err)
	}

	if id != 1 {
		return 0, fmt.Errorf("the id must be 1, but current: {%d} ", id)
	}

	return id, nil
}

// Reading log table names from the main table
func readLogTablesName(db *sql.DB) (nameI, nameW, nameE string, err error) {
	if db == nil {
		return "", "", "", errors.New("missed db pointer")
	}

	row := db.QueryRow("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1")

	err = row.Scan(&nameI, &nameW, &nameE)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return "", "", "", sql.ErrNoRows
	}
	if err != nil {
		return "", "", "", fmt.Errorf("fault reading log table names from the main table: {%v}", err)
	}

	return nameI, nameW, nameE, nil
}

// Create tables
func checkCreateMainTable(db *sql.DB) error {
	if db == nil {
		return errors.New("missed db pointer")
	}

	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS main (
	id INTEGER PRIMARY KEY,
	nameTableI string UNIQUE,
	nameTableW string UNIQUE,
	nameTableE string UNIQUE,
	timestamp TEXT DEFAULT CURRENT_TIMESTAMP);
	`)
	if err != nil {
		return fmt.Errorf("fault create the main table: %v", err)
	}

	return nil
}

// Change the name of log table
func changeLogTableNameCreate(db *sql.DB, typeTable string) error {
	if db == nil {
		return errors.New("missed db pointer")
	}

	nameI, nameW, nameE, err := readLogTablesName(db)
	if err != nil {
		return fmt.Errorf("fault read names of tables")
	}
	var newName string

	switch typeTable {
	case "I":
		newName, err = incrementIdInName(nameI)
		if err != nil {
			return fmt.Errorf("fault change name of I table: {%v}", err)
		}
		err = updateNameLogTable(db, newName, typeTable)
		if err != nil {
			return fmt.Errorf("fault update the name of I table: {%v}", err)
		}

	case "W":
		newName, err = incrementIdInName(nameW)
		if err != nil {
			return fmt.Errorf("fault change name of W table: {%v}", err)
		}
		err = updateNameLogTable(db, newName, typeTable)
		if err != nil {
			return fmt.Errorf("fault update the name of W table: {%v}", err)
		}

	case "E":
		newName, err = incrementIdInName(nameE)
		if err != nil {
			return fmt.Errorf("fault change name of E table: {%v}", err)
		}
		err = updateNameLogTable(db, newName, typeTable)
		if err != nil {
			return fmt.Errorf("fault update the name of E table: {%v}", err)
		}

	default:
		return fmt.Errorf("error in type of table. want I or W or E, recieve: {%s}", typeTable)
	}

	// Create new table
	err = checkCreateLogTable(db, newName)
	if err != nil {
		return fmt.Errorf("fault create new table: {%s}", newName)
	}

	return nil
}

// Increment an index in the name log table
func incrementIdInName(name string) (string, error) {

	sl := strings.Split(name, "_")
	if len(sl) != 2 {
		return "", fmt.Errorf("not correct the name of table: {%s}", name)
	}

	index, err := strconv.Atoi(sl[1])
	if err != nil {
		return "", fmt.Errorf("name table not have the index table: {%s}", name)
	}

	index += 1

	return fmt.Sprintf("%s_%d", sl[0], index), nil
}

// Update the name log table in the main table
func updateNameLogTable(db *sql.DB, newName, typeTable string) error {
	if db == nil {
		return errors.New("missed db pointer")
	}

	switch typeTable {
	case "I":
		_, err := db.Exec("UPDATE main SET nameTableI=? WHERE id=1", newName)
		if err != nil {
			return fmt.Errorf("fault update the name of I table: {%v}", err)
		}

	case "W":
		_, err := db.Exec("UPDATE main SET nameTableW=? WHERE id=1", newName)
		if err != nil {
			return fmt.Errorf("fault update the name of W table: {%v}", err)
		}

	case "E":
		_, err := db.Exec("UPDATE main SET nameTableE=? WHERE id=1", newName)
		if err != nil {
			return fmt.Errorf("fault update the name of E table: {%v}", err)
		}

	default:
		return fmt.Errorf("error in type of table. want I or W or E, recieve: {%s}", typeTable)
	}

	return nil
}

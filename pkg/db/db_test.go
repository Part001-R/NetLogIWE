package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =======================
// ==       PUBLIC      ==
// =======================

// =======================
// ==      SUCCESS      ==
// =======================

// Test - Saving the received message in the database. Return error
func Test_SavingMessage_SUCCESS(t *testing.T) {

	msg := []MessageT{
		{
			TypeMessage:   "I",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:65",
			BodyMessage:   "Not equal",
		},
		{
			TypeMessage:   "W",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:66",
			BodyMessage:   "Not equal",
		},
		{
			TypeMessage:   "E",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:67",
			BodyMessage:   "Not equal",
		},
	}

	tests := []struct {
		nameTest string
		initMock func(mock sqlmock.Sqlmock)
		index    int
	}{
		{
			nameTest: "Save I msg. Not Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			index: 0,
		},
		{
			nameTest: "Save I msg. Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(20, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logI_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))

			},
			index: 0,
		},
		{
			nameTest: "Save W msg. Not Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			index: 1,
		},
		{
			nameTest: "Save W msg. Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(20, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logW_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))

			},
			index: 1,
		},
		{
			nameTest: "Save E msg. Not Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			index: 2,
		},
		{
			nameTest: "Save E msg. Over ",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(20, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logE_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))

			},
			index: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.initMock(mock)

			instAct, err := RepoDB(db)
			require.NoError(t, err)

			instAct.SavingMessage(msg[tt.index])
			require.NoError(t, err)

		})
	}
}

// Test - Working with database tables
func Test_Tables_SUCCESS(t *testing.T) {

	tests := []struct {
		nameTest string
		initMock func(mock sqlmock.Sqlmock)
	}{
		{
			nameTest: "Tables is missed",
			initMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnError(sql.ErrNoRows)

				mock.ExpectExec("INSERT INTO main").
					WithArgs("logI_1", "logW_1", "logE_1").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
		{
			nameTest: "Tables is exists",
			initMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.initMock(mock)

			instAct, err := RepoDB(db)
			require.NoError(t, err)

			err = instAct.Tables()
			require.NoError(t, err)

		})
	}
}

// =======================
// ==       FAULT       ==
// =======================

// Test - Saving the received message in the database. Return error
func Test_SavingMessage_FAULT(t *testing.T) {

}

// =======================
// ==      INTERNAL     ==
// =======================

// =======================
// ==      SUCCESS      ==
// =======================

// Test - Check overload the log table
func Test_checkOverloadLogTable_SUCCESS(t *testing.T) {

	tests := []struct {
		nameTest   string
		typeTable  string
		curStrNumb int64
		maxI       string
		maxW       string
		maxE       string
		wantFlag   bool
	}{
		{
			nameTest:   "Not over I",
			typeTable:  "I",
			curStrNumb: 1,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   false,
		},
		{
			nameTest:   "Not over W",
			typeTable:  "W",
			curStrNumb: 1,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   false,
		},
		{
			nameTest:   "Not over E",
			typeTable:  "E",
			curStrNumb: 1,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   false,
		},
		{
			nameTest:   "Over I",
			typeTable:  "I",
			curStrNumb: 11,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   true,
		},
		{
			nameTest:   "Over W",
			typeTable:  "W",
			curStrNumb: 11,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   true,
		},
		{
			nameTest:   "Over E",
			typeTable:  "E",
			curStrNumb: 11,
			maxI:       "10",
			maxW:       "10",
			maxE:       "10",
			wantFlag:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			flag, err := checkOverloadLogTable(tt.typeTable, tt.maxI, tt.maxW, tt.maxE, tt.curStrNumb)
			require.NoError(t, err)
			assert.Equalf(t, tt.wantFlag, flag, "want:{%t}, recieved:{%t}", tt.wantFlag, flag)
		})
	}
}

// Test - TypeMessage
func Test_doSaving_SUCCESS(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	var msg = MessageT{
		TypeMessage:   "I",
		NameProject:   "project",
		LocationEvent: "cmd/main.go:65",
		BodyMessage:   "Not equal",
	}
	var tableName = "logI_1"

	mock.ExpectExec("INSERT INTO").WithArgs(msg.NameProject, msg.LocationEvent, msg.BodyMessage).WillReturnResult(sqlmock.NewResult(1, 1))

	ind, err := doSaving(db, tableName, msg)
	require.NoError(t, err)
	assert.Equal(t, int64(1), ind)
}

// Test - Save message + check overload log table + update name log table + create new table
func Test_savingMessageCheckResult_SUCCESS(t *testing.T) {

	msg := []MessageT{
		{
			TypeMessage:   "I",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:65",
			BodyMessage:   "Not equal",
		},
		{
			TypeMessage:   "W",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:66",
			BodyMessage:   "Not equal",
		},
		{
			TypeMessage:   "E",
			NameProject:   "project",
			LocationEvent: "cmd/main.go:67",
			BodyMessage:   "Not equal",
		},
	}

	tests := []struct {
		nameTest  string
		mockInit  func(mock sqlmock.Sqlmock)
		index     int
		nameTable string
		maxI      string
		maxW      string
		maxE      string
	}{
		{
			nameTest: "Store msg I. Not over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			index:     0,
			nameTable: "logI_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
		{
			nameTest: "Store msg I. Over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[0].NameProject, msg[0].LocationEvent, msg[0].BodyMessage).
					WillReturnResult(sqlmock.NewResult(6, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logI_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			index:     0,
			nameTable: "logI_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
		{
			nameTest: "Store msg W. Not over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[1].NameProject, msg[1].LocationEvent, msg[1].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			index:     1,
			nameTable: "logW_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
		{
			nameTest: "Store msg W. Over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[1].NameProject, msg[1].LocationEvent, msg[1].BodyMessage).
					WillReturnResult(sqlmock.NewResult(6, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logW_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			index:     1,
			nameTable: "logW_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
		{
			nameTest: "Store msg E. Not over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[2].NameProject, msg[2].LocationEvent, msg[2].BodyMessage).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			index:     2,
			nameTable: "logE_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
		{
			nameTest: "Store msg E. Over",
			mockInit: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO").
					WithArgs(msg[2].NameProject, msg[2].LocationEvent, msg[2].BodyMessage).
					WillReturnResult(sqlmock.NewResult(6, 1))

				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logE_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			index:     2,
			nameTable: "logE_1",
			maxI:      "5",
			maxW:      "5",
			maxE:      "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockInit(mock)

			err = savingMessageCheckResult(db, tt.nameTable, tt.maxI, tt.maxW, tt.maxE, msg[tt.index])
			require.NoError(t, err)
		})
	}
}

// Test - Check create table by name
func Test_checkCreateLogTable_SUCCESS(t *testing.T) {

	var nameDB = "test.db"
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").WillReturnResult(sqlmock.NewResult(0, 0))

	err = checkCreateLogTable(db, nameDB)
	require.NoError(t, err)
}

// Test - Initialisation the main table
func Test_initMainTable_SUCCESS(t *testing.T) {

	nameI := "logI_1"
	nameW := "logW_1"
	nameE := "logE_1"

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO main").
		WithArgs(nameI, nameW, nameE).
		WillReturnResult(sqlmock.NewResult(1, 1))

	num, err := initMainTable(db, nameI, nameW, nameE)
	require.NoError(t, err)
	assert.Equalf(t, int64(1), num, "wait 1, recieved:{%d}", num)

}

// Test - Reading log table names from the main table
func Test_readLogTablesName_SUCCESS(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoErrorf(t, err, "recieved error: {%v}", err)
	defer db.Close()

	mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
		WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
			AddRow("logI_1", "logW_1", "logE_1"))

	nameI, nameW, nameE, err := readLogTablesName(db)
	require.NoErrorf(t, err, "recieved err: {%v}", err)
	assert.Equalf(t, "logI_1", nameI, "wait {logI_1}, recieved {%s}", nameI)
	assert.Equalf(t, "logW_1", nameW, "wait {logW_1}, recieved {%s}", nameW)
	assert.Equalf(t, "logE_1", nameE, "wait {logE_1}, recieved {%s}", nameE)
}

// Test - Check-create main tables
func Test_checkCreateMainTable_SUCCESS(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = checkCreateMainTable(db)
	require.NoError(t, err)
}

// Test - Change the name of log table
func Test_changeLogTableNameCreate_SUCCESS(t *testing.T) {

	tests := []struct {
		nameTest  string
		initMock  func(mock sqlmock.Sqlmock)
		typeTable string
	}{
		{
			nameTest: "Change name I table",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logI_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			typeTable: "I",
		},
		{
			nameTest: "Change name W table",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logW_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			typeTable: "W",
		},
		{
			nameTest: "Change name E table",
			initMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT nameTableI, nameTableW, nameTableE FROM main WHERE id = 1").
					WillReturnRows(sqlmock.NewRows([]string{"nameTableI", "nameTableW", "nameTableE"}).
						AddRow("logI_1", "logW_1", "logE_1"))

				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logE_2").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("CREATE TABLE IF NOT EXISTS").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			typeTable: "E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.initMock(mock)

			err = changeLogTableNameCreate(db, tt.typeTable)
			require.NoError(t, err)
		})
	}
}

// Test - // Increment an index in the name log table
func Test_incrementIdInName_SUCCESS(t *testing.T) {

	curName := "logI_1"
	waitName := "logI_2"

	newName, err := incrementIdInName(curName)
	require.NoError(t, err)
	assert.Equalf(t, waitName, newName, "wait:{%s} recieved:{%s}", waitName, newName)
}

// Test - Update the name log table in the main table
func Test_updateNameLogTable_SUCCESS(t *testing.T) {

	tests := []struct {
		nameTest  string
		nameTable string
		nameType  string
		mocks     func(mock sqlmock.Sqlmock)
	}{
		{
			nameTest:  "SUCCESS change name LogI_1",
			nameTable: "logI_2",
			nameType:  "I",
			mocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logI_2").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			nameTest:  "SUCCESS change name LogW_1",
			nameTable: "logW_2",
			nameType:  "W",
			mocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logW_2").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			nameTest:  "SUCCESS change name LogE_1",
			nameTable: "logE_2",
			nameType:  "E",
			mocks: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE main SET nameTable").
					WithArgs("logE_2").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nameTest, func(t *testing.T) {

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mocks(mock)
			err = updateNameLogTable(db, tt.nameTable, tt.nameType)
			require.NoErrorf(t, err, "wait no error, but recieved: {%v}", err)
		})
	}
}

// =======================
// ==       FAULT       ==
// =======================

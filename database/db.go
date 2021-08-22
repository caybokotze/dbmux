package database

import (
	"database/sql"
	"fmt"
	"github.com/caybokotze/dbmux/configuration"
	"github.com/caybokotze/dbmux/logging"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func CreateConnectionToDbHost(configuration configuration.Configuration) (db *sql.DB, err error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		configuration.DbUser,
		configuration.DbPassword,
		configuration.DbHostIp,
		configuration.DbSchema)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	}
	return db, nil
}

func Query(db *sql.DB, q string) (*sql.Rows, error) {
	return db.Query(q)
}

func QueryRow(db *sql.DB, q string) *sql.Row {
	return db.QueryRow(q)
}

func ExecQuery(db *sql.DB, q string) (sql.Result, error) {
	if main.VerbosityEnabled {
		log.Printf("ExecQuery: %s\n", q)
	}
	return db.Exec(q)
}

func InsertLog(db *sql.DB, t *logging.Query) bool {
	insertSql := `
	insert into query_log(bindport, client, client_port, server, server_port, sql_type, 
	sql_string, create_time) values (%d, '%s', %d, '%s', %d, '%s', '%s', now())
	`
	_, err := ExecQuery(db, fmt.Sprintf(insertSql, t.BindPort, t.ClientIP, t.ClientPort, t.ServerIP, t.ServerPort, t.SqlType, t.SqlString))
	if err != nil {
		return false
	}
	return true
}
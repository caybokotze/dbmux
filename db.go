package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func databaseHost(configuration Configuration) (db *sql.DB, err error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		configuration.DbUser,
		configuration.DbPassword,
		configuration.DbHost,
		configuration.DbSchema)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		return db, err
	}
	return db, nil
}

func Query(db *sql.DB, q string) (*sql.Rows, error) {
	if VerbosityEnabled {
		log.Printf("Query: %s\n", q)
	}
	return db.Query(q)
}

func QueryRow(db *sql.DB, q string) *sql.Row {
	if VerbosityEnabled {
		log.Printf("Query: %s", q)
	}
	return db.QueryRow(q)
}

func ExecQuery(db *sql.DB, q string) (sql.Result, error) {
	if VerbosityEnabled {
		log.Printf("ExecQuery: %s\n", q)
	}
	return db.Exec(q)
}

func InsertLog(db *sql.DB, t *query) bool {
	insertSql := `
	insert into query_log(bindport, client, client_port, server, server_port, sql_type, 
	sql_string, create_time) values (%d, '%s', %d, '%s', %d, '%s', '%s', now())
	`
	_, err := ExecQuery(db, fmt.Sprintf(insertSql, t.bindPort, t.client, t.cport, t.server, t.sport, t.sqlType, t.sqlString))
	if err != nil {
		return false
	}
	return true
}
package database

import (
	"database/sql"
	"fmt"
	"github.com/caybokotze/dbmux/config"
	_ "github.com/go-sql-driver/mysql"
)

func CreateConnectionToDbHost(configuration config.Configuration) (db *sql.DB, err error) {
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
	return db.Exec(q)
}
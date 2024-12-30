package config

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectMySQL() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return errors.New("error opening database: " + err.Error())
	}
	err = DB.Ping()
	if err != nil {
		return errors.New("error connecting database: " + err.Error())
	}
	return nil
}

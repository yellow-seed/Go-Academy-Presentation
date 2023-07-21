package db

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	var err error
	user := os.Getenv("DBUser")
	pass := os.Getenv("DBPass")
	host := os.Getenv("DBHost")
	name := os.Getenv("DBName")

	db, err := sql.Open("mysql", user+":"+pass+"@("+host+":3306)/"+name+"?parseTime=true")
	if err != nil {
		return db, err
	}

	if err = db.Ping(); err != nil {
		return db, err
	}

	return db, nil
}
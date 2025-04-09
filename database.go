package godatabase

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {
	// Test opening a connection to the database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_golang?parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)  // 5 minutes
	db.SetConnMaxLifetime(60 * time.Minute) // 60 minutes

	return db
}

package godatabase

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMain(t *testing.T) {
	//
}

func TestOpenConnection(t *testing.T) {
	// Test opening a connection to the database
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_golang")
	if err != nil {
		fmt.Println("Error opening connection:", err)
		return
	}
	defer db.Close()

	// Check if the connection is valid
	if db == nil {
		fmt.Println("Database connection is nil")
		return
	}
	fmt.Println("Database connection opened successfully")
}

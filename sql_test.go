package godatabase

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := `INSERT INTO customer(id, name) VALUES('pajir', 'Pajir')`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Insert data success")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := `SELECT id, name FROM customer`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("ID:", id, "Name:", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := `SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance float32
		var rating float64
		var birth_date, created_at time.Time
		var married bool

		err = rows.Scan(&id, &name, &email, &balance, &rating, &birth_date, &married, &created_at)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("==============================================")
		fmt.Println("ID:", id)
		fmt.Println("Name:", name)
		if email.Valid {
			fmt.Println("Email:", email.String)
		}
		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)
		fmt.Println("Birth Date:", birth_date.Format("2006-01-02"))
		fmt.Println("Married:", married)
		fmt.Println("Created At:", created_at.Format("2006-01-02"))
	}

	fmt.Println("==============================================")
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	query := "SELECT username FROM users WHERE username = '" + username + "' AND password = '" + password + "' LIMIT 1"
	fmt.Println("Query:", query)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Login success, username:", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin"
	password := "admin"

	query := "SELECT username FROM users WHERE username = ? AND password = ? LIMIT 1"
	fmt.Println("Query:", query)
	rows, err := db.QueryContext(ctx, query, username, password)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err = rows.Scan(&username)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println("Login success, username:", username)
	} else {
		fmt.Println("Login failed")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "agan"
	password := "ugun"

	query := `INSERT INTO users(username, password) VALUES(?, ?)`
	_, err := db.ExecContext(ctx, query, username, password)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Insert data success")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "agan@gmail.com"
	comment := "oyyy"

	query := `INSERT INTO comments(email, comment) VALUES(?, ?)`
	result, err := db.ExecContext(ctx, query, email, comment)
	if err != nil {
		panic(err.Error())
	}

	id, err := result.LastInsertId()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Insert data success", "ID:", id)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	query := `INSERT INTO comments(email, comment) VALUES(?, ?)`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		email := "agan" + fmt.Sprint(i) + "@gmail.com"
		comment := "ugun" + fmt.Sprint(i) + " comment"

		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err.Error())
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Insert data success", "ID:", id)
		fmt.Println("Insert data success", "Email:", email, "Comment:", comment)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection() // mengambil koneksi database
	defer db.Close()      // menutup koneksi database setelah selesai

	ctx := context.Background() // membuat context untuk query

	tx, err := db.BeginTx(ctx, nil) // memulai transaksi
	if err != nil {                 // memeriksa error saat memulai transaksi
		panic(err.Error())
	}

	query := `INSERT INTO comments(email, comment) VALUES(?, ?)`
	stmt, err := tx.PrepareContext(ctx, query) // menyiapkan statement untuk transaksi
	if err != nil {                            // memeriksa error saat menyiapkan statement
		panic(err.Error())
	}
	defer stmt.Close() // menutup statement setelah selesai

	for i := 0; i < 10; i++ { // loop untuk memasukkan data 10 kali
		email := "agan" + fmt.Sprint(i) + "@gmail.com"
		comment := "ugun" + fmt.Sprint(i) + " comment"

		result, err := stmt.ExecContext(ctx, email, comment) // mengeksekusi statement dengan parameter email dan comment
		if err != nil {
			tx.Rollback() // jika terjadi error, rollback/batalkan transaksi
			panic(err.Error())
		}

		id, err := result.LastInsertId()
		if err != nil {
			tx.Rollback() // jika terjadi error, rollback/batalkan transaksi
			panic(err.Error())
		}

		fmt.Println("Insert data success", "ID:", id)
		fmt.Println("Insert data success", "Email:", email, "Comment:", comment)
	}

	err = tx.Commit() // commit/konfirmasi transaksi
	if err != nil {
		panic(err.Error())
	}
}

package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	//Need pointer to database database driver
	Database string
}

func DatabaseInit() *sql.DB {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//Logs Error to DB and prints error message for user to view.
func LogError(errorMessage error) {
	db := DB{}

	db.Insert()

	fmt.Errorf(errorMessage.Error())
}

func Log(message string) {

}

func (db *DB) Update() {

}

func (db *DB) Retrieve() {

}

func (db *DB) Insert() {

}

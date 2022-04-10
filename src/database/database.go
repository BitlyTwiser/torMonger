package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type DB struct {
	//Need pointer to database database driver
	Database *sql.DB
}

func DatabaseInit() DB {
	databaseUser := os.Getenv("POSTGRES_USER")
	databaseName := os.Getenv("POSTGRES_DB")
	databasePort := os.Getenv("5432")
	databaseUserPassword := os.Getenv("POSTGRES_PASSWORD")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=verify-full",
		databaseUser,
		databaseName,
		databaseUserPassword,
		databasePort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(fmt.Errorf("error initiating database connectino: %s", err.Error()))
	}

	database := DB{Database: db}

	return database
}

//Logs Error to DB and prints error message for user to view.
func (db *DB) LogError(errorMessage error) {

	db.Insert()

	fmt.Errorf(errorMessage.Error())
}

func (db *DB) Log(message string) {

}

func (db *DB) Update() {

}

func (db *DB) Retrieve() {

}

func (db *DB) Insert() {

}

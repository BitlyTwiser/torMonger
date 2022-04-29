package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
	"tor/src/logging"
)

type DB struct {
	//Need pointer to database database driver
	Database *pgxpool.Pool
}

func DatabaseInit() DB {
	databaseUser := os.Getenv("POSTGRES_USER")
	databaseName := os.Getenv("POSTGRES_DB")
	databaseHost := os.Getenv("POSTGRES_HOST")
	databasePort := os.Getenv("5432")
	databaseUserPassword := os.Getenv("POSTGRES_PASSWORD")
	//postgres://username:password@localhost:5432/database_name
	db, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", databaseUser, databaseUserPassword, databaseHost, databasePort, databaseName))
	if err != nil {
		logging.LogError(fmt.Errorf("Unable to connect to database: %v\n", err))
		os.Exit(1)
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

func (db *DB) Find(query string, resultSet interface{}) (interface{}, error) {
	defer db.Database.Close()

	err := db.Database.QueryRow(context.Background(), query).Scan(&resultSet)
	if err != nil {
		logging.LogError(fmt.Errorf("QueryRow failed: %v\n", err))
	}

	return resultSet, err
}

func (db *DB) Update() {

}

func (db *DB) Retrieve() {

}

func (db *DB) Insert() {

}

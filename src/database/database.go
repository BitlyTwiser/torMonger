package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

type DB struct {
	//Need pointer to database database driver
	Database *pgxpool.Pool
}

func loadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func DatabaseInit() DB {
	//Load Env Vars
	loadEnvFile()

	databaseUser := os.Getenv("POSTGRES_USER")
	databaseName := os.Getenv("POSTGRES_DB")
	databaseHost := os.Getenv("POSTGRES_HOST")
	databasePort := os.Getenv("POSTGRES_PORT")
	databaseUserPassword := os.Getenv("POSTGRES_PASSWORD")
	//postgres://username:password@localhost:5432/database_name
	databaseConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", databaseUser, databaseUserPassword, databaseHost, databasePort, databaseName)
	db, err := pgxpool.Connect(context.Background(), databaseConnectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	database := DB{Database: db}

	return database
}

func (db *DB) GenerateUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		db.LogError(fmt.Errorf("error generating UUID: %s", err.Error()))
	}

	return fmt.Sprintf("%v", uuid)
}

//Logs Error to DB and prints error message for user to view.
func (db *DB) LogError(errorMessage error) {
	_, err := db.Database.Exec(context.Background(), "INSERT INTO logs(id, error_message, notes) VALUES($1, $2, $3)", db.GenerateUUID(), errorMessage.Error(), "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in creating log record%v\n", err)
	}
}

func (db *DB) Log(message string) {

}

func (db *DB) FindLinkReference(link string, resultSet LinkReference) (LinkReference, error) {
	defer db.Database.Close()

	err := db.Database.QueryRow(context.Background(), "SELECT * FROM tormonger_data WHERE link_hash=$1", link).Scan(&resultSet)
	if err != nil {
		db.LogError(fmt.Errorf("queryRow failed: %v", err))
	}

	return resultSet, err
}

func (db *DB) Update() {

}

func (db *DB) Retrieve() {

}

func (db *DB) Insert() {

}

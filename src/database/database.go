package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"tor/src/types"

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
	_, err := db.Database.Exec(context.Background(), "INSERT INTO logs(id, log_message, log_type) VALUES($1, $2, $3)", db.GenerateUUID(), errorMessage.Error(), "Error")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in creating log record%v\n", err)
	}
}

func (db *DB) Log(message string) {
	_, err := db.Database.Exec(context.Background(), "INSERT INTO logs(id, log_message, log_type) VALUES($1, $2, $3)", db.GenerateUUID(), message, "Log")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in creating log record%v\n", err)
	}
}

func (db *DB) FindLinkReference(link string, resultSet LinkReference) (LinkReference, error) {
	rows, err := db.Database.Query(context.Background(), "SELECT id, link_hash, link FROM tormonger_data WHERE link_hash=$1", link)
	if err != nil {
		db.LogError(fmt.Errorf("queryRow failed: %v", err))
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&resultSet.Id, &resultSet.LinkHash, &resultSet.Link)
	}

	return resultSet, err
}

func (db *DB) FindSubDirectoryMatch(tormongerDataId, subdirectoryHash string, resultSet SubdirctoryReference) (SubdirctoryReference, error) {
	rows, err := db.Database.Query(context.Background(), "SELECT id, tormonger_data_id, subdirectory_path FROM tormonger_data_sub_directories WHERE tormonger_data_id=$1 AND subdirectory_path=$2", tormongerDataId, subdirectoryHash)
	if err != nil {
		db.LogError(fmt.Errorf("queryRow failed: %v", err))
	}
	defer rows.Close()

	//If no rows, make value in the tormong_data_subdirectories table
	for rows.Next() {
		rows.Scan(&resultSet.Id, &resultSet.TormongerDataId, &resultSet.SubdirectoryPath)
	}

	return resultSet, err
}

// Return ID
func (db *DB) CreateTormongDataRecord(linkHash, link string) string {
	var id string
	err := db.Database.QueryRow(context.Background(), "INSERT INTO tormonger_data(id, link_hash, link) VALUES($1, $2, $3) RETURNING id", db.GenerateUUID(), linkHash, link).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in creating log record%v\n", err)
	}

	return id
}

// Return ID
func (db *DB) CreateSubDirectoryRecord(link, subdirectoriesMatch, tormonger_id string) string {
	var id string
	err := db.Database.QueryRow(context.Background(), "INSERT INTO tormonger_data_sub_directories(id, tormonger_data_id, html_data_id, subdirectory_path) VALUES($1, $2, $3, $4) RETURNING id", db.GenerateUUID(), tormonger_id, nil, subdirectoriesMatch).Scan(&id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in creating log record%v\n", err)
	}

	return id
}

func (db *DB) CreateOrUpdateHtmlData(htmlData string, tormongerData types.TormongerDataValues, htmlReferenceData HtmlDataReference) {

}

func (db *DB) FindHtmlRecordForLink(tormongerDataId string, resultSet HtmlDataReference) (HtmlDataReference, error) {
	rows, err := db.Database.Query(context.Background(), "SELECT id, tormonger_data_id, tormonger_data_sub_directories_id, html_data FROM html_data WHERE tormonger_data_id=$1", tormongerDataId)
	if err != nil {
		db.LogError(fmt.Errorf("queryRow failed: %v", err))
	}
	defer rows.Close()

	//If no rows, make value in the tormong_data_subdirectories table
	for rows.Next() {
		rows.Scan(&resultSet.Id, &resultSet.TormongerDataId, &resultSet.TormongerDataSubDirectoriesId, &resultSet.HtmlData)
	}

	return resultSet, err
}

func (db *DB) Insert() {

}

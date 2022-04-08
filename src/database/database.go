package database

import "fmt"

type DB struct {
	//Need pointer to database database driver
	Database string
}

//Logs Error to DB and prints error message for user to view.
func LogError(errorMessage error) {
	db := DB{}

	db.Insert()

	fmt.Errorf(errorMessage.Error())
}

func (db *DB) Update() {

}

func (db *DB) Retrieve() {

}

func (db *DB) Insert() {

}

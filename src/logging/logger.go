package logging

import (
	"log"
	"tor/src/database"
)

var db = database.DatabaseInit()

// LogError Wraps the application in a generic logging utility that pushes logs into the database.
func LogError(message error) {
	log.Println(message)
	db.LogError(message)
}

func Log(message string) {
	log.Println(message)
	db.Log(message)
}

package logging

import (
	"log"
	"tor/src/database"
)

//Wraps the application in a generic logging utility that pushes logs into the database.
func LogError(message error) {
	log.Println(message)
	database.LogError(message)
}

func Log(message string) {
	log.Println(message)
	database.Log(message)
}

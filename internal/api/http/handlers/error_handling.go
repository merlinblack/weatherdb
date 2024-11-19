package handlers

import "log"

func quitOnError(message string, err error) {
	log.Fatalf("%s: %v\n", message, err)
}

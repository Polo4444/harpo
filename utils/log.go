package utils

import "log"

// LogFatalIfErr logs the error and exits the program if the error is not nil.
func LogFatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

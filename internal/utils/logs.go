package utils

import "log"

var Verbosity bool = true

func LogInfo(message string) {
	log.Println(message)
}

func LogVerbose(message string) {
	if Verbosity {
		log.Println(message)
	}
}

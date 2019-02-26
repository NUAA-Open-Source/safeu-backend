package common

import (
	"log"
	"os"
)

var (
	logFile *os.File
	err     error
)

func InitLogger() {

	os.Mkdir("log", os.ModePerm|os.ModeDir)
	logFile, err = os.OpenFile("log/safeu.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if !DEBUG {
		log.SetOutput(logFile)
	}
}

func GetLogFile() *os.File {
	return logFile
}

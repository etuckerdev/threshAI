package logging

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "threshAI: ", log.LstdFlags|log.Lshortfile)

type LogInterface interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

func NewLogger() LogInterface {
	return Logger
}

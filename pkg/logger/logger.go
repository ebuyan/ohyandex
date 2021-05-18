package logger

import (
	logger "log"
	"net/http"
)

func Error(r *http.Request, err error) {
	logger.Println("")
	logger.Println(r.URL)
	if err != nil {
		logger.Println("Error: " + err.Error())
	}
	logger.Println("")
}

func Info(r *http.Request, msg string) {
	logger.Println("")
	logger.Println(r.URL)
	logger.Println("Info: " + msg)
	logger.Println("")
}

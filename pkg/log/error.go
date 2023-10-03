package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func WriteErrorLog(errMessage error) {
	fileName := fmt.Sprintf("./storage/error_logs/error-%s.log", time.Now().Format("2006-01-02"))

	// open log file
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	defer logFile.Close()

	// set log out put
	log.SetOutput(logFile)

	log.SetFlags(log.LstdFlags)

	_, fileName, line, _ := runtime.Caller(1)
	log.Printf("[Error] in [%s:%d] %v", fileName, line, errMessage.Error())
}

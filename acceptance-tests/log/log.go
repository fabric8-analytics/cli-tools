package acclog

import (
	"log"
	"os"
)

// WarningLogger Logs warnings
var (
    WarningLogger *log.Logger
    InfoLogger    *log.Logger
    ErrorLogger   *log.Logger
)

// Initlog initializes the Log
func Initlog() {
	file, logerr := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
    if logerr != nil {
        log.Fatal(logerr)
    }

    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}
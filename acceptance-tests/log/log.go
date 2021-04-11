package acc_log

import (
	"log"
	"os"
)


var (
    WarningLogger *log.Logger
    InfoLogger    *log.Logger
    ErrorLogger   *log.Logger
)


func Init_log() {
	file, log_err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY, 0666)
    if log_err != nil {
        log.Fatal(log_err)
    }

    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}
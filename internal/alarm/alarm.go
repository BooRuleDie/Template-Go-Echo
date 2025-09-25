package alarm

import (
	"log"
)

// Usage: go alarm.GlobalAlarmer.Alarm("some alarm message")
var GlobalAlarmer Alarmer

const maxRetryCount = 3

// Default logger - will be changed to GlobalLogger in the future
var defaultLogger = log.New(log.Writer(), "[ALARM] ", log.LstdFlags|log.Lshortfile)

type Alarmer interface {
	Alarm(message string)
}

package utils

import "fmt"

var LoggingEnv = "LOGGING"

type MyLogger struct {
	logging bool
}

func NewMyLogger(logging bool) MyLogger {
	return MyLogger{logging: logging}
}

func (l *MyLogger) Log(message string) {
	if l.logging {
		fmt.Println(message)
	}
}

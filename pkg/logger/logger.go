package logger

import (
    "log"
)

// Info logs an informational message with a "[INFO]:" prefix.
// It takes a single parameter:
//   - message: The informational message to be logged.
func Info(message string) {
    log.Println("[INFO]:", message)
}

// Error logs the provided error message with an "[ERROR]" prefix.
// It uses the standard log package to print the error message to the standard output.
//
// Parameters:
//   - err: The error to be logged.
func Error(err error) {
    log.Println("[ERROR]:", err)
}
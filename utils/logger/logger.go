package logger

import (
	"fmt"
	"time"
)

func Info(message string) {
	fmt.Printf("[%s] INFO: %s\n", time.Now().Format("15:04:05"), message)
}

func Error(err error) {
	fmt.Printf("[%s] ERROR: %s\n", time.Now().Format("15:04:05"), err)
}

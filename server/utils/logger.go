package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func LogInfo(message string, args ...interface{}) {
	log.Printf("[INFO] "+message, args...)
}

func LogError(message string, args ...interface{}) {
	log.Printf("[ERROR] "+message, args...)
}

func LogWarning(message string, args ...interface{}) {
	log.Printf("[WARNING] "+message, args...)
}

func LogDebug(message string, args ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("[DEBUG] %s:%d - "+message, append([]interface{}{file, line}, args...)...)
		} else {
			log.Printf("[DEBUG] "+message, args...)
		}
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func PanicIfErr(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", message, err))
	}
}

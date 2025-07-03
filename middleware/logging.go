package middleware

import (
	"go.uber.org/zap"
	"log"
	"os"
)

var ZapLogger *zap.Logger

func InitLogger() {
	var err error
	if os.Getenv("ENV") == "production" {
		ZapLogger, err = zap.NewProduction()
	} else {
		ZapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Printf("Failed to initialize logger: %v", err)
	}
}

func CleanupLogger() {
	if ZapLogger != nil {
		ZapLogger.Sync() // Flushes any buffered log entries
	}
}

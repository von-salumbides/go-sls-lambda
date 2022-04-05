package logger

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() {
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Zap init failure")
	}
	zap.ReplaceGlobals(l)
	defer l.Sync()
}

func PANIC(message string, err error) {
	if err != nil {
		zap.L().Panic(message, zap.Any("error", err))
	}
}
func INFO(message string, data interface{}) {
	zap.L().Info(message, zap.Any("data", data))
}

package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	levelStr := os.Getenv("LOG_LEVEL")
	if levelStr == "" {
		levelStr = "info"
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}

	log.SetLevel(level)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}

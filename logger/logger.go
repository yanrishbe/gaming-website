package logger

import "github.com/sirupsen/logrus"

func New(level string) *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

package app

import "github.com/sirupsen/logrus"

func NewLogger() *logrus.Entry {
	logger := logrus.StandardLogger()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return logrus.NewEntry(logger)
}

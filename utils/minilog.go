package utils

import (
	"github.com/sirupsen/logrus"
)

var logLevel = logrus.DebugLevel

var log *logrus.Logger = nil

func init() {
	if log == nil {
		log = logrus.New()
	}
	log.SetLevel(logLevel)
	log.SetReportCaller(false)
}

func Info(args ...interface{}) {
	log.Info(args...);
}

func Warn(args ...interface{}) {
	log.Warn(args...);
}

func Debug(args ...interface{}) {
	log.Debug(args...);
}

func Fatal(args ...interface{}) {
	log.Fatal(args...);
}

func Error(args ...interface{}) {
	log.Error(args...);
}
package utils

import (
	"github.com/sirupsen/logrus"
)

func Info(args ...interface{}) {
	logrus.Info(args...);
}

func Warn(args ...interface{}) {
	logrus.Warn(args...);
}

func Debug(args ...interface{}) {
	logrus.Debug(args...);
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...);
}

func Error(args ...interface{}) {
	logrus.Error(args...);
}
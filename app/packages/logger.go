package packages

import "go.uber.org/zap"

var logger  *zap.SugaredLogger

func SetLogger(log *zap.SugaredLogger) {
	logger = log
}

func GetLogger() *zap.SugaredLogger {
	return logger
}


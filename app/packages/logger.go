package packages

import "go.uber.org/zap"

var Log  *zap.SugaredLogger

func SetLogger(log *zap.SugaredLogger) {
	Log = log
}


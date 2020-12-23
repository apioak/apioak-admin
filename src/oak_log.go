package src

import (
	log "github.com/sirupsen/logrus"
	"os"
)

type OakLog struct {
	config ConfigLog
}

func (ol *OakLog) New() error {
	log.SetFormatter(&log.JSONFormatter{})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.WarnLevel)

	return nil
}

func initLog(config ConfigLog) (*OakLog, error) {
	return nil, nil
}

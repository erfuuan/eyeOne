package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

func GetLogger() *zap.Logger {
	once.Do(func() {
		var err error
		log, err = zap.NewProduction()
		if err != nil {
			panic("cannot initialize zap logger: " + err.Error())
		}
	})
	return log
}

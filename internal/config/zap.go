package config

import (
	"go.uber.org/zap"
)

func InitZap() {
	if *GetConfig().Debug {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(logger)
	} else {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		zap.ReplaceGlobals(logger)
	}
}

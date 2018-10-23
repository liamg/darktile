package main

import (
	"fmt"

	"github.com/liamg/aminal/config"
	"go.uber.org/zap"
)

func getLogger(conf *config.Config) (*zap.SugaredLogger, error) {

	var logger *zap.Logger
	var err error
	if conf.DebugMode {
		logger, err = zap.NewDevelopment()
	} else {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Encoding = "console"
		logger, err = loggerConfig.Build()
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to create logger: %s", err)
	}
	return logger.Sugar(), nil
}

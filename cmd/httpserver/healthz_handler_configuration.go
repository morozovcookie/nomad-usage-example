package main

import (
	"go.uber.org/zap"
)

type HealthHandlerConfiguration struct {
	LoggerInstance *zap.Logger
}

func (cfg *HealthHandlerConfiguration) Logger() *zap.Logger {
	return cfg.LoggerInstance
}

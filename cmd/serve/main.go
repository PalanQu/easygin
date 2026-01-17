package main

import (
	"easygin/cmd/serve/cmd"
	"easygin/pkg/logging"

	"go.uber.org/zap"
)

func main() {
	logger := logging.GetGlobalLogger()
	if err := cmd.InitCmd(); err != nil {
		logger.Fatal("failed to initialize command", zap.Error(err))
	}
	if err := cmd.Execute(); err != nil {
		logger.Fatal("failed to execute command", zap.Error(err))
	}
}

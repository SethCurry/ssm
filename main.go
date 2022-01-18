package main

import (
	"github.com/SethCurry/ssm/internal/ssm"
	"go.uber.org/zap"
)

type ServiceConfig struct {
	StartCommand string
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	supConfig, err := ssm.LoadSupervisorConfig("./supervisor.yaml")
	if err != nil {
		logger.Fatal("failed to read supervisor config",
			zap.Error(err),
		)
	}

	serviceConfigs, err := ssm.LoadServiceDirectory(supConfig.ServiceDirectory)
	if err != nil {
		logger.Fatal("failed to load service configs",
			zap.Error(err),
		)
	}

	sup := ssm.NewSupervisor(logger)

	sup.Run()
}

package ssm

import (
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

type ServiceConfig struct {
	Name           string   `json:"name"`
	StartCommand   string   `json:"start_command"`
	StartArguments []string `json:"start_arguments"`
	StartOnBoot    bool     `json:"start_on_boot"`
}

type ServiceEntry struct {
	IsRunning       func() (bool, error)
	ShouldBeRunning func() bool
	Process         *os.Process
	Config          ServiceConfig
}

func NewSupervisor(logger *zap.Logger) *Supervisor {
	return &Supervisor{
		Services: make(map[string]*ServiceEntry),
		Logger:   logger,
	}
}

type Supervisor struct {
	Services map[string]*ServiceEntry
	Logger   *zap.Logger
}

func getServiceLogger(baseLogger *zap.Logger, svcName string) *zap.Logger {
	return baseLogger.With(
		zap.String("service", svcName),
	)
}

func (s *Supervisor) Run() error {
	for {
		for svcName, entry := range s.Services {
			logger := getServiceLogger(s.Logger, svcName)
			isRunning, err := entry.IsRunning()
			if err != nil {
				logger.Warn("failed to check whether service is running",
					zap.Error(err),
				)
				continue
			}

			if !isRunning && entry.ShouldBeRunning() {
				cmd := exec.Command(entry.Config.StartCommand, entry.Config.StartArguments...)
				err := cmd.Start()
				if err != nil {
					logger.Warn("failed to start service",
						zap.Error(err),
					)
					continue
				}

				entry.Process = cmd.Process
				logger.Info("successfully started service")
			}
		}

		time.Sleep(time.Second)
	}
}

package ssm

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// SupervisorConfig stores the configuration options for a supervisor instance.
type SupervisorConfig struct {
	ServiceDirectory string `yaml:"service_directory"`
}

func LoadSupervisorConfig(configPath string) (SupervisorConfig, error) {
	var ret SupervisorConfig

	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		return ret, errors.Wrapf(err, "failed to read config file: %s", configPath)
	}

	err = yaml.Unmarshal(contents, &ret)
	if err != nil {
		return ret, errors.Wrapf(err, "failed to unmarshal config file: %s", configPath)
	}

	return ret, nil
}

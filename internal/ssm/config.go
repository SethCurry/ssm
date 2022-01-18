package ssm

import (
	"io/ioutil"
	"path/filepath"
	"strings"

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

func LoadServiceDirectory(dirPath string) (map[string]ServiceConfig, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list services in directory \"%s\"", dirPath)
	}

	configs := make(map[string]ServiceConfig)

	for _, v := range files {
		fullPath := filepath.Join(dirPath, v.Name())

		contents, err := ioutil.ReadFile(fullPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read service file \"%s\"", fullPath)
		}

		var conf ServiceConfig

		err = yaml.Unmarshal(contents, &conf)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal service file \"%s\"", fullPath)
		}

		serviceName := strings.TrimSuffix(v.Name(), filepath.Ext(v.Name()))
		configs[serviceName] = conf
	}

	return configs, nil
}

package utils

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	SecretToken string
}

type ConfigManager struct {
	path string
}

func NewConfigManager(path string) *ConfigManager {
	return &ConfigManager{path}
}

func (c *ConfigManager) ReadConfig() (*Config, error) {
	result := &Config{}
	b, err := ioutil.ReadFile(c.path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

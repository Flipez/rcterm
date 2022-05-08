package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type RctermConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

func ReadConfig() RctermConfig {
	var rctermConfig RctermConfig
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	configPath := configDir + "/rcterm/config.yml"
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configFile, &rctermConfig)
	if err != nil {
		panic(err)
	}

	return rctermConfig
}

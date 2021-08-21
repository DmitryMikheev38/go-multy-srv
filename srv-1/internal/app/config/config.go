package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Port string `yaml:"port"`
	Postgres struct{
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		User string `yaml:"user"`
		Pwd string `yaml:"pwd"`
		DB string `yaml:"db"`
	} `yaml:"postgres"`
}

func GetConf() (*Config, error) {
	config := &Config{}
	yamlFile, err := ioutil.ReadFile("configs/app.yml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
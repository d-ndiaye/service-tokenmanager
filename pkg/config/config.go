package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	Port    int     `yaml:"port"`
	Sslport int     `yaml:"sslport"`
	Mongo   Mongodb `yaml:"mongodb"`
	Healthy Health  `yaml:"health"`
}
type Mongodb struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	Dbname string `yaml:"dbname"`
}

type Health struct {
	Period int `yaml:"period"`
}

func Load(file string) (error, Config) {
	var config Config
	_, err := os.Stat(file)
	if err != nil {
		return err, config
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("can't load config file: %s", err.Error()), config
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("can't unmarshal config file: %s", err.Error()), config
	}
	return nil, config
}

package service

import (
	"fmt"
	"io/ioutil"

	"github.com/realbucksavage/robin/pkg/database"
	"github.com/realbucksavage/robin/pkg/manage"
	"github.com/realbucksavage/robin/pkg/traffic"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Traffic    traffic.Config  `yaml:"traffic"`
	Management manage.Config   `yaml:"management"`
	Database   database.Config `yaml:"database"`
}

func readConfig(file string) (Config, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %s", err)
	}

	var c Config
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %s", err)
	}

	return c, nil
}

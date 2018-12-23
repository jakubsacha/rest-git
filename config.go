package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Repositories map[string]ConfigRepo
}

type ConfigRepo struct {
	Name string
	Url  string
}

func reloadConfig() (*Config, error) {
	tomlData, err := ioutil.ReadFile("rest-git.toml")
	if err != nil {
		return nil, err
	}

	var conf Config
	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

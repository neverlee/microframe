package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SrvConf struct {
	PluginPath string      `yaml:"plugin_path"`
	Plugins    boardConfig `yaml:"plugins"`
}

func NewSrvConf(configPath string) (rconf *SrvConf, rerr error) {
	var configBytes []byte
	var conf SrvConf
	var err error
	if configBytes, err = ioutil.ReadFile(configPath); err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(configBytes, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

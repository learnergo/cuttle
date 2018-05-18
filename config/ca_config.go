/**
ca 文件解析
**/

package config

import (
	"io/ioutil"

	"github.com/learnergo/loach/config"
	"gopkg.in/yaml.v2"
)

type Ca_Config struct {
	ECertClientConfig   *config.ClientConfig `yaml:"ecert"`
	TlsCertClientConfig *config.ClientConfig `yaml:"tlscert"`
}

func NewCaConfig(path string) (*Ca_Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := new(Ca_Config)
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

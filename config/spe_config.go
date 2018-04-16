package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type SpeConfig struct {
	Nodes []NodeConfig `yaml:"Nodes"`
}

type NodeConfig struct {
	Name     string         `yaml:"Name"`
	CaFile   string         `yaml:"CaFile"`
	Output   string         `yaml:"Output"`
	Register RegisterConfig `yaml:"Register"`
	Enroll   EnrollConfig   `yaml:"Enroll"`
}

type RegisterConfig struct {
	Registered     bool          `yaml:"Registered"`
	EnrollID       string        `yaml:"EnrollID"`
	Type           string        `yaml:"Type"`
	Secret         string        `yaml:"Secret"`
	MaxEnrollments int           `yaml:"MaxEnrollments"`
	Affiliation    string        `yaml:"Affiliation"`
	Attrs          []AttrsConfig `yaml:"Attrs"`
}

type AttrsConfig struct {
	Name  string `yaml:"Name"`
	Value string `yaml:"Value"`
	ECert bool   `yaml:"ECert,omitempty"`
}

type EnrollConfig struct {
	EnrollID string  `yaml:"EnrollID"`
	Secret   string  `yaml:"Secret"`
	Subject  Subject `yaml:"Subject"`
}

func NewSpeConfig(path string) (*SpeConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	speConfig := &SpeConfig{}

	err = yaml.Unmarshal(file, speConfig)
	if err != nil {
		return nil, err
	}

	return speConfig, nil
}

func (speConfig *SpeConfig) Marshal(path string) error {
	byteData, err := yaml.Marshal(speConfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, byteData, os.ModePerm)
}

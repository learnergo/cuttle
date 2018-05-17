/**
crypto-config 文件对应解析
**/
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type CryptoConfig struct {
	OrdererOrgs []OrdererConfig `yaml:"OrdererOrgs"`
	PeerOrgs    []PeerConfig    `yaml:"PeerOrgs"`
	Output      string          `yaml:"Output"`
	Subject     Subject         `yaml:"Subject"`
}

type OrdererConfig struct {
	Name   string `yaml:"Name"`
	CaFile string `yaml:"CaFile"`
	Domain string `yaml:"Domain"`
	Specs  []Spec `yaml:"Specs"`
}

type PeerConfig struct {
	Name     string   `yaml:"Name"`
	CaFile   string   `yaml:"CaFile"`
	Domain   string   `yaml:"Domain"`
	Specs    []Spec   `yaml:"Specs"`
	Template Template `yaml:"Template"`
	Users    Amount   `yaml:"Users"`
}

type Spec struct {
	Hostname   string `yaml:"Hostname"`
	CommonName string `yaml:"CommonName"`
}

type Template struct {
	Count int `yaml:"Count"`
	Start int `yaml:"Start"`
}

type Amount struct {
	Count int `yaml:"Count"`
}

type Subject struct {
	Country            string `yaml:"Country"`
	Province           string `yaml:"Province"`
	Locality           string `yaml:"Locality"`
	Organization       string `yaml:"Organization"`
	OrganizationalUnit string `yaml:"OrganizationalUnit"`
}

func NewCryptoConfig(path string) (*CryptoConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cConfig := &CryptoConfig{}

	err = yaml.Unmarshal(file, cConfig)
	if err != nil {
		return nil, err
	}

	return cConfig, nil
}

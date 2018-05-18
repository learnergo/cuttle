package loach

import (
	"github.com/learnergo/loach/config"
	"github.com/learnergo/loach/invoke"
	"github.com/learnergo/loach/model"
)

func NewClient(path string) (model.Client, error) {

	config, err := config.NewClientConfig(path)
	if err != nil {
		return nil, err
	}

	return invoke.NewClient(config)
}

func NewClientFromConfig(config *config.ClientConfig) (model.Client, error) {
	return invoke.NewClient(config)
}

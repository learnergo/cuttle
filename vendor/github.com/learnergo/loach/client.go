package loach

import (
	"github.com/learnergo/loach/config"
	"github.com/learnergo/loach/invoke"
	"github.com/learnergo/loach/model"
)

func NewClients(path string) (model.Clients, error) {

	config, err := config.NewClientConfig(path)
	if err != nil {
		return model.Clients{}, err
	}

	return invoke.NewClients(config)
}

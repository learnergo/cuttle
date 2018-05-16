package loach

import (
	"errors"

	"github.com/learnergo/loach/config"
	"github.com/learnergo/loach/crypto"
	"github.com/learnergo/loach/invoke"
	"github.com/learnergo/loach/model"
)

func NewClient(path string) (model.Client, error) {

	config, err := config.NewClientConfig(path)
	if err != nil {
		return nil, err
	}

	c, err := getCrypto(config.CryptoConfig)
	if err != nil {
		return nil, err
	}

	return invoke.NewClient(c, config)
}

func getCrypto(cc config.CryptoConfig) (crypto.Crypto, error) {
	switch cc.Family {
	case "ecdsa":
		c, err := crypto.NewCrypto(cc)
		if err != nil {
			return nil, err
		}
		return c, nil
	default:
		return nil, errors.New("Error Crypto")
	}
}

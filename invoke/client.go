package invoke

import (
	"github.com/learnergo/cuttle/config"
	"github.com/learnergo/loach"
	"github.com/learnergo/loach/model"
)

type cuttleClient struct {
	ECertClient   model.Client
	TlsCertClient model.Client
}

func newClient(path string) (cuttleClient, error) {
	caConfig, err := config.NewCaConfig(path)
	if err != nil {
		return cuttleClient{}, err
	}
	eCertClinet, err := loach.NewClientFromConfig(caConfig.ECertClientConfig)
	if err != nil {
		return cuttleClient{}, err
	}
	tlsCertClinet, err := loach.NewClientFromConfig(caConfig.TlsCertClientConfig)
	if err != nil {
		return cuttleClient{}, err
	}
	return cuttleClient{
		ECertClient:   eCertClinet,
		TlsCertClient: tlsCertClinet,
	}, nil
}

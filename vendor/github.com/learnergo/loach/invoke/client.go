/**
ca client真正实现
**/
package invoke

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/learnergo/loach/config"
	"github.com/learnergo/loach/crypto"
	"github.com/learnergo/loach/model"
)

type clientImpl struct {
	Url        string
	ServerName string
	Profile    string
	Algorithm  string
	Crypto     crypto.Crypto
	AdminKey   string
	AdminCert  string
}

func (client *clientImpl) GetAdmin() (*model.Identity, error) {
	return loadIdentity(client.AdminKey, client.AdminCert)
}

func (client *clientImpl) GetServer() (string, string) {
	return client.Url, client.ServerName
}

func (client *clientImpl) GetProfile() string {
	return client.Profile
}

func (client *clientImpl) Register(request *model.RegisterRequest) (*model.RegisterResponse, error) {
	return register(client, request)
}

func (client *clientImpl) Enroll(request *model.EnrollRequest) (*model.EnrollResponse, error) {
	return enroll(client, request)
}

func NewClient(config *config.ClientConfig) (model.Client, error) {
	//初始化ECertClient
	c, err := getCrypto(config.CryptoConfig)
	if err != nil {
		return nil, err
	}

	return &clientImpl{
		Url:        config.Url,
		ServerName: config.ServerName,
		Profile:    config.Profile,
		Crypto:     c,
		Algorithm:  config.Algorithm,
		AdminKey:   config.AdminKey,
		AdminCert:  config.AdminCert,
	}, nil
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

func (client *clientImpl) createAuthToken(identity *model.Identity, request []byte) (string, error) {

	encPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: identity.Cert.Raw})
	encCert := base64.StdEncoding.EncodeToString(encPem)
	body := base64.StdEncoding.EncodeToString(request)
	sigString := body + "." + encCert
	sig, err := client.Crypto.Sign([]byte(sigString), identity.Key)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", encCert, base64.StdEncoding.EncodeToString(sig)), nil
}

func (client *clientImpl) getTransport() *http.Transport {
	var tr *http.Transport
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return tr
}

func loadIdentity(keyPath string, certPath string) (*model.Identity, error) {
	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	pCert, _ := pem.Decode(certData)
	pKey, _ := pem.Decode(keyData)

	cert, err := x509.ParseCertificate(pCert.Bytes)
	if err != nil {
		return nil, err
	}
	key, err := x509.ParsePKCS8PrivateKey(pKey.Bytes)
	if err != nil {
		return nil, err
	}
	return &model.Identity{Cert: cert, Key: key}, nil
}

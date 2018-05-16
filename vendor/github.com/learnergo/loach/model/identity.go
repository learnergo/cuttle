/**
Identity 模型
**/
package model

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

type Identity struct {
	Key  interface{}
	Cert *x509.Certificate
	Ski  string
}

func MarshalIdentity(identity *Identity) (string, error) {

	key, cert, err := SplitIdentity(identity)
	if err != nil {
		return "", err
	}
	str, err := json.Marshal(map[string]string{"key": key, "cert": cert})
	if err != nil {
		return "", err
	}

	return string(str), nil
}

func SplitIdentity(identity *Identity) (string, string, error) {
	var key, cert string

	key, err := KeyToString(identity.Key)
	if err != nil {
		return "", "", err
	}
	cert = CertToString(identity.Cert)

	return key, cert, nil
}

func KeyToString(key interface{}) (string, error) {
	var result string
	switch key.(type) {
	case *ecdsa.PrivateKey:
		cast := key.(*ecdsa.PrivateKey)
		b, err := x509.MarshalECPrivateKey(cast)
		if err != nil {
			return "", err
		}
		block := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
		result = base64.StdEncoding.EncodeToString(block)

	default:
		return "", errors.New("Error peivate key Type")
	}
	return result, nil
}

func CertToString(cert *x509.Certificate) string {
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	certPem := pem.EncodeToMemory(certBlock)
	return base64.StdEncoding.EncodeToString(certPem)
}

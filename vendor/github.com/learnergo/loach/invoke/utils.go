package invoke

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"

	"github.com/learnergo/loach/model"
)

func stringToCert(data string) (*x509.Certificate, error) {
	rawCert, err := base64.StdEncoding.DecodeString(data)
	if err != nil {

		return nil, err
	}
	pemResult, _ := pem.Decode(rawCert)
	return x509.ParseCertificate(pemResult.Bytes)
}

func createCertificateRequest(request *model.CreateCsrRequest) ([]byte, error) {
	commonName := request.Name.CommonName
	if commonName == "" {
		return nil, errors.New("Missing CommonName")
	}
	subj := request.Name
	rawSubj := subj.ToRDNSequence()

	asn1Subj, err := asn1.Marshal(rawSubj)
	if err != nil {
		return nil, err
	}

	dnsAddr := make([]string, 0)

	hosts := request.Hosts

	for i := range hosts {
		dnsAddr = append(dnsAddr, hosts[i])
	}
	algorithm := request.Algorithm
	al, err := stringToAlgorithm(algorithm)
	if err != nil {
		return nil, err
	}
	template := x509.CertificateRequest{
		RawSubject:         asn1Subj,
		SignatureAlgorithm: al,
		DNSNames:           dnsAddr,
	}
	key := request.Key
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, err
	}
	csr := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})
	return csr, nil
}

func stringToAlgorithm(data string) (x509.SignatureAlgorithm, error) {
	var result x509.SignatureAlgorithm
	switch data {
	case "P256-SHA256":
		result = x509.ECDSAWithSHA256
	case "P384-SHA384":
		result = x509.ECDSAWithSHA384
	case "P521-SHA512":
		result = x509.ECDSAWithSHA512
	default:
		return 0, errors.New("Error algorithm")
	}
	return result, nil
}

func getSki(key interface{}) string {
	priKey := key.(*ecdsa.PrivateKey)
	raw := elliptic.Marshal(priKey.Curve, priKey.X, priKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)

	return hex.EncodeToString(hash.Sum(nil))
}

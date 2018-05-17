package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	_ "io/ioutil"
	"log"
	_ "os"

	"github.com/learnergo/cuttle/config"
	"github.com/learnergo/cuttle/node"
	"github.com/learnergo/cuttle/utils"

	"github.com/learnergo/loach"
	"github.com/learnergo/loach/model"
)

type CertType string

const (
	ConfigPath    string = "static/crypto-config.yaml"
	SpeConfigPath string = "static/cuttle.yaml"

	TlsCert CertType = "tls"
	ECert   CertType = "ecert"
	CaCert  CertType = "ca"
)

func main() {
	RunConfig()
}

func RunConfig() {
	//加载节点集
	cryptoSys, err := node.NewNode(ConfigPath)
	if err != nil {
		log.Printf("Failed to load cryptoSys ,err=%s", err)
		return
	}

	for _, value := range cryptoSys.PeerOrgs {
		if err := generatePeerOrg(value); err != nil {
			log.Printf("Failed to load nodes ,err=%s", err)
			return
		}
	}

	for _, value := range cryptoSys.OrdererOrgs {
		if err := generateOrdererOrg(value); err != nil {
			log.Printf("Failed to load nodes ,err=%s", err)
			return
		}
	}
}

func RunSpeConfig() {
	speConfig, err := config.NewSpeConfig(SpeConfigPath)
	if err != nil {
		log.Fatalf("Failed to load SpeConfig ,err=%s", err)
		return
	}
	for _, value := range speConfig.Nodes {
		register(value)
		enrollCert(ECert, value)
		enrollCert(TlsCert, value)
	}
}

func generatePeerOrg(peerOrg node.PeerOrg) error {
	//generate Admin
	register(peerOrg.Admin)
	enrollCert(ECert, peerOrg.Admin)
	enrollCert(TlsCert, peerOrg.Admin)
	//generate Users
	for _, value := range peerOrg.Users {
		register(value)
		enrollCert(ECert, value)
		enrollCert(TlsCert, value)
	}
	//generate Peers
	for _, value := range peerOrg.Peers {
		register(value)
		enrollCert(ECert, value)
		enrollCert(TlsCert, value)
	}
	return nil
}

func generateOrdererOrg(ordererOrg node.OrdererOrg) error {
	//generate Admin
	register(ordererOrg.Admin)
	enrollCert(ECert, ordererOrg.Admin)
	enrollCert(TlsCert, ordererOrg.Admin)
	//generate Orderers
	for _, value := range ordererOrg.Orderers {
		register(value)
		enrollCert(ECert, value)
		enrollCert(TlsCert, value)
	}
	return nil
}

func register(value config.NodeConfig) error {
	//加载ca客户端
	client, err := loach.NewClient(value.CaFile)
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return err
	}

	var attrs []model.RegisterAttribute
	for _, value := range value.Register.Attrs {
		attrs = append(attrs, model.RegisterAttribute{
			Name:  value.Name,
			Value: value.Value,
			ECert: value.ECert,
		})
	}

	request := &model.RegisterRequest{
		EnrollID:       value.Name,
		Type:           value.Register.Type,
		Secret:         value.Register.Secret,
		MaxEnrollments: value.Register.MaxEnrollments,
		Affiliation:    value.Register.Affiliation,
		Attrs:          attrs,
	}
	response, err := client.Register(request)
	if err != nil {
		log.Printf("Register %s Failed,because err=%s", value.Name, err)
		return err
	} else {
		if !response.Success || len(response.Errors) > 0 {
			log.Printf("Register %s Failed,because err=%s", value.Name, response.Error())
			return response.Error()
		}
		log.Printf("Register %s success,password=%s", value.Name, response.Result.Secret)
	}
	return nil
}

func enrollCert(certType CertType, value config.NodeConfig) error {

	//加载ca客户端
	client, err := loach.NewClient(value.CaFile)
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return err
	}

	request := &model.EnrollRequest{
		EnrollID: value.Enroll.EnrollID,
		Secret:   value.Enroll.Secret,
		Profile:  "tls",
		Name: pkix.Name{
			Country:            []string{value.Enroll.Subject.Country},
			Province:           []string{value.Enroll.Subject.Province},
			Locality:           []string{value.Enroll.Subject.Locality},
			Organization:       []string{value.Enroll.Subject.Organization},
			OrganizationalUnit: []string{value.Enroll.Subject.OrganizationalUnit},
			CommonName:         value.Enroll.EnrollID,
		},
	}

	if certType == TlsCert {
		request.Profile = "tls"
	}

	response, err := client.Enroll(request)
	if err != nil {
		log.Printf("Failed Enroll %s,err=%s", value.Name, err)
		return err
	} else {
		key, cert, err := model.SplitIdentity(response.Identity)
		ski := getSki(response.Identity.Key)
		if err != nil {
			log.Printf("Failed Enroll %s,err=%s", value.Name, err)
			return err
		}
		chain := model.CertToString(response.CertChain)
		log.Printf("Succeed to Enroll %s", value.Name)
		_, caName := client.GetServer()
		SaveIdentity(certType, ski+"_sk", value.Enroll.EnrollID, caName, value.Output, key, cert, chain)
	}
	return nil
}

func SaveIdentity(certType CertType, keyName, certName, caName, outPut, key, cert, chain string) {

	keyData, _ := base64.StdEncoding.DecodeString(key)
	key = string(keyData)
	certData, _ := base64.StdEncoding.DecodeString(cert)
	cert = string(certData)
	chainData, _ := base64.StdEncoding.DecodeString(chain)
	chain = string(chainData)

	switch certType {
	case ECert:

		outPut += "/" + "msp"
		//保存私钥
		keyPath := outPut + "/" + "keystore"
		utils.SaveFile(key, keyPath+"/"+keyName+".key")

		//保存证书
		certPath := outPut + "/" + "signcerts"
		utils.SaveFile(cert, certPath+"/"+certName+"-cert.pem")

		//保存证书链
		chainPath := outPut + "/" + "cacerts"
		utils.SaveFile(chain, chainPath+"/"+caName+"-cert.pem")

	case TlsCert:
		tlsOutPut := outPut + "/" + "tls"
		//保存私钥
		utils.SaveFile(key, tlsOutPut+"/"+"server.key")

		//保存证书
		utils.SaveFile(cert, tlsOutPut+"/"+"server.crt")

		//保存证书链
		utils.SaveFile(chain, tlsOutPut+"/"+"ca.crt")

		//保存tls证书链
		tlsChainPath := outPut + "/" + "msp" + "/" + "tlscacerts"
		utils.SaveFile(chain, tlsChainPath+"/"+caName+"-cert.pem")
	case CaCert:
	}
}

func getSki(key interface{}) string {
	priKey := key.(*ecdsa.PrivateKey)
	raw := elliptic.Marshal(priKey.Curve, priKey.X, priKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)

	return hex.EncodeToString(hash.Sum(nil))
}

//func ArrangeFiles(basePath string) {
//	files, err := ioutil.ReadDir(basePath)
//	if err != nil {
//		log.Printf("Failed to read %s dir,err=%s", basePath, err)
//	}
//	for _, value := range files {

//	}
//}

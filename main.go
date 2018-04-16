package main

import (
	"crypto/x509/pkix"
	"encoding/base64"
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
	nodes, err := node.NewNode(ConfigPath)
	if err != nil || len(nodes) == 0 {
		log.Printf("Failed to load nodes ,err=%s", err)
		return
	}

	//将节点集转为操作集
	speConfig, err := node.ParseNodesToSpeConfig(nodes)
	if err != nil || len(nodes) == 0 {
		log.Printf("Failed to Parse nodes to speConfig ,err=%s", err)
		return
	}

	//注册节点信息
	register(speConfig)
	//登记身份证书

	enrollCert(ECert, speConfig)

	//登记通讯证书
	enrollCert(TlsCert, speConfig)
}

func RunSpeConfig() {
	speConfig, err := config.NewSpeConfig(SpeConfigPath)
	if err != nil {
		log.Fatalf("Failed to load SpeConfig ,err=%s", err)
		return
	}
	//注册节点信息
	register(speConfig)
	//登记身份证书

	enrollCert(ECert, speConfig)

	//登记通讯证书
	enrollCert(TlsCert, speConfig)
}

func register(speConfig *config.SpeConfig) {
	for _, value := range speConfig.Nodes {
		//加载ca客户端
		client, err := loach.NewClient(value.CaFile)
		if err != nil {
			log.Fatalf("Failed to load ca client ,err=%s", err)
			continue
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
		} else {
			if !response.Success || len(response.Errors) > 0 {
				log.Printf("Register %s Failed,because err=%s", value.Name, response.Error())
			}
			log.Printf("Register %s success,password=%s", value.Name, response.Result.Secret)
		}
	}
}

func enrollCert(certType CertType, speConfig *config.SpeConfig) {
	for _, value := range speConfig.Nodes {

		//加载ca客户端
		client, err := loach.NewClient(value.CaFile)
		if err != nil {
			log.Fatalf("Failed to load ca client ,err=%s", err)
			continue
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
		} else {
			key, cert, err := model.SplitIdentity(response.Identity)
			if err != nil {
				log.Printf("Failed Enroll %s,err=%s", value.Name, err)
				return
			}
			chain := model.CertToString(response.CertChain)
			SaveIdentity(certType, value.Output, key, cert, chain)
		}
	}
}

func SaveIdentity(certType CertType, outPut, key, cert, chain string) {

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
		utils.SaveFile(key, keyPath+"/"+"key.key")

		//保存证书
		certPath := outPut + "/" + "signcerts"
		utils.SaveFile(cert, certPath+"/"+"cert.crt")

		//保存证书链
		chainPath := outPut + "/" + "cacerts"
		utils.SaveFile(chain, chainPath+"/"+"ca-cert.crt")

	case TlsCert:
		outPut += "/" + "tls"
		//保存私钥
		utils.SaveFile(key, outPut+"/"+"server.key")

		//保存证书
		utils.SaveFile(cert, outPut+"/"+"server.crt")

		//保存证书链
		utils.SaveFile(chain, outPut+"/"+"ca.crt")
	case CaCert:
	}
}

//func ArrangeFiles(basePath string) {
//	files, err := ioutil.ReadDir(basePath)
//	if err != nil {
//		log.Printf("Failed to read %s dir,err=%s", basePath, err)
//	}
//	for _, value := range files {

//	}
//}

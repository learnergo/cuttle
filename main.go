package main

import (
	"encoding/base64"
	"log"

	"github.com/learnergo/cuttle/node"
	"github.com/learnergo/cuttle/utils"
	"github.com/learnergo/loach"
	"github.com/learnergo/loach/model"
)

type CertType string

const (
	FilePath  string = "static\\crypto-config.yaml"
	CaPath    string = "static\\ca.yaml"
	AdminKey  string = "static\\admin.key"
	AdminCert string = "static\\admin.crt"

	TlsCert CertType = "tls"
	ECert   CertType = "ecert"
	CaCert  CertType = "ca"
)

func main() {
	//加载节点集
	nodes, err := node.NewNode(FilePath)
	if err != nil || len(nodes) == 0 {
		log.Fatalf("Failed to load nodes ,err=%s", err)
		return
	}
	//加载ca客户端
	client, err := loach.NewClient(CaPath)
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return
	}

	//加载ca管理员身份
	admin, err := loach.LoadIdentity(AdminKey, AdminCert)
	if err != nil {
		log.Fatalf("Failed to load ca admin ,err=%s", err)
		return
	}

	//注册节点信息
	register(nodes, client, admin)
	//登记身份证书

	enrollECert(nodes, client, admin)

	//登记通讯证书
	enrollTlsCert(nodes, client, admin)

}

func register(nodes []node.Node, client model.Client, admin *model.Identity) map[string]string {
	failList := make(map[string]string)

	for _, value := range nodes {
		nodeType := string([]byte(value.Type))
		var attrs []model.RegisterAttribute
		attrs = append(attrs, model.RegisterAttribute{Name: "hf.Registrar.Roles", Value: nodeType})
		attrs = append(attrs, model.RegisterAttribute{Name: "hf.Revoker", Value: "false"})

		request := &model.RegisterRequest{
			EnrollID:       value.Name,
			Type:           nodeType,
			Secret:         "adminpwd",
			MaxEnrollments: -1,
			Affiliation:    "*",
			Attrs:          attrs,
		}
		response, err := client.Register(admin, request)
		if err != nil {
			failList[value.Name] = err.Error()
		} else {
			if !response.Success || len(response.Errors) > 0 {
				failList[value.Name] = response.Error().Error()
			}
			log.Printf("Register %s success,password=%s", value.Name, response.Result.Secret)
		}
	}
	return failList
}

func enrollECert(nodes []node.Node, client model.Client, admin *model.Identity) map[string]string {
	failList := make(map[string]string)

	for _, value := range nodes {
		request := &model.EnrollRequest{
			EnrollID: value.Name,
			Secret:   "adminpwd",
			Name:     *value.Subject,
		}
		response, err := client.Enroll(request)
		if err != nil {
			failList[value.Name] = err.Error()
			log.Printf("Failed Enroll %s,err=%s", value.Name, err)
		} else {
			key, cert, err := model.SplitIdentity(response.Identity)
			if err != nil {
				log.Printf("Failed Enroll %s,err=%s", value.Name, err)
			}
			chain := model.CertToString(response.CertChain)
			SaveIdentity(ECert, value.Output, key, cert, chain)
		}
	}
	return failList
}

func enrollTlsCert(nodes []node.Node, client model.Client, admin *model.Identity) map[string]string {
	failList := make(map[string]string)

	for _, value := range nodes {
		request := &model.EnrollRequest{
			EnrollID: value.Name,
			Secret:   "adminpwd",
			Profile:  "tls",
			Name:     *value.Subject,
		}
		response, err := client.Enroll(request)
		if err != nil {
			log.Printf("Failed Enroll %s,err=%s", value.Name, err)
		} else {
			key, cert, err := model.SplitIdentity(response.Identity)

			if err != nil {
				failList[value.Name] = err.Error()
				log.Printf("Failed Enroll %s,err=%s", value.Name, err)
			}
			chain := model.CertToString(response.CertChain)

			SaveIdentity(TlsCert, value.Output, key, cert, chain)
		}
	}
	return failList
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

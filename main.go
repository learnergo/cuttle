package main

import (
	"encoding/base64"
	"log"

	"github.com/learnergo/cuttle/node"
	"github.com/learnergo/cuttle/utils"
	"github.com/learnergo/loach"
	"github.com/learnergo/loach/model"
)

const (
	FilePath  string = "static\\file.yaml"
	CaPath    string = "static\\ca.yaml"
	AdminKey  string = "static\\admin.key"
	AdminCert string = "static\\admin.crt"
)

func main() {
	//加载节点集
	nodes, err := node.NewNode(FilePath)
	if err != nil {
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
	for _, value := range nodes {
		t := string([]byte(value.Type))
		var attrs []model.RegisterAttribute
		attrs = append(attrs, model.RegisterAttribute{Name: "hf.Registrar.Roles", Value: t})
		attrs = append(attrs, model.RegisterAttribute{Name: "hf.Revoker", Value: "false"})

		request := &model.RegisterRequest{
			EnrolmentId:    value.Name,
			Type:           t,
			Secret:         "adminpwd",
			MaxEnrollments: -1,
			Affiliation:    "org1.department1",
			Attrs:          attrs,
		}
		response, err := client.Register(admin, request)
		if err != nil {
			log.Printf("Failed Register,err=%s", err)
		} else {
			if !response.Success {
				log.Printf("Failed Register,err=%s", response.Error())
			}
			if len(response.Errors) > 0 {
				log.Printf("Failed Register,err=%s", response.Error())
			}
			result := response.Result.Secret
			log.Printf("Register success,password=%s", result)
		}
	}

	//登记身份证书

	for _, value := range nodes {
		request := &model.EnrollRequest{
			EnrollID: value.Name,
			Secret:   "adminpwd",
			Name:     *value.Subject,
		}
		response, err := client.Enroll(request)
		if err != nil {
			log.Printf("Failed Enroll %s,err=%s", value.Name, err)
		} else {
			key, cert, err := model.SplitIdentity(response.Identity)
			if err != nil {
				log.Printf("Failed Enroll %s,err=%s", value.Name, err)
			}
			chain := model.CertToString(response.CertChain)
			SaveIdentity(true, value.Output, key, cert, chain)
		}
	}

	//登记通讯证书
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
				log.Printf("Failed Enroll %s,err=%s", value.Name, err)
			}
			chain := model.CertToString(response.CertChain)

			SaveIdentity(false, value.Output, key, cert, chain)
		}
	}
}

func SaveIdentity(eCert bool, outPut, key, cert, chain string) {

	keyData, _ := base64.StdEncoding.DecodeString(key)
	key = string(keyData)
	certData, _ := base64.StdEncoding.DecodeString(cert)
	cert = string(certData)
	chainData, _ := base64.StdEncoding.DecodeString(chain)
	chain = string(chainData)

	if eCert {
		outPut += "/" + "msp"
		//保存私钥
		keyPath := outPut + "/" + "keystore"
		utils.Mkdir(keyPath) //这两个方法应该合并
		utils.SaveFile(key, keyPath+"/"+"key.key")

		//保存证书
		certPath := outPut + "/" + "signcerts"
		utils.Mkdir(certPath)
		utils.SaveFile(cert, certPath+"/"+"cert.crt")

		//保存证书链
		chainPath := outPut + "/" + "cacerts"
		utils.Mkdir(chainPath)
		utils.SaveFile(chain, chainPath+"/"+"ca-cert.crt")

	} else {
		outPut += "/" + "tls"
		utils.Mkdir(outPut)
		//保存私钥
		utils.SaveFile(key, outPut+"/"+"server.key")

		//保存证书
		utils.SaveFile(cert, outPut+"/"+"server.crt")

		//保存证书链
		utils.SaveFile(chain, outPut+"/"+"ca.crt")
	}
}

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"log"

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
		clients, err := loach.NewClients(value.CaFile)

		//generate ECert
		err = generateECert(clients.ECertClient, value)
		if err != nil {
			return
		}
		//generate TlsCert
		err = generateTlsCert(clients.TlsCertClient, value)
		if err != nil {
			return
		}
	}
}

func generatePeerOrg(peerOrg node.PeerOrg) error {

	//加载ca客户端
	clients, err := loach.NewClients(peerOrg.Admin.CaFile)
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return err
	}

	//copy Admin
	sourceTarget := peerOrg.Admin.Output + "/" + "msp" + "/" + "signcerts" + "/" + peerOrg.Admin.Enroll.EnrollID + "-cert.pem"

	//generate Admin
	//generate ECert
	err = generateECert(clients.ECertClient, peerOrg.Admin)
	if err != nil {
		return err
	}
	//generate TlsCert
	err = generateTlsCert(clients.TlsCertClient, peerOrg.Admin)
	if err != nil {
		return err
	}

	//复制admin文件
	utils.CopyFile(sourceTarget, peerOrg.Admin.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")

	//创建org msp
	utils.CopyDir(peerOrg.Admin.Output+"/"+"msp"+"/"+"admincerts", peerOrg.RootPath+"/"+"admincerts")
	utils.CopyDir(peerOrg.Admin.Output+"/"+"msp"+"/"+"cacerts", peerOrg.RootPath+"/"+"cacerts")
	utils.CopyDir(peerOrg.Admin.Output+"/"+"msp"+"/"+"tlscacerts", peerOrg.RootPath+"/"+"tlscacerts")
	//generate Users
	for _, value := range peerOrg.Users {

		//generate ECert
		err = generateECert(clients.ECertClient, value)
		if err != nil {
			return err
		}
		//复制admin文件
		err = utils.CopyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")
		if err != nil {
			return err
		}

		//generate TlsCert
		err := generateTlsCert(clients.TlsCertClient, value)
		if err != nil {
			return err
		}
	}
	//generate Peers
	for _, value := range peerOrg.Peers {
		//generate ECert
		err = generateECert(clients.ECertClient, value)
		if err != nil {
			return err
		}
		//复制admin文件
		err = utils.CopyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")
		if err != nil {
			return err
		}

		//generate TlsCert
		err := generateTlsCert(clients.TlsCertClient, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateOrdererOrg(ordererOrg node.OrdererOrg) error {
	//加载ca客户端
	clients, err := loach.NewClients(ordererOrg.Admin.CaFile)
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return err
	}

	//copy Admin
	sourceTarget := ordererOrg.Admin.Output + "/" + "msp" + "/" + "signcerts" + "/" + ordererOrg.Admin.Enroll.EnrollID + "-cert.pem"

	//generate Admin
	//generate ECert
	err = generateECert(clients.ECertClient, ordererOrg.Admin)
	if err != nil {
		return err
	}
	//generate TlsCert
	err = generateTlsCert(clients.TlsCertClient, ordererOrg.Admin)
	if err != nil {
		return err
	}
	//复制admin文件
	utils.CopyFile(sourceTarget, ordererOrg.Admin.Output+"/"+"msp"+"/"+"admincerts/"+ordererOrg.Admin.Enroll.EnrollID+"-cert.pem")

	//创建org msp
	utils.CopyDir(ordererOrg.Admin.Output+"/"+"msp"+"/"+"admincerts", ordererOrg.RootPath+"/"+"admincerts")
	utils.CopyDir(ordererOrg.Admin.Output+"/"+"msp"+"/"+"cacerts", ordererOrg.RootPath+"/"+"cacerts")
	utils.CopyDir(ordererOrg.Admin.Output+"/"+"msp"+"/"+"tlscacerts", ordererOrg.RootPath+"/"+"tlscacerts")

	//generate Orderers
	for _, value := range ordererOrg.Orderers {
		//generate ECert
		err = generateECert(clients.ECertClient, value)
		if err != nil {
			return err
		}
		//复制admin文件
		err = utils.CopyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+ordererOrg.Admin.Enroll.EnrollID+"-cert.pem")
		if err != nil {
			return err
		}

		//generate TlsCert
		err := generateTlsCert(clients.TlsCertClient, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateECert(client model.Client, value config.NodeConfig) error {
	register(client, value)
	return enrollCert(client, ECert, value)
}

func generateTlsCert(client model.Client, value config.NodeConfig) error {
	register(client, value)
	return enrollCert(client, TlsCert, value)
}

func register(client model.Client, value config.NodeConfig) error {

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

func enrollCert(client model.Client, certType CertType, value config.NodeConfig) error {

	request := &model.EnrollRequest{
		EnrollID: value.Enroll.EnrollID,
		Profile:  client.GetProfile(),
		Secret:   value.Enroll.Secret,
		Name: pkix.Name{
			Country:            []string{value.Enroll.Subject.Country},
			Province:           []string{value.Enroll.Subject.Province},
			Locality:           []string{value.Enroll.Subject.Locality},
			Organization:       []string{value.Enroll.Subject.Organization},
			OrganizationalUnit: []string{value.Enroll.Subject.OrganizationalUnit},
			CommonName:         value.Enroll.EnrollID,
		},
	}

	var response *model.EnrollResponse
	var caName string

	if certType == TlsCert {
		request.Hosts = []string{value.Enroll.EnrollID}

	}
	_, caName = client.GetServer()
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

		saveIdentity(certType, ski+"_sk", value.Enroll.EnrollID, caName, value.Output, key, cert, chain)
	}
	return nil
}

func saveIdentity(certType CertType, keyName, certName, caName, outPut, key, cert, chain string) {

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
		utils.SaveFile(key, keyPath+"/"+keyName)

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

//保持与cryptogen中key名称生成规则一致
func getSki(key interface{}) string {
	priKey := key.(*ecdsa.PrivateKey)
	raw := elliptic.Marshal(priKey.Curve, priKey.X, priKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)

	return hex.EncodeToString(hash.Sum(nil))
}

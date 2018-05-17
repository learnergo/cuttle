package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path"

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

		err = enrollCert(ECert, value)
		if err != nil {
			log.Fatalf("Failed to enroll ECert ,err=%s", err)
			return
		}

		err = enrollCert(TlsCert, value)
		if err != nil {
			log.Fatalf("Failed to lenroll TlsCert ,err=%s", err)
			return
		}
	}
}

func generatePeerOrg(peerOrg node.PeerOrg) error {

	//copy Admin
	sourceTarget := peerOrg.Admin.Output + "/" + "msp" + "/" + "signcerts" + "/" + peerOrg.Admin.Enroll.EnrollID + "-cert.pem"

	//generate Admin
	//注册
	register(peerOrg.Admin)
	//登记cert
	err := enrollCert(ECert, peerOrg.Admin)
	if err != nil {
		return err
	}
	//登记tls cert
	err = enrollCert(TlsCert, peerOrg.Admin)
	if err != nil {
		return err
	}
	//复制admin文件
	err = copyFile(sourceTarget, peerOrg.Admin.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")

	//generate Users
	for _, value := range peerOrg.Users {
		//注册
		register(value)
		//登记ecert
		err = enrollCert(ECert, value)
		if err != nil {
			return err
		}
		//复制admin文件
		err = copyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")
		if err != nil {
			return err
		}
		//登记tlscert
		err = enrollCert(TlsCert, value)
		if err != nil {
			return err
		}
	}
	//generate Peers
	for _, value := range peerOrg.Peers {
		//注册
		register(value)
		//登记ecert
		err = enrollCert(ECert, value)
		if err != nil {
			return err
		}
		//复制admin文件
		err = copyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+peerOrg.Admin.Enroll.EnrollID+"-cert.pem")
		if err != nil {
			return err
		}
		//登记tlscert
		err = enrollCert(TlsCert, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateOrdererOrg(ordererOrg node.OrdererOrg) error {
	//copy Admin
	sourceTarget := ordererOrg.Admin.Output + "/" + "msp" + "/" + "signcerts" + "/" + ordererOrg.Admin.Enroll.EnrollID + "-cert.pem"

	//generate Admin
	//注册
	register(ordererOrg.Admin)
	//登记ecert
	err := enrollCert(ECert, ordererOrg.Admin)
	if err != nil {
		return err
	}
	//复制admin文件
	copyFile(sourceTarget, ordererOrg.Admin.Output+"/"+"msp"+"/"+"admincerts/"+ordererOrg.Admin.Enroll.EnrollID+"-cert.pem")

	err = enrollCert(TlsCert, ordererOrg.Admin)
	if err != nil {
		return err
	}
	//generate Orderers
	for _, value := range ordererOrg.Orderers {
		//注册
		register(value)
		//登记ecert
		err = enrollCert(ECert, value)
		if err != nil {
			return err
		}
		//复制admin文件
		copyFile(sourceTarget, value.Output+"/"+"msp"+"/"+"admincerts/"+ordererOrg.Admin.Enroll.EnrollID+"-cert.pem")
		//登记tlscert
		err = enrollCert(TlsCert, value)
		if err != nil {
			return err
		}
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

//保持与cryptogen中key名称生成规则一致
func getSki(key interface{}) string {
	priKey := key.(*ecdsa.PrivateKey)
	raw := elliptic.Marshal(priKey.Curve, priKey.X, priKey.Y)

	// Hash it
	hash := sha256.New()
	hash.Write(raw)

	return hex.EncodeToString(hash.Sum(nil))
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	dir := path.Dir(dst)
	utils.Mkdir(dir)
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

//func ArrangeFiles(basePath string) {
//	files, err := ioutil.ReadDir(basePath)
//	if err != nil {
//		log.Printf("Failed to read %s dir,err=%s", basePath, err)
//	}
//	for _, value := range files {

//	}
//}

package node

import (
	"crypto/x509/pkix"
	"fmt"

	"github.com/learnergo/cuttle/config"
	"github.com/learnergo/cuttle/constant"
)

type CryptoSys struct {
	PeerOrgs    []PeerOrg
	OrdererOrgs []OrdererOrg
}

type PeerOrg struct {
	Peers []config.NodeConfig
	Admin config.NodeConfig
	Users []config.NodeConfig
}

type OrdererOrg struct {
	Orderers []config.NodeConfig
	Admin    config.NodeConfig
}

type Node struct {
	Name    string
	CaFile  string
	Type    constant.NodeType
	Subject *pkix.Name
	Output  string
	OrgName string
}

func NewNode(path string) (*CryptoSys, error) {
	cConfig, err := config.NewCryptoConfig(path)
	if err != nil {
		return nil, err
	}

	return parseConfigToNodes(cConfig)
}

func parseConfigToNodes(cConfig *config.CryptoConfig) (*CryptoSys, error) {

	cryptoSys := &CryptoSys{}

	//解析orderer
	for _, value := range cConfig.OrdererOrgs {
		ordererOrg := OrdererOrg{}
		for _, o := range value.Specs {
			n := Node{}
			n.Name = fmt.Sprintf("%s.%s", o.Hostname, value.Domain)
			if o.CommonName != "" {
				//CommonName 覆盖
				n.Name = o.CommonName
			}
			n.Type = constant.Orderer
			n.CaFile = value.CaFile
			n.Subject = getSubject(&cConfig.Subject, n.Name)
			n.Output = getOrdererOutput(n.Name, value.Domain, n.Type, cConfig.Output)
			n.OrgName = value.Name
			v, err := parseNodeToNodeConfig(n)
			if err != nil {
				return nil, err
			}
			ordererOrg.Orderers = append(ordererOrg.Orderers, v)
		}
		//orderer 默认加一个Admin
		for _, item := range []string{"Admin"} {
			n := Node{}
			n.Name = fmt.Sprintf("%s@%s", item, value.Domain)
			if item == "Admin" {
				n.Type = constant.Admin
			} else {
				n.Type = constant.User
			}
			n.CaFile = value.CaFile
			n.Subject = getSubject(&cConfig.Subject, n.Name)
			n.Output = getOrdererOutput(n.Name, value.Domain, n.Type, cConfig.Output)
			n.OrgName = value.Name
			v, err := parseNodeToNodeConfig(n)
			if err != nil {
				return nil, err
			}
			ordererOrg.Admin = v
		}
		cryptoSys.OrdererOrgs = append(cryptoSys.OrdererOrgs, ordererOrg)
	}

	//解析peer
	for _, value := range cConfig.PeerOrgs {
		peerOrg := PeerOrg{}
		peers := getPeers(value.Specs, &value.Template, value.Domain)
		for _, p := range peers {
			n := Node{}
			n.Name = p
			n.Type = constant.Peer
			n.CaFile = value.CaFile
			n.Subject = getSubject(&cConfig.Subject, n.Name)
			n.Output = getPeerOutput(n.Name, value.Domain, n.Type, cConfig.Output)
			n.OrgName = value.Name
			v, err := parseNodeToNodeConfig(n)
			if err != nil {
				return nil, err
			}
			peerOrg.Peers = append(peerOrg.Peers, v)
		}

		u := []string{"Admin"}
		users := getUsers(value.Users.Count)
		u = append(u, users...)
		for _, item := range u {
			n := Node{}
			n.Name = fmt.Sprintf("%s@%s", item, value.Domain)
			if item == "Admin" {
				n.Type = constant.Admin
			} else {
				n.Type = constant.User
			}
			n.CaFile = value.CaFile
			n.Subject = getSubject(&cConfig.Subject, n.Name)
			n.Output = getPeerOutput(n.Name, value.Domain, n.Type, cConfig.Output)
			n.OrgName = value.Name
			v, err := parseNodeToNodeConfig(n)
			if err != nil {
				return nil, err
			}
			if item == "Admin" {
				peerOrg.Admin = v
			} else {
				peerOrg.Users = append(peerOrg.Users, v)
			}
		}
		cryptoSys.PeerOrgs = append(cryptoSys.PeerOrgs, peerOrg)
	}

	return cryptoSys, nil
}

func getSubject(subject *config.Subject, commonName string) *pkix.Name {
	return &pkix.Name{
		Country: []string{subject.Country},
		//Organization:[]string{subject.Organization},
		//OrganizationalUnit:[]string{subject.OrganizationalUnit},
		Locality:   []string{subject.Country},
		Province:   []string{subject.Country},
		CommonName: commonName,
	}
}

func getOrdererOutput(name string, domain string, nodeType constant.NodeType, basePath string) string {
	path0 := basePath
	path1 := "ordererOrganizations"
	path2 := domain
	path3 := ""
	if nodeType == constant.Orderer {
		path3 = "orderers"
	} else {
		path3 = "users"
	}
	path4 := name
	return fmt.Sprintf("%s/%s/%s/%s/%s", path0, path1, path2, path3, path4)
}

func getPeerOutput(name string, domain string, nodeType constant.NodeType, basePath string) string {
	path0 := basePath
	path1 := "peerOrganizations"
	path2 := domain
	path3 := ""
	if nodeType == constant.Peer {
		path3 = "peers"
	} else {
		path3 = "users"
	}
	path4 := name
	return fmt.Sprintf("%s/%s/%s/%s/%s", path0, path1, path2, path3, path4)
}

func getPeers(specs []config.Spec, temp *config.Template, domain string) []string {
	var result []string
	if specs != nil && len(specs) > 0 {
		for _, value := range specs {
			if value.CommonName != "" {
				result = append(result, value.CommonName)
			} else {
				result = append(result, fmt.Sprintf("%s.%s", value.Hostname, domain))
			}
		}
		return result
	}
	count := temp.Count
	start := temp.Start
	for i := start; i < count; i++ {
		result = append(result, fmt.Sprintf("peer%d.%s", i, domain))
	}
	return result
}

func getUsers(count int) []string {
	var result []string
	if count > 0 {
		for i := 1; i <= count; i++ {
			result = append(result, fmt.Sprintf("User%d", i))
		}
		return result
	}
	return nil
}

func parseNodeToNodeConfig(value Node) (config.NodeConfig, error) {

	return config.NodeConfig{
		Name:   value.Name,
		CaFile: value.CaFile,
		Output: value.Output,
		Register: config.RegisterConfig{
			Registered:     false,
			EnrollID:       value.Name,
			Type:           string([]byte(value.Type)),
			Secret:         "adminpwd",
			MaxEnrollments: -1,
			Affiliation:    ".",
			Attrs: []config.AttrsConfig{config.AttrsConfig{
				Name:  "hf.Registrar.Roles",
				Value: string([]byte(value.Type)),
			}, config.AttrsConfig{
				Name:  "hf.Revoker",
				Value: "false",
			}},
		},
		Enroll: config.EnrollConfig{
			EnrollID: value.Name,
			Secret:   "adminpwd",
			Subject:  ParseSubject(value.Subject),
		},
	}, nil
}

func ParseSubject(subject *pkix.Name) config.Subject {
	result := config.Subject{}
	if subject.Country != nil && len(subject.Country) > 0 {
		result.Country = subject.Country[0]
	}
	if subject.Province != nil && len(subject.Province) > 0 {
		result.Province = subject.Province[0]
	}
	if subject.Locality != nil && len(subject.Locality) > 0 {
		result.Locality = subject.Locality[0]
	}
	if subject.Organization != nil && len(subject.Organization) > 0 {
		result.Organization = subject.Organization[0]
	}
	if subject.OrganizationalUnit != nil && len(subject.OrganizationalUnit) > 0 {
		result.OrganizationalUnit = subject.OrganizationalUnit[0]
	}
	return result
}

package main

import (
	"log"

	"github.com/learnergo/cuttle/node"
	"github.com/learnergo/loach"
	"github.com/learnergo/loach/model"
)

func main() {
	//加载节点集
	nodes, err := node.NewNode("static\\file.yaml")
	if err != nil {
		log.Fatalf("Failed to load nodes ,err=%s", err)
		return
	}
	//加载ca客户端
	client, err := loach.NewClient("static\\ca.yaml")
	if err != nil {
		log.Fatalf("Failed to load ca client ,err=%s", err)
		return
	}

	//加载ca管理员身份
	admin, err := loach.LoadIdentity("static\\admin.key", "static\\admin.crt")
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

	//登记通讯证书
}

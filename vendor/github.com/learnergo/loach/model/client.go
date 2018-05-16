/**
ca client 接口定义
**/
package model

type Client interface {
	GetAdmin() (*Identity, error)
	GetServer() (url string, serverName string)

	Register(*RegisterRequest) (*RegisterResponse, error)
	Enroll(*EnrollRequest) (*EnrollResponse, error)
}

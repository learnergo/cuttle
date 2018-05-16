/**
一些常量定义
**/
package constant

type RequestApi string

const (
	Register RequestApi = "/api/v1/register"
	Enroll   RequestApi = "/api/v1/enroll"
)

const RequestMethod string = "POST"

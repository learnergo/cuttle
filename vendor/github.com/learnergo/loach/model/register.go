/**
register 过程所用模型
**/
package model

import (
	"crypto/x509/pkix"
)

type RegisterRequest struct {
	EnrollID       string              `json:"id"`
	Type           string              `json:"type"`
	Secret         string              `json:"secret,omitempty"`
	MaxEnrollments int                 `json:"max_enrollments,omitempty"`
	Affiliation    string              `json:"affiliation"`
	Attrs          []RegisterAttribute `json:"attrs"`
	CAName         string              `json:"caname,omitempty"`
}

type RegisterAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	ECert bool   `json:"ecert,omitempty"`
}

type RegisterResponse struct {
	Response
	Result RegisterCredentialResponse `json:"result"`
}

type RegisterCredentialResponse struct {
	Secret string `json:"secret"`
}

type CreateCsrRequest struct {
	Name      pkix.Name
	Key       interface{}
	Hosts     []string
	Algorithm string
}

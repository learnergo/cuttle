/**
enroll 过程所用模型
**/
package model

import (
	"crypto/x509"
	"crypto/x509/pkix"
)

type EnrollRequest struct {
	EnrollID string
	Secret   string
	Name     pkix.Name
	Profile  string            `json:"profile,omitempty"`
	Label    string            `json:"label,omitempty"`
	CAName   string            `json:"caname,omitempty"`
	Hosts    []string          `json:"hosts,omitempty"`
	Attrs    []EnrollAttribute `json:"attr_reqs,omitempty"`
}

type EnrollAttribute struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional,omitempty"`
}

type EnrollResponse struct {
	Identity  *Identity
	CertChain *x509.Certificate
}

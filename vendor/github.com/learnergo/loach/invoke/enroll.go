/**
enroll 过程实现
**/
package invoke

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/learnergo/loach/constant"
	"github.com/learnergo/loach/model"
)

type enrollmentResponse struct {
	model.Response
	Result enrollmentResponseResult `json:"result"`
}

type enrollmentResponseResult struct {
	Cert       string
	ServerInfo enrollmentResponseServerInfo
	Version    string
}

type certificateRequest struct {
	model.EnrollRequest
	CR string `json:"certificate_request"`
}

func (e *enrollmentResponseResult) UnmarshalJSON(b []byte) error {
	type tmpStruct struct {
		Cert       string
		ServerInfo enrollmentResponseServerInfo
		Version    string
	}
	if len(b) > 2 {
		r := new(tmpStruct)
		err := json.Unmarshal(b, r)
		if err != nil {
			return err
		}
		e.Cert = r.Cert
		e.ServerInfo = r.ServerInfo
		e.Version = r.Version
	}
	return nil
}

type enrollmentResponseServerInfo struct {
	CAName  string
	CAChain string
}

func enroll(client *clientImpl, request *model.EnrollRequest) (*model.EnrollResponse, error) {
	//create private key
	key, err := client.Crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	var hosts []string
	if len(request.Hosts) == 0 {
		parsedUrl, err := url.Parse(client.Url)
		if err != nil {
			return nil, err
		}
		hosts = []string{parsedUrl.Host}
	} else {
		hosts = request.Hosts
	}
	csrRequest := &model.CreateCsrRequest{
		Name:      request.Name,
		Key:       key,
		Hosts:     hosts,
		Algorithm: client.Algorithm,
	}

	csr, err := createCertificateRequest(csrRequest)
	if err != nil {
		return nil, err
	}

	crm, err := json.Marshal(certificateRequest{CR: string(csr), EnrollRequest: *request})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(constant.RequestMethod, fmt.Sprintf("%s%s", client.Url, constant.Enroll), bytes.NewBuffer(crm))

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(request.EnrollID, request.Secret)

	httpClient := &http.Client{Transport: client.getTransport()}
	resp, err := httpClient.Do(req)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		enrResp := new(enrollmentResponse)
		if err := json.Unmarshal(body, enrResp); err != nil {
			return nil, err
		}
		if !enrResp.Success {
			return nil, enrResp.Error()
		}
		cert, err := stringToCert(enrResp.Result.Cert)
		if err != nil {

			return nil, err
		}
		certChain, err := stringToCert(enrResp.Result.ServerInfo.CAChain)
		if err != nil {
			return nil, err
		}
		return &model.EnrollResponse{Identity: &model.Identity{Cert: cert, Key: key}, CertChain: certChain}, nil
	}
	return nil, fmt.Errorf("non 200 response: %v message is: %s", resp.StatusCode, string(body))
}

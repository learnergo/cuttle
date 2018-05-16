/**
register 过程实现
**/
package invoke

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/learnergo/loach/constant"
	"github.com/learnergo/loach/model"
)

func register(client *clientImpl, request *model.RegisterRequest) (*model.RegisterResponse, error) {
	if request.EnrollID == "" {
		return nil, errors.New("Missing EnrollmentID")
	}
	if request.Affiliation == "" {
		return nil, errors.New("Missing Affiliation")
	}
	if request.Type == "" {
		return nil, errors.New("Missing Type")
	}
	admin, _ := client.GetAdmin()

	if admin == nil {
		return nil, errors.New("nil admin")
	}
	reqJson, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(constant.RequestMethod, fmt.Sprintf("%s%s", client.Url, constant.Register), bytes.NewBuffer(reqJson))

	token, err := client.createAuthToken(admin, reqJson)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("authorization", token)
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Transport: client.getTransport()}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		result := new(model.RegisterResponse)
		if err := json.Unmarshal(body, result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, fmt.Errorf("non 200 response: %v message is: %s", resp.StatusCode, string(body))
}

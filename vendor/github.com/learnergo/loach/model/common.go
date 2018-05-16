/**
通用模型
**/
package model

import (
	"fmt"
)

type Response struct {
	Success  bool          `json:"success"`
	Errors   []ResponseErr `json:"errors"`
	Messages []string      `json:"messages"`
}

type ResponseErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r Response) Error() error {
	errors := ""
	for _, e := range r.Errors {
		errors += e.Message + ": "
	}
	return fmt.Errorf(errors)
}

package response

import (
	"encoding/json"
	"fmt"
)

type Response struct {
	Message string `json:"message"`
}

func NewResponse(message string) *Response {
	return &Response{
		Message: message,
	}
}

func ResponseBody(message string) string {
	res := NewResponse(message)
	body, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		body = []byte("error")
	}
	return string(body)
}

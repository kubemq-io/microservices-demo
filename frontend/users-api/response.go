package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"runtime"

	"github.com/labstack/echo"
)

var goRuntime = runtime.Version()

type Response struct {
	IsError  bool            `json:"is_error"`
	Message  string          `json:"message"`
	Data     json.RawMessage `json:"data"`
	httpCode int
	c        echo.Context
	data     interface{}
	request  interface{}
	method   string
	his      *History
	kube     *KubeMQ
}

func NewResponse(c echo.Context, mq *KubeMQ, method, requestType string) *Response {
	res := &Response{
		c:        c,
		httpCode: 200,
		method:   c.Request().URL.Path,
		his: &History{
			Id:           uuid.New().String(),
			Source:       "users-api",
			Time:         time.Now(),
			Type:         requestType,
			Method:       method,
			Request:      "",
			Response:     "",
			IsError:      false,
			ErrorMessage: "",
		},
		kube: mq,
	}
	res.setResponseHeaders()
	return res
}
func (res *Response) setResponseHeaders() *Response {
	res.c.Response().Header().Set("X-Runtime", goRuntime)
	return res
}

func (res *Response) SetError(err error) *Response {
	res.IsError = true
	res.Message = err.Error()
	res.his.IsError = true
	res.his.ErrorMessage = err.Error()
	return res
}

func (res *Response) SetErrorWithText(errText string) *Response {
	res.IsError = true
	res.Message = errText
	return res
}

func (res *Response) SetResponseBody(data interface{}) *Response {
	res.data = data
	res.his.Response = PrettyJson(data)
	return res
}
func (res *Response) SetRequestBody(data interface{}) *Response {
	res.request = data
	res.his.Request = PrettyJson(data)
	return res
}
func (res *Response) SetHttpCode(value int) *Response {
	res.httpCode = value
	return res
}
func (res *Response) Send() error {
	buffer, err := json.Marshal(res.data)
	if err != nil {
		res.SetError(err)
		return res.c.JSONPretty(res.httpCode, res, "\t")
	}
	res.Data = buffer
	if !res.IsError {
		res.Message = "OK"
	}
	log.Println(fmt.Sprintf("New Call Received:\nmethod: %s\nrequest: %sresponse: %smessage: %s", res.method, PrettyJson(res.request), PrettyJson(res.data), res.Message))
	go res.kube.SendHistory(context.Background(), res.his)
	return res.c.JSONPretty(res.httpCode, res, "\t")
}

func (res *Response) Marshal() []byte {
	buffer, _ := json.Marshal(res)
	return buffer
}

func (res *Response) Unmarshal(v interface{}) error {
	err := json.Unmarshal(res.Data, v)
	return err
}

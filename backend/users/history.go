package main

import (
	"encoding/json"
	"time"
)

type History struct {
	Id           string    `json:"id"`
	Source       string    `json:"source"`
	Time         time.Time `json:"time"`
	Type         string    `json:"type"`
	Method       string    `json:"method"`
	Request      string    `json:"request"`
	Response     string    `json:"response"`
	IsError      bool      `json:"is_error"`
	ErrorMessage string    `json:"error_message"`
}

func (h *History) Data() []byte {
	data, _ := json.Marshal(h)
	return data
}

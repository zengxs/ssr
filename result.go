package ssr

import (
	"encoding/json"
	"time"
)

type Result struct {
	ContextID  string        `json:"id"`
	Title      string        `json:"title,omitempty"`
	Meta       string        `json:"meta,omitempty"`
	Body       string        `json:"body"`
	RenderTime time.Duration `json:"-"`
}

func NewResultFromMap(obj map[string]interface{}) (*Result, error) {
	encoded, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	res := &Result{}
	err = json.Unmarshal(encoded, res)
	return res, err
}

package ssr

import (
	"encoding/json"
	"html/template"
	"reflect"
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

// ToHtmlSafeMap make all string field of Result to html safe (no escape in template)
// 转换所有字符串字段到 template.HTML 来确保在模版中不被转义
func (res *Result) toHtmlSafeMap() map[string]interface{} {
	result := make(map[string]interface{})

	t := reflect.TypeOf(res).Elem()
	v := reflect.ValueOf(res).Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		field := t.Field(i)
		value := v.FieldByName(field.Name).Interface()
		if field.Type.Kind() == reflect.String {
			result[field.Name] = template.HTML(value.(string))
		} else {
			result[field.Name] = value
		}
	}

	return result
}

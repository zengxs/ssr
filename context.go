package ssr

import (
	`encoding/base32`
	`encoding/json`
	`fmt`
	`net/http`

	`github.com/google/uuid`
)

type Context struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	QueryString string `json:"query_string"`
}

func NewContext(r *http.Request) *Context {
	ctx := &Context{
		Path:        r.URL.Path,
		QueryString: r.URL.RawQuery,
	}

	// generate context id: base32(uuid())
	data, err := uuid.Must(uuid.NewUUID()).MarshalBinary() // new uuid, convert to 16 bytes
	if err != nil {
		panic(fmt.Errorf("failed to marshal uuid to text: %v", err))
	}
	ctx.ID = base32.
		// base32.HexEncoding (lower case)
		NewEncoding("0123456789abcdefghijklmnopqrstuv").
		// strip padding characters
		WithPadding(base32.NoPadding).
		// should return string of 26 characters
		EncodeToString(data)

	return ctx
}

func (ctx *Context) ToMap() (map[string]interface{}, error) {
	data, err := json.Marshal(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(data, &result)
	return result, err
}

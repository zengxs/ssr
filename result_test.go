package ssr

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResultFromMap(t *testing.T) {
	obj1 := map[string]interface{}{
		"id":   "yxnrqagojqi55ciy7jqgndbbsq",
		"body": "hello",
	}

	if res1, err := NewResultFromMap(obj1); err != nil {
		t.Fatal(err)
	} else {
		assert.Equal(t, obj1["id"], res1.ContextID)
		assert.Equal(t, obj1["body"], res1.Body)
	}
}

func TestResult_ToHtmlSafeMap(t *testing.T) {
	tmpl := template.Must(template.New("").Parse("{{ .Body }}"))

	res1, _ := NewResultFromMap(map[string]interface{}{
		"id":   "yxnrqagojqi55ciy7jqgndbbsq",
		"body": "<h1>hello</h1>",
	})

	rr1_0 := new(bytes.Buffer)
	_ = tmpl.Execute(rr1_0, res1)
	assert.Equal(t, "&lt;h1&gt;hello&lt;/h1&gt;", rr1_0.String())

	rr1_1 := new(bytes.Buffer)
	obj1 := res1.toHtmlSafeMap()
	_ = tmpl.Execute(rr1_1, obj1)
	assert.Equal(t, "<h1>hello</h1>", rr1_1.String())
}

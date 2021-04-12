package ssr

import (
	`html/template`
	`io/ioutil`
	`net/http`
	`net/http/httptest`
	`testing`

	`github.com/stretchr/testify/assert`
)

func TestRenderer_ServeHTTP(t *testing.T) {
	tmpl := template.Must(template.ParseFiles("./testdata/test.gohtml"))
	bundleContent, err := ioutil.ReadFile("./testdata/render-pass.js")
	if err != nil {
		t.Error(err)
	}
	bundle := &ScriptBundle{
		Content:  bundleContent,
		FuncName: "render",
	}
	renderer := NewRenderer(bundle, tmpl)

	req, _ := http.NewRequest("GET", "/testing", nil)
	rr := httptest.NewRecorder()
	renderer.ServeHTTP(rr, req)

	assert.Contains(t, rr.Body.String(), "server")
}

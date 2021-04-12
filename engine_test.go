package ssr

import (
	`io/ioutil`
	`net/http`
	`testing`
	`time`

	`github.com/stretchr/testify/assert`
)

func TestJSEngine(t *testing.T) {
	bundleContent, err := ioutil.ReadFile("./testdata/render-pass.js")
	if err != nil {
		t.Error(err)
	}

	bundle := &ScriptBundle{
		Content:  bundleContent,
		FuncName: "render",
	}

	je, err := NewJSEngine(bundle, 2*time.Second)
	if err != nil {
		t.Error(err)
	}

	req1, _ := http.NewRequest("GET", "https://github.com/zengxs", nil)
	ctx1 := NewContext(req1)

	select {
	case result := <-je.Handle(ctx1):
		assert.Equal(t, "server", result.Body)
		assert.Equal(t, ctx1.ID, result.ContextID)
	case <-time.After(je.Timeout):
		t.Errorf("bundle script render timeout")
	}
}

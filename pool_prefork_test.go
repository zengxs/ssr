package ssr

import (
	`io/ioutil`
	`net/http`
	`testing`
	`time`

	`github.com/stretchr/testify/assert`
)

func TestPoolPrefork(t *testing.T) {
	bundleContent, err := ioutil.ReadFile("./testdata/render-pass.js")
	if err != nil {
		t.Error(err)
	}
	bundle := &ScriptBundle{
		Content:  bundleContent,
		FuncName: "render",
	}

	pool := NewPoolPrefork(bundle, 3*time.Second, 4)
	je1 := pool.Get()
	assert.NotNil(t, je1)

	req1, _ := http.NewRequest("GET", "https://github.com/zengxs", nil)
	ctx1 := NewContext(req1)

	select {
	case res1 := <-je1.Handle(ctx1):
		pool.Put(je1)
		assert.Equal(t, ctx1.ID, res1.ContextID)
		assert.Equal(t, "server", res1.Body)
	case <-time.After(je1.Timeout):
		t.Fatal("js-engine render timeout")
	}

	pool.Drop(pool.Get())
}

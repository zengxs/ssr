package ssr

import (
	`net/http`
	`testing`

	`github.com/stretchr/testify/assert`
)

func TestNewContext(t *testing.T) {
	req1, _ := http.NewRequest("GET", "https://github.com/zengxs/ssr.git", nil)
	ctx1 := NewContext(req1)
	assert.Equal(t, 26, len(ctx1.ID)) // len(ctx.id) == 26
	assert.Equal(t, "/zengxs/ssr.git", ctx1.Path)
	assert.Equal(t, "", ctx1.QueryString)

	req2, _ := http.NewRequest("GET", "https://github.com/zengxs/ssr/issues?q=is%3Aissue+is%3Aopen", nil)
	ctx2 := NewContext(req2)
	assert.Equal(t, "/zengxs/ssr/issues", ctx2.Path)
	assert.Equal(t, "q=is%3Aissue+is%3Aopen", ctx2.QueryString)
}

func TestContext_ToMap(t *testing.T) {
	req1, _ := http.NewRequest("GET", "https://github.com/", nil)
	if obj1, err := NewContext(req1).ToMap(); err != nil {
		t.Error(err)
	} else {
		assert.Len(t, obj1["id"], 26)
		assert.Equal(t, "/", obj1["path"])
		assert.Equal(t, "", obj1["query_string"])
	}
}

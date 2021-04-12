package ssr

import (
	`testing`

	`github.com/stretchr/testify/assert`
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

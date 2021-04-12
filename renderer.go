package ssr

import (
	`html/template`
	`net/http`
	`runtime`
	`time`
)

// Renderer implement http.Handler
type Renderer struct {
	// template to render html
	tmpl *template.Template

	// js-engine pool
	pool Pool

	// server render script bundle
	bundle *ScriptBundle
}

func NewRenderer(bundle *ScriptBundle, tmpl *template.Template) *Renderer {
	r := &Renderer{
		tmpl:   tmpl,
		pool:   NewPoolPrefork(bundle, 3*time.Second, uint(runtime.NumCPU())),
		bundle: bundle,
	}
	return r
}

func (r *Renderer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req)

	je := r.pool.Get()
	var res *Result
	select {
	case res = <-je.Handle(ctx):
		r.pool.Put(je)
	case <-time.After(je.Timeout):
		r.pool.Drop(je)
	}

	_ = r.tmpl.Execute(w, res)
}

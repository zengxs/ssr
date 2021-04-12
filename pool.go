package ssr

type Pool interface {
	// Get a js-engine from pool
	Get() *JSEngine

	// Put back js-engine to pool for re-use
	Put(je *JSEngine)

	// Drop js-engine if render failed
	Drop(je *JSEngine)
}

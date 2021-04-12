package ssr

import (
	`fmt`
	`time`
)

type PoolPrefork struct {
	size   uint
	ch     chan *JSEngine
	bundle *ScriptBundle
}

func NewPoolPrefork(bundle *ScriptBundle, timeout time.Duration, size uint) *PoolPrefork {
	pool := &PoolPrefork{
		size:   size,
		ch:     make(chan *JSEngine, size),
		bundle: bundle,
	}

	go func() {
		for i := uint(0); i < size; i++ {
			je, err := NewJSEngine(bundle, timeout)
			if err != nil {
				panic(fmt.Errorf("failed to create new js-engine: %v", err))
			}
			pool.ch <- je
		}
	}()

	return pool
}

func (p *PoolPrefork) Get() *JSEngine {
	return <-p.ch
}

func (p *PoolPrefork) Put(je *JSEngine) {
	p.ch <- je
}

func (p *PoolPrefork) Drop(je *JSEngine) {
	timeout := je.Timeout
	// release this js-engine
	je.Stop()
	je = nil
	// create a new one
	newJe, err := NewJSEngine(p.bundle, timeout)
	if err != nil {
		panic(fmt.Errorf("failed to create new js-engine: %v", err))
	}
	p.ch <- newJe
}

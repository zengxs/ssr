package ssr

import (
	crand `crypto/rand`
	`encoding/binary`
	`fmt`
	`math/rand`
	`time`

	`github.com/dop251/goja`
	`github.com/dop251/goja_nodejs/eventloop`
)

// Inject callback function with this name to js runtime
var callbackFuncName = "__go_ssr_callback__"

type JSEngine struct {
	*eventloop.EventLoop
	bundle *ScriptBundle
	ch     chan *Result
	fn     goja.Callable

	// Maximum render time
	Timeout time.Duration
}

func (je *JSEngine) Handle(ctx *Context) <-chan *Result {
	je.RunOnLoop(func(vm *goja.Runtime) {
		obj, err := ctx.ToMap()
		if err != nil {
			panic(fmt.Errorf("failed to convert ssr context to map: %v", err))
		}

		if _, err := je.fn(nil, vm.ToValue(obj), vm.ToValue(callbackFuncName)); err != nil {
			panic(fmt.Errorf("failed to execute bundle entry function: %v", err))
		}
	})
	return je.ch
}

func NewJSEngine(bundle *ScriptBundle, timeout time.Duration) (*JSEngine, error) {
	je := &JSEngine{
		EventLoop: eventloop.NewEventLoop(),
		ch:        make(chan *Result, 1),
		bundle:    bundle,
		Timeout:   timeout,
	}

	je.Start()
	je.initRuntimeRandSource()
	je.initRuntimeGlobalProcess()
	je.initRuntimeScriptBundle()
	je.initRuntimeExtractEntryFunc()
	je.initRuntimeInjectCallback()

	return je, nil
}

func (je *JSEngine) initRuntimeRandSource() {
	je.RunOnLoop(func(vm *goja.Runtime) {
		var seed int64
		if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
			panic(fmt.Errorf("failed to read random bytes: %v", err))
		}
		vm.SetRandSource(rand.New(rand.NewSource(seed)).Float64)
	})
}

func (je *JSEngine) initRuntimeGlobalProcess() {
	je.RunOnLoop(func(vm *goja.Runtime) {
		if err := vm.Set("process", map[string]interface{}{
			"env": map[string]string{
				"NODE_ENV": "production",
				"VUE_ENV":  "server",
			},
		}); err != nil {
			panic(fmt.Errorf("failed to inject globa.process to js runtime: %v", err))
		}
	})
}

func (je *JSEngine) initRuntimeScriptBundle() {
	je.RunOnLoop(func(vm *goja.Runtime) {
		if _, err := vm.RunString(string(je.bundle.Content)); err != nil {
			panic(fmt.Errorf("failed to execute bundle script: %v", err))
		}
	})
}

func (je *JSEngine) initRuntimeExtractEntryFunc() {
	je.RunOnLoop(func(vm *goja.Runtime) {
		if fn, ok := goja.AssertFunction(vm.Get(je.bundle.FuncName)); ok {
			je.fn = fn
		} else {
			panic("failed to extract bundle script entry func")
		}
	})
}

func (je *JSEngine) initRuntimeInjectCallback() {
	je.RunOnLoop(func(vm *goja.Runtime) {
		callback := func(call goja.FunctionCall) goja.Value {
			obj := call.Argument(0).Export().(map[string]interface{})
			res, err := NewResultFromMap(obj)
			if err != nil {
				panic(fmt.Errorf("failed to convert map to ssr result: %v", err))
			}
			je.ch <- res
			return nil
		}

		if err := vm.Set(callbackFuncName, callback); err != nil {
			panic(fmt.Errorf("failed to inject callback function: %v", err))
		}
	})
}

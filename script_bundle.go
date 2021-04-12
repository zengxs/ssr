package ssr

type ScriptBundle struct {
	// js script content
	Content []byte

	// js bundle script entry function name
	FuncName string
}

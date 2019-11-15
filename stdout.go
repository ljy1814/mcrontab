package main

import (
	"context"
	"fmt"
	"os"
)

var (
	_defaultStdout = NewStdout()
	defaultPattern = "%L %d-%T %f %M"
)

type StdoutHandler struct {
	render Render
}

func NewStdout() *StdoutHandler {
	return &StdoutHandler{
		render: newPatternRender(defaultPattern),
	}
}

func (h *StdoutHandler) Log(ctx context.Context, lv Level, args ...interface{}) {
	// handle data
	//d[_time] = time.Now().Format(_timeFormat)
	//h.render.Render(os.Stderr, args)
	fmt.Fprint(os.Stderr, args)
	os.Stderr.Write([]byte("\n"))
}

func (h *StdoutHandler) Close() error {
	return nil
}

func (h *StdoutHandler) SetFormat(f string) {
	h.render = newPatternRender(f)
}

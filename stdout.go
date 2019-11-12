package main

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

func (h *StdoutHandler) Log() {

}

func (h *StdoutHandler) Close() {
	return nil
}

func (h *StdoutHandler) SetFormat(f string) {
	h.render = newPatternRender(f)
}

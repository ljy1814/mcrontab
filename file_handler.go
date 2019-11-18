package main

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
)

const (
	_infoIdx = iota
	_warnIdx
	_errorIdx
	_totalIdx
)

var (
	_fileNames = map[int]string{
		_infoIdx:  "info.log",
		_warnIdx:  "warn.log",
		_errorIdx: "error.log",
	}
)

type FileHandler struct {
	// 渲染
	render Render
	fws    []*FileWriter
}

func NewFile(dir string, buffersize, rotateSize int64, maxLogFile int) *FileHandler {
	newWriter := func(name string) *FileWriter {
		var options []Option
		if rotateSize > 0 {
			options = append(options, MaxSize(rotateSize))
		}
		if maxLogFile > 0 {
			options = append(options, MaxFile(maxLogFile))
		}

		w, err := New(filepath.Join(dir, name), options...)
		if err != nil {
			panic(err)
		}
		return w
	}

	handler := &FileHandler{
		//render: newPatternRender("[%D %T] [%L] [%S] %M"),
	}

	for idx, name := range _fileNames {
		handler.fws[idx] = newWriter(name)
	}

	return handler
}

func (fh *FileHandler) Close() error {
	return nil
}

func (fh *FileHandler) SetFormat(string) {

}

func (h *FileHandler) Log(ctx context.Context, lv Level, args ...interface{}) {

	var w io.Writer

	switch lv {
	case _warnLevel:
		w = h.fws[_warnIdx]
	case _errorLevel:
		w = h.fws[_errorIdx]
	default:
		w = h.fws[_infoIdx]
	}

	bs := []byte(fmt.Sprint(args))

	//logrus.Infof("************FileHandler[fws:%d] data:%s", len(h.fws), bs)

	w.Write(bs)
	return
}

func (h *FileHandler) Write(p []byte) (n int, err error) {
	var w io.Writer

	w = h.fws[_infoIdx]

	w.Write(p)
	return
}

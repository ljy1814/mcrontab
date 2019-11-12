package main

import "path/filepath"

const (
	_infoIdx = iota
	_warnIdx
	_errorIdx
	_totalIdx
)

var (
	_fileNames = map[int]string{}
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
		render: newPatternRender("[%D %T] [%L] [%S] %M"),
	}

	for idx, name := range _fileNames {
		handler.fws[idx] = newWriter(name)
	}

	return handler
}

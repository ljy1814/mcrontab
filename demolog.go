package main

import "path/filepath"

func NewDemoFile(dir string, buffersize, rotateSize int64, maxLogFile int) *FileHandler {
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

	handler := &FileHandler{}
	handler.fws = make([]*FileWriter, len(_fileNames))
	for idx, name := range _fileNames {
		handler.fws[idx] = newWriter(name)
	}

	return handler
}

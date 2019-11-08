package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"
)

type wrapFile struct {
	fsize int64
	fp    *os.File
}

func (w *wrapFile) write(p []byte) (n int, err error) {
	n, err = w.fp.Write(p)
	w.fsize += int64(n)

	return
}

func (w *wrapFile) size() int64 {
	return w.fsize
}

func newWrapFile(fpath string) (*wrapFile, error) {
	fp, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	return &wrapFile{
		fp:    fp,
		fsize: fi.Size(),
	}, err
}

type FileWriter struct {
	dir    string
	fname  string
	ch     chan *bytes.Buffer
	stdlog *log.Logger
	pool   *sync.Pool

	lastRotateFormat string
	lastSplitNum     int

	current *wrapFile
	files   *list.List

	closed int32
	wg     sync.WaitGroup
	opt    option
}

type option struct {
	RotateFormat string
	MaxFile      int
	MaxSize      int64
	ChanSize     int

	// 检查时间间隔
	RotateInterval time.Duration
	WriteTimeout   time.Duration
}

func (f *FileWriter) Write(p []byte) (int, error) {
	if atomic.LoadInt32(&f.closed) == 1 {
		f.stdlog.Printf("%s", p)
		return 0, fmt.Errorf("filewriter already closed")
	}

	buf := f.getBuf()
	buf.Write(p)

	// timeout 写入量太大
	if f.opt.WriteTimeout == 0 {
		select {
		case f.ch <- buf:
			return len(p), nil
		default:
			return 0, fmt.Errorf("log channel is full, discard log")
		}
	}

	timeout := time.NewTimer(f.opt.WriteTimeout)
	select {
	case f.ch <- buf:
		return len(p), nil
	case <-timeout.C:
		return 0, fmt.Errorf("log channel is full, discard log")
	}
}

func (f *FileWriter) Close() error {
	atomic.StoreInt32(&f.closed, 1)
	close(f.ch)
	f.wg.Wait()

	return nil
}

func (f *FileWriter) write(p []byte) error {
	if f.current == nil {
		f.stdlog.Printf("can't write log to file, please check stderr log for detail")
		f.stdlog.Printf("%s", p)
	}

	_, err := f.current.write(p)
	return err
}

func (f *FileWriter) getBuf() *bytes.Buffer {
	return f.pool.Get().(*bytes.Buffer)
}

func (f *FileWriter) putBuf(buf *bytes.Buffer) {
	buf.Reset()
	f.pool.Put(buf)
}

func (f *FileWriter) daemon() {
	aggsbuf := &bytes.Buffer{}
	tk := time.NewTicker(f.opt.RotateInterval)

	aggstk := time.NewTicker(10 * time.Millisecond)
	var err error

	for {
		select {
		case t := <-tk.C:
			f.checkRotate(t)
		case buf, ok := <-f.ch:
			if ok {
				aggsbuf.Write(buf.Bytes())
				f.putBuf(buf)
			}
		case <-aggstk.C:
			if aggsbuf.Len() > 0 {
				if err = f.write(aggsbuf.Bytes()); err != nil {
					f.stdlog.Printf("write log error:%v", err)
				}
				aggsbuf.Reset()
			}
		}

		if atomic.LoadInt32(&f.closed) != 1 {
			continue
		}

		// 处理尾巴
		if err = f.write(aggsbuf.Bytes()); err != nil {
			f.stdlog.Printf("write log error:%v", err)
		}

		for buf := range f.ch {
			if err = f.write(buf.Bytes()); err != nil {
				f.stdlog.Printf("write log error:%v", err)
			}
			f.putBuf(buf)
		}
		break
	}

	f.wg.Done()
}

func (f *FileWriter) checkRotate(t time.Time) {
	formatFname := func(format string, num int) string {
		if num == 0 {
			return fmt.Sprintf("%s.%s", f.fname, format)
		}
		return fmt.Sprintf("%s.%s.%03d", f.fname, format, num)
	}

	format := t.Format(f.opt.RotateFormat)

	// 有文件个数线上
	if f.opt.MaxFile != 0 {
		for f.files.Len() > f.opt.MaxFile {
			rt := f.files.Remove(f.files.Front()).(rotateItem)
			fpath := filepath.Join(f.dir, rt.fname)
			if err := os.Remove(fpath); err != nil {
				f.stdlog.Printf("remove file:%s err:%v", fpath, err)
			}
		}
	}

	if format != f.lastRotateFormat || (f.opt.MaxSize != 0 && f.current.size() > f.opt.MaxSize) {
		var err error
		if err = f.current.fp.Close(); err != nil {
			f.stdlog.Printf("close current file error:%s", err)
		}

		fname := formatFname(f.lastRotateFormat, f.lastSplitNum)
		oldpath := filepath.Join(f.dir, f.fname)
		newpath := filepath.Join(f.dir, fname)

		if err = os.Rename(oldpath, newpath); err != nil {
			f.stdlog.Printf("rename file %s to %s err:%v", oldpath, newpath, err)
			return
		}

		f.files.PushBack(rotateItem{
			fname: fname,
		})

		if format != f.lastRotateFormat {
			f.lastRotateFormat = format
			f.lastSplitNum = 0
		} else {
			f.lastSplitNum++
		}

		f.current, err = newWrapFile(filepath.Join(f.dir, f.fname))
		if err != nil {
			f.stdlog.Printf("create log file err:%v", err)
		}
	}
}

type rotateItem struct {
	rotateTime int64
	rotateNum  int
	fname      string
}

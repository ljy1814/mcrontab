package main

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Sirupsen/logrus"
)

// 日志写入器,代理写入器
type wrapFile struct {
	fsize int64
	fp    *os.File
}

// 日志写入
func (w *wrapFile) write(p []byte) (n int, err error) {
	logrus.Infof("wrapFile:%s", p)
	n, err = w.fp.Write(p)
	w.fsize += int64(n)

	return
}

func (w *wrapFile) size() int64 {
	return w.fsize
}

// 使用已有的serv.log或者新建一个serv.log
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
	dir   string
	fname string
	// buffer 队列
	ch     chan *bytes.Buffer
	stdlog *log.Logger
	// 缓冲池
	pool *sync.Pool

	lastRotateFormat string
	lastSplitNum     int

	// 当前写入文件
	current *wrapFile
	// 历史滚动文件
	files *list.List

	closed int32
	wg     sync.WaitGroup
	opt    option
}

type option struct {
	RotateFormat string
	MaxFile      int
	MaxSize      int64
	ChanSize     int

	// 日志滚动检查时间间隔
	RotateInterval time.Duration
	WriteTimeout   time.Duration
}

// 日志写入,并非真的写入而是写入缓存
func (f *FileWriter) Write(p []byte) (int, error) {
	if atomic.LoadInt32(&f.closed) == 1 {
		f.stdlog.Printf("%s", p)
		return 0, fmt.Errorf("filewriter already closed")
	}

	// 从内存池里面取一个buffer
	buf := f.getBuf()
	// 写入buffer
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

// 后台写入
func (f *FileWriter) daemon() {
	aggsbuf := &bytes.Buffer{}
	tk := time.NewTicker(f.opt.RotateInterval)

	aggstk := time.NewTicker(10 * time.Millisecond)
	var err error

	for {
		select {
		// 检查是否需要滚动日志文件
		case t := <-tk.C:
			f.checkRotate(t)
		case buf, ok := <-f.ch:
			// 写入到buffer
			if ok {
				aggsbuf.Write(buf.Bytes())
				f.putBuf(buf)
			}
		case <-aggstk.C:
			// 时间片到了写入日志文件
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

// 滚动检查
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
		logrus.Errorf("close current file error:%v", err)

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

type Option func(opt *option)

var (
	RotateDaily   = "2006-01-02-15"
	defaultOption = option{
		RotateFormat:   RotateDaily,
		MaxSize:        1 << 30,
		ChanSize:       1024 * 8,
		RotateInterval: 10 * time.Second,
	}
)

func RotateFormat(f string) Option {
	if strings.Contains(f, ".") {
		panic(fmt.Sprintf("rotate format can not contain '.' format:%s", f))
	}

	return func(o *option) {
		o.RotateFormat = f
	}
}

func MaxFile(n int) Option {
	return func(o *option) {
		o.MaxFile = n
	}
}

func MaxSize(n int64) Option {
	return func(o *option) {
		o.MaxSize = n
	}
}

func ChanSize(n int) Option {
	return func(o *option) {
		o.ChanSize = n
	}
}

// 新建日志写入器
func New(fpath string, fns ...Option) (*FileWriter, error) {
	opt := defaultOption
	// 处理额外配置系列
	for _, fn := range fns {
		fn(&opt)
	}

	fname := filepath.Base(fpath)
	if fname == "" {
		return nil, fmt.Errorf("filename cannot empty")
	}

	dir := filepath.Dir(fpath)
	fi, err := os.Stat(dir)
	if err == nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s already exists and not a directory", dir)
	}

	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0744); err != nil {
			return nil, fmt.Errorf("create directory:%s err:%v", dir, err)
		}
	}

	current, err := newWrapFile(fpath)
	if err != nil {
		return nil, err
	}

	stdlog := log.New(os.Stderr, "flog ", log.LstdFlags)
	ch := make(chan *bytes.Buffer, opt.ChanSize)

	files, err := parseRotateItem(dir, fname, opt.RotateFormat)
	if err != nil {
		files = list.New()
		stdlog.Printf("parseRotateItem err:%s", err)
	}

	lastRotateFormat := time.Now().Format(opt.RotateFormat)
	var splitNum int
	if files.Len() > 0 {
		rt := files.Front().Value.(rotateItem)
		if strings.Contains(rt.fname, lastRotateFormat) {
			splitNum = rt.rotateNum
		}
	}

	fw := &FileWriter{
		opt:    opt,
		dir:    dir,
		fname:  fname,
		stdlog: stdlog,
		ch:     ch,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		lastSplitNum:     splitNum,
		lastRotateFormat: lastRotateFormat,
		files:            files,
		current:          current,
	}

	fw.wg.Add(1)
	go fw.daemon()

	return fw, nil
}

func parseRotateItem(dir, fname, rotateFormat string) (*list.List, error) {
	// 读取日志目录
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	parse := func(s string) (rt rotateItem, err error) {
		rt.fname = s
		// error.log.2019-10-31.001 -> 2019-10-31.001
		s = strings.TrimLeft(s[len(fname):], ".")
		seqs := strings.Split(s, ".")
		var t time.Time
		switch len(seqs) {
		case 2:
			if rt.rotateNum, err = strconv.Atoi(seqs[1]); err != nil {
				return
			}
			fallthrough
		case 1:
			if t, err = time.Parse(rotateFormat, seqs[0]); err != nil {
				return
			}
			rt.rotateTime = t.Unix()
		}
		return
	}

	var items []rotateItem

	// 处理已存在的日文件
	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), fname) && fi.Name() != fname {
			rt, err := parse(fi.Name())
			if err != nil {
				// TODO error handle
				continue
			}
			items = append(items, rt)
		}
	}

	// 排序已存在的日志文件
	sort.Slice(items, func(i, j int) bool {
		if items[i].rotateTime == items[j].rotateTime {
			return items[i].rotateNum > items[j].rotateNum
		}
		return items[i].rotateTime > items[j].rotateTime
	})

	l := list.New()

	// 放入链表
	for _, item := range items {
		l.PushBack(item)
	}
	return l, nil
}

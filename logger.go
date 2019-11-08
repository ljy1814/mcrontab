package main

import (
	"context"
	"time"

	pkgerr "github.com/pkg/errors"
)

type Config struct {
	Family         string
	Host           string
	Stdout         bool
	Dir            string
	FileBufferSize int64
	MaxLogFile     int
	RotateSize     int64
	V              int32
	Module         map[string]int32
	Filter         []string
}

var (
	_stdout   bool
	_v        int
	_dir      string
	_agentDSN string
)

func InitLogger(conf *Config) {
	var isNil bool
	if conf == nil {
		isNil = true
		conf = &Config{
			Stdout: _stdout,
			Dir:    _dir,
			V:      int32(_v),
		}
	}

	if isNil {

	}

	//var hs []Handler
	//if conf.Stdout {
	//	hs = append(hs, NewStdout())
	//}
}

const (
	_debugLevel Level = iota
	_infoLevel
	_warnLevel
	_errorLevel
	_fatalLevel
)

var (
	levelNames = [...]string{
		_debugLevel: "DEBUG",
		_infoLevel:  "INFO",
		_warnLevel:  "WARN",
		_errorLevel: "ERROR",
		_fatalLevel: "FATAL",
	}
)

type Level int

func (l Level) String() string {
	return levelNames[l]
}

type FieldType int32

type ArrayEncoder struct {
}

type ArrayMarshaler interface {
	MarshalLogArray(ArrayEncoder) error
}

type ObjectEncoder interface {
	AddArray(key string)
}

type Field struct {
	Key       string
	Value     interface{}
	Type      FieldType
	StringVal string
	Int64Val  int64
}

const (
	UnknownType FieldType = iota
	StringType
	IntType
	Int64Type
	UintType
	Uint64Type
	Float32Type
	Float64Type
	DurationType
)

func (f Field) AddTo(enc ObjectEncoder) {

}

type D Field

func KV(k string, v interface{}) D {
	return D{
		Key:   k,
		Value: v,
	}
}

func KVString(k string, v string) D {
	return D{
		Key:       k,
		Type:      StringType,
		StringVal: v,
	}
}

func KVInt64(k string, v int64) D {
	return D{
		Key:      k,
		Type:     Int64Type,
		Int64Val: v,
	}
}

type Handler interface {
	Log(context.Context, Level, ...D)
	Close() error
	SetFormat(string)
}

const (
	_source     = "source"
	_level      = "level"
	_time       = "time"
	_levelValue = "level_value"
)

type Handlers struct {
	filters  map[string]struct{}
	handlers []Handler
}

func (hs Handlers) Log(ctx context.Context, lv Level, d ...D) {
	hasSource := false
	for i := range d {
		if _, ok := hs.filters[d[i].Key]; ok {
			d[i].Value = "****"
		}

		if d[i].Key == _source {
			hasSource = true
		}
	}

	if !hasSource {
		fn := funcName(3)
		errIncr(lv, fn)
		d = append(d, KVString(_source, fn))
	}

	d = append(d, KV(_time, time.Now()), KVInt64(_levelValue, int64(lv)), KVString(_level, lv.String()))

	for _, h := range hs.handlers {
		h.Log(ctx, lv, d...)
	}
}

func (hs Handlers) Close() (err error) {
	for _, h := range hs.handlers {
		if e := h.Close(); e != nil {
			err = pkgerr.WithStack(e)
		}
	}

	return
}

func (hs Handlers) SetFormat(f string) {
	for _, h := range hs.handlers {
		h.SetFormat(f)
	}
}

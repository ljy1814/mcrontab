package main

import (
	"bytes"
	"fmt"
	"io"
	"path"
	"strings"
	"sync"
	"time"
)

type pattern struct {
	funcs   []func(map[string]interface{}) string
	bufPool sync.Pool
}

var (
	patternMap = map[string]func(map[string]interface{}) string{
		"T": longTime,
		"t": shortTime,
		"D": longDate,
		"d": shortDate,
		"L": keyFactory(_level),
		"f": keyFactory(_source),
		"i": keyFactory(_instanceID),
		"e": keyFactory(_deplyEnv),
		"z": keyFactory(_zone),
		"S": longSource,
		"s": shortSource,
		"M": message,
	}
)

func newPatternRender(format string) *pattern {
	p := &pattern{
		bufPool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
	b := make([]byte, 0, len(format))

	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			b = append(b, format[i])
			continue
		}
		if i+1 > len(format) {
			b = append(b, format[i])
			continue
		}

		f, ok := patternMap[string(format[i+1])]
		if !ok {
			b = append(b, format[i])
			continue
		}

		if len(b) != 0 {
			p.funcs = append(p.funcs, textFactory(string(b)))
			b = b[:0]
		}

		p.funcs = append(p.funcs, f)
		i++
	}
	return p
}

func (p *pattern) Render(w io.Writer, d map[string]interface{}) error {
	buf := p.bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		p.bufPool.Put(buf)
	}()

	for _, f := range p.funcs {
		buf.WriteString(f(d))
	}

	_, err := buf.WriteTo(w)
	return err
}

func (p *pattern) RenderString(d map[string]interface{}) string {
	buf := p.bufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		p.bufPool.Put(buf)
	}()

	for _, f := range p.funcs {
		buf.WriteString(f(d))
	}

	return buf.String()
}

func textFactory(text string) func(map[string]interface{}) string {
	return func(d map[string]interface{}) string {
		return text
	}
}

func keyFactory(key string) func(map[string]interface{}) string {
	return func(d map[string]interface{}) string {
		v, ok := d[key]
		if !ok {
			return ""
		}

		s, ok := v.(string)
		if ok {
			return s
		}
		return fmt.Sprint(v)
	}
}

func longSource(d map[string]interface{}) string {
	fn, ok := d[_source].(string)
	if ok {
		return fn
	}
	return "unknown:0"
}

func shortSource(d map[string]interface{}) string {
	fn, ok := d[_source].(string)
	if ok {
		return path.Base(fn)
	}
	return "unknown:0"
}

func longTime(d map[string]interface{}) string {
	return time.Now().Format("15:04:05.000")
}

func shortTime(d map[string]interface{}) string {
	return time.Now().Format("15:04")
}

func longDate(d map[string]interface{}) string {
	return time.Now().Format("2006/01/02")
}

func shortDate(d map[string]interface{}) string {
	return time.Now().Format("01/02")
}

func isIntervalKey(k string) bool {
	switch k {
	case _level, _levelValue, _time, _source, _instanceID, _appID, _deplyEnv, _zone:
		return true
	default:
		return false
	}
}

func message(d map[string]interface{}) string {
	var m string
	var s []string

	for k, v := range d {
		if k == _log {
			m = fmt.Sprint(v)
			continue
		}

		if isIntervalKey(k) {
			continue
		}

		s = append(s, fmt.Sprintf("%s:%v", k, v))
	}

	s = append(s, m)
	return strings.Join(s, " ")
}

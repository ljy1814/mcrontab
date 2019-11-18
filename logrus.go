package main

import (
	"github.com/Sirupsen/logrus"
)

// 对特定类型日志做特殊处理
// 比如发邮件,发短信,钉钉通知等等
type logrusHook struct {
}

func (logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (lh logrusHook) Fire(entry *logrus.Entry) error {
	lv := _infoLevel
	switch entry.Level {
	case logrus.PanicLevel:
	case logrus.FatalLevel:
	case logrus.ErrorLevel:
		lv = _errorLevel
	case logrus.WarnLevel:
		lv = _warnLevel
	case logrus.InfoLevel:
		lv = _infoLevel
	case logrus.DebugLevel:
	case logrus.TraceLevel:
	}

	if lv == _infoLevel {
	}

	//buf := bytes.NewBufferString("")
	//buf.WriteString(entry.Message)

	//for k, v := range entry.Data {
	//	buf.WriteString(" ")
	//	buf.WriteString(k)
	//	buf.WriteString(":")
	//	buf.WriteString(fmt.Sprint(v))
	//	buf.WriteString(" ")
	//}

	//GLogger.Log(context.Background(), lv, buf.String())
	return nil
}

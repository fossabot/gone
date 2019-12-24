package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"runtime"
	"sync/atomic"
)

type (
	Fields = logrus.Fields
	Level  = logrus.Level
)

type wrapperLogger struct {
	*logrus.Logger
	mux             logrus.MutexWrap
	depth           int
	fileLineNoLevel Level
}

func NewLog() (logger *wrapperLogger) {
	logger = &wrapperLogger{
		Logger:          logrus.New(),
		depth:           0,
		fileLineNoLevel: logrus.WarnLevel,
	}
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(NewFormatter())
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "log/foo.log",
		MaxSize:    500, // 500M
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})
	return
}

func (w *wrapperLogger) SetFileLineNoLevel(l Level) {
	atomic.StoreUint32((*uint32)(&w.fileLineNoLevel), uint32(l))
}

func (w *wrapperLogger) getFileAndLineno() string {
	fl := "unknown:0"
	_, file, line, ok := runtime.Caller(w.depth + 3)
	if ok {
		file = filepath.Base(file)
		fl = fmt.Sprintf("%s:%d", file, line)
	}
	return fl
}

func (w *wrapperLogger) defaultFields(level Level) Fields {
	var fields = Fields{}
	if w.fileLineNoLevel >= level {
		fields[fileLinenoStr] = w.getFileAndLineno()
	} else {
		fields[fileLinenoStr] = "unknown:0"
	}
	fields[hostnameStr] = getHostname()
	w.WithField("", fields)
	return fields
}

func (w *wrapperLogger) Info(args ...interface{}) {
	w.WithFields(w.defaultFields(logrus.InfoLevel)).Info(args...)
}

func (w *wrapperLogger) Warn(args ...interface{}) {
	w.WithFields(w.defaultFields(logrus.WarnLevel)).Info(args...)
}

func (w *wrapperLogger) Panic(args ...interface{}) {
	w.WithFields(w.defaultFields(logrus.PanicLevel)).Info(args...)
}

func (w *wrapperLogger) Error(args ...interface{}) {
	w.WithFields(w.defaultFields(logrus.ErrorLevel)).Info(args...)
}

func (w *wrapperLogger) Debug(args ...interface{}) {
	w.WithFields(w.defaultFields(logrus.DebugLevel)).Info(args...)
}

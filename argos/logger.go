package argos

import (
	"context"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (ll *Logger) SetOutput(w io.Writer) {
	ll.logger.SetOutput(w)
}

func (ll *Logger) SetLevel(lv klog.Level) {
}

func (ll *Logger) Fatal(v ...interface{}) {
	ll.logger.Fatal(v...)
}

func (ll *Logger) Error(v ...interface{}) {
	ll.logger.Error(v...)
}

func (ll *Logger) Warn(v ...interface{}) {
	ll.logger.Warn(v...)
}

func (ll *Logger) Notice(v ...interface{}) {
	ll.logger.Info(v...)
}

func (ll *Logger) Info(v ...interface{}) {
	ll.logger.Info(v...)
}

func (ll *Logger) Debug(v ...interface{}) {
	ll.logger.Debug(v...)
}

func (ll *Logger) Trace(v ...interface{}) {
	ll.logger.Trace(v...)
}

func (ll *Logger) Fatalf(format string, v ...interface{}) {
	ll.logger.Fatalf(format, v...)
}

func (ll *Logger) Errorf(format string, v ...interface{}) {
	ll.logger.Errorf(format, v...)
}

func (ll *Logger) Warnf(format string, v ...interface{}) {
	ll.logger.Warnf(format, v...)
}

func (ll *Logger) Noticef(format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *Logger) Infof(format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *Logger) Debugf(format string, v ...interface{}) {
	ll.logger.Debugf(format, v...)
}

func (ll *Logger) Tracef(format string, v ...interface{}) {
	ll.logger.Tracef(format, v...)
}

func (ll *Logger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Fatalf(format, v...)
}

func (ll *Logger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Errorf(format, v...)
}

func (ll *Logger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Warnf(format, v...)
}

func (ll *Logger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *Logger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *Logger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Debugf(format, v...)
}

func (ll *Logger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Tracef(format, v...)
}

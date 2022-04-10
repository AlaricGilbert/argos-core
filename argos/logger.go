package argos

import (
	"context"
	"io"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/sirupsen/logrus"
)

var logger = logrus.StandardLogger()

// WarppedLogger is a warp of logrus.Logger, making it can be used as a kitex/klog logger
type WarppedLogger struct {
	logger *logrus.Logger
}

func SetLogger(l *logrus.Logger) {
	klog.SetLogger(WarpLogger(l))
	logger = l
}

func StandardLogger() *logrus.Logger {
	return logger
}

func WarpLogger(logger *logrus.Logger) *WarppedLogger {
	return &WarppedLogger{
		logger: logger,
	}
}

func (ll *WarppedLogger) SetOutput(w io.Writer) {
	ll.logger.SetOutput(w)
}

func (ll *WarppedLogger) SetLevel(lv klog.Level) {
}

func (ll *WarppedLogger) Fatal(v ...interface{}) {
	ll.logger.Fatal(v...)
}

func (ll *WarppedLogger) Error(v ...interface{}) {
	ll.logger.Error(v...)
}

func (ll *WarppedLogger) Warn(v ...interface{}) {
	ll.logger.Warn(v...)
}

func (ll *WarppedLogger) Notice(v ...interface{}) {
	ll.logger.Info(v...)
}

func (ll *WarppedLogger) Info(v ...interface{}) {
	ll.logger.Info(v...)
}

func (ll *WarppedLogger) Debug(v ...interface{}) {
	ll.logger.Debug(v...)
}

func (ll *WarppedLogger) Trace(v ...interface{}) {
	ll.logger.Trace(v...)
}

func (ll *WarppedLogger) Fatalf(format string, v ...interface{}) {
	ll.logger.Fatalf(format, v...)
}

func (ll *WarppedLogger) Errorf(format string, v ...interface{}) {
	ll.logger.Errorf(format, v...)
}

func (ll *WarppedLogger) Warnf(format string, v ...interface{}) {
	ll.logger.Warnf(format, v...)
}

func (ll *WarppedLogger) Noticef(format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *WarppedLogger) Infof(format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *WarppedLogger) Debugf(format string, v ...interface{}) {
	ll.logger.Debugf(format, v...)
}

func (ll *WarppedLogger) Tracef(format string, v ...interface{}) {
	ll.logger.Tracef(format, v...)
}

func (ll *WarppedLogger) CtxFatalf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Fatalf(format, v...)
}

func (ll *WarppedLogger) CtxErrorf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Errorf(format, v...)
}

func (ll *WarppedLogger) CtxWarnf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Warnf(format, v...)
}

func (ll *WarppedLogger) CtxNoticef(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *WarppedLogger) CtxInfof(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Infof(format, v...)
}

func (ll *WarppedLogger) CtxDebugf(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Debugf(format, v...)
}

func (ll *WarppedLogger) CtxTracef(ctx context.Context, format string, v ...interface{}) {
	ll.logger.Tracef(format, v...)
}

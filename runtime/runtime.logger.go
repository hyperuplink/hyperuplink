package runtime

import (
	"log/slog"
)

type AsyncLogger struct {
	logger *slog.Logger
}

func NewAsyncLogger(logger *slog.Logger) (l AsyncLogger) {
	l = AsyncLogger{}

	l.logger = logger

	return l
}

func (l AsyncLogger) fixArgs(args ...interface{}) []interface{} {
	if len(args)%2 != 0 {
		args = append([]interface{}{"args"}, args...)
	}
	return args
}

func (l AsyncLogger) Debug(args ...interface{}) {
	l.logger.Debug("Worker.Async", l.fixArgs(args)...)
}

func (l AsyncLogger) Info(args ...interface{}) {
	l.logger.Info("Worker.Async", l.fixArgs(args)...)
}

func (l AsyncLogger) Warn(args ...interface{}) {
	l.logger.Warn("Worker.Async", l.fixArgs(args)...)
}

func (l AsyncLogger) Error(args ...interface{}) {
	l.logger.Error("Worker.Async", l.fixArgs(args)...)
}

func (l AsyncLogger) Fatal(args ...interface{}) {
	l.logger.Error("Worker.Async", l.fixArgs(args)...)
}

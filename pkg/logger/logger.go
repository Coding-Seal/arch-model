package logger

import "log/slog"

type Logger struct {
	serviceName string
	log         *slog.Logger
}

func New(log *slog.Logger, serviceName string) *Logger {
	return &Logger{log: log, serviceName: serviceName}
}

func (l *Logger) SetServiceName(str string) {
	l.serviceName = str
}

func (l *Logger) Info(msg string, args ...slog.Attr) {
	l.log.Info(msg, l.addAttr(args...)...)
}

func (l *Logger) Error(msg string, args ...slog.Attr) {
	l.log.Error(msg, l.addAttr(args...)...)
}

func (l *Logger) Debug(msg string, args ...slog.Attr) {
	l.log.Debug(msg, l.addAttr(args...)...)
}

func (l *Logger) Warning(msg string, args ...slog.Attr) {
	l.log.Warn(msg, l.addAttr(args...)...)
}

func (l *Logger) addAttr(args ...slog.Attr) []any {
	allArgs := make([]any, 0, len(args)+1)
	allArgs = append(allArgs, slog.String("service", l.serviceName))

	for _, attr := range args {
		allArgs = append(allArgs, attr)
	}

	return allArgs
}

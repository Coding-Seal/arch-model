package logger

import "log/slog"

type Logger struct {
	serviceName string
	log         *slog.Logger
}

func New(log *slog.Logger, serviceName string) *Logger {
	return &Logger{log: log, serviceName: serviceName}
}

func (l *Logger) Info(msg string, args ...slog.Attr) {
	l.log.Info(msg, slog.String("service", l.serviceName))
}

func (l *Logger) Error(msg string, args ...slog.Attr) {
	l.log.Error(msg, slog.String("service", l.serviceName))
}

func (l *Logger) Debug(msg string, args ...slog.Attr) {
	l.log.Debug(msg, slog.String("service", l.serviceName))
}
func (l *Logger) Warning(msg string, args ...slog.Attr) {
	l.log.Warn(msg, slog.String("service", l.serviceName))
}

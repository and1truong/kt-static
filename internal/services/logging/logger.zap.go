package logging

import "go.uber.org/zap"

// zapLogger is an implementation of the Logger interface using zap.SugaredLogger.
type zapLogger struct {
	*zap.SugaredLogger
}

// NewLogger creates a new Logger instance using a production-ready zap configuration.
func NewLogger() Logger {
	// Use zap.NewProduction() for a robust, high-performance logger.
	logger, _ := zap.NewProduction()
	// AddCallerSkip(1) is necessary because the zapLogger methods (Debug, Info, etc.)
	// are wrappers around the SugaredLogger methods, and we want the caller of the
	// zapLogger method, not the zapLogger method itself.
	return &zapLogger{logger.WithOptions(zap.AddCallerSkip(1)).Sugar()}
}

// Debug implements the Logger interface's Debug method.
func (l *zapLogger) Debug(msg string, keysAndValues ...any) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

// Info implements the Logger interface's Info method.
func (l *zapLogger) Info(msg string, keysAndValues ...any) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

// Warn implements the Logger interface's Warn method.
func (l *zapLogger) Warn(msg string, keysAndValues ...any) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

// Error implements the Logger interface's Error method.
func (l *zapLogger) Error(msg string, keysAndValues ...any) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

// With implements the Logger interface's With method.
func (l *zapLogger) With(keysAndValues ...any) Logger {
	// The With method of SugaredLogger returns a new SugaredLogger with the added context.
	return &zapLogger{l.SugaredLogger.With(keysAndValues...)}
}

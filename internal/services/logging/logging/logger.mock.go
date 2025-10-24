package logging

import "htruong/kt-crawler/internal/services/logging"

// MockLogger implements the logging.Logger interface for testing.
type MockLogger struct{}

func (m *MockLogger) Debug(msg string, keysAndValues ...any) {}
func (m *MockLogger) Info(msg string, keysAndValues ...any)  {}
func (m *MockLogger) Warn(msg string, keysAndValues ...any)  {}
func (m *MockLogger) Error(msg string, keysAndValues ...any) {}
func (m *MockLogger) With(keysAndValues ...any) logging.Logger {
	return m
}

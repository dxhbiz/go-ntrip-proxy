package log

var logger Logger

// Logger is a logger interface
type Logger interface {
	Sync() error

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Sync
func Sync() {
	if logger == nil {
		return
	}

	logger.Sync()
}

// Debug
func Debug(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Debug(args...)
}

// Debugf
func Debugf(format string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Debugf(format, args...)
}

// Info
func Info(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Info(args...)
}

// Infof
func Infof(format string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Infof(format, args...)
}

// Warn
func Warn(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Warn(args...)
}

// Warnf
func Warnf(format string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Warnf(format, args...)
}

// Error
func Error(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Error(args...)
}

// Errorf
func Errorf(format string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Errorf(format, args...)
}

// Fatal
func Fatal(args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Fatal(args...)
}

// Fatalf
func Fatalf(format string, args ...interface{}) {
	if logger == nil {
		return
	}

	logger.Fatalf(format, args...)
}

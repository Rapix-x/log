package log

var logger = MustNewLogger(Configuration{MinimumLogLevel: DebugLevel})

// Debug logs all inputs on the debug level.
func Debug(v ...any) {
	logger.Debug(v...)
}

// Debugf formats and logs all inputs on the debug level.
func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
}

// Debugw logs all inputs and fields on the debug level.
func Debugw(msg string, keyValuePairs ...any) {
	logger.Debugw(msg, keyValuePairs...)
}

// Error logs all inputs on the error level.
func Error(v ...any) {
	logger.Error(v...)
}

// Errorf formats and logs all inputs on the error level.
func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

// Errorw logs all inputs and fields on the error level.
func Errorw(msg string, keyValuePairs ...any) {
	logger.Errorw(msg, keyValuePairs...)
}

// Fatal logs all inputs on the fatal level and runs os.exit(1) at
// the end.
func Fatal(v ...any) {
	logger.Fatal(v...)
}

// Fatalf formats and logs all inputs on the fatal level and runs
// os.exit(1) at the end.
func Fatalf(format string, v ...any) {
	logger.Fatalf(format, v...)
}

// Fatalw logs all inputs and fields on the fatal level and runs
// os.exit(1) at the end.
func Fatalw(msg string, keyValuePairs ...any) {
	logger.Fatalw(msg, keyValuePairs...)
}

// Info logs all inputs on the info level.
func Info(v ...any) {
	logger.Info(v...)
}

// Infof formats and logs all inputs on the info level.
func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

// Infow logs all inputs and fields on the info level.
func Infow(msg string, keyValuePairs ...any) {
	logger.Infow(msg, keyValuePairs...)
}

// Warn logs all inputs on the warn level.
func Warn(v ...any) {
	logger.Warn(v...)
}

// Warnf formats and logs all inputs on the warn level.
func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

// Warnw logs all inputs and fields on the warn level.
func Warnw(msg string, keyValuePairs ...any) {
	logger.Warnw(msg, keyValuePairs...)
}

func Sync() error {
	return logger.Sync()
}

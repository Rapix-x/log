package log

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level specifies a log level. Usually it is used to indicate the
// minimum log level for a logger.
type Level zapcore.Level

const (
	DebugLevel = Level(zapcore.DebugLevel)
	InfoLevel  = Level(zapcore.InfoLevel)
	WarnLevel  = Level(zapcore.WarnLevel)
	ErrorLevel = Level(zapcore.ErrorLevel)
	PanicLevel = Level(zapcore.PanicLevel)
	FatalLevel = Level(zapcore.FatalLevel)
)

var (
	logLevels = map[Level]struct{}{
		DebugLevel: {},
		InfoLevel:  {},
		WarnLevel:  {},
		ErrorLevel: {},
		PanicLevel: {},
		FatalLevel: {},
	}
)

type encoding string

const (
	encodingConsole encoding = "console" // nolint
	encodingJSON    encoding = "json"
)

var encoderConfig = zapcore.EncoderConfig{
	MessageKey:          "msg",
	LevelKey:            "lvl",
	TimeKey:             "ts",
	NameKey:             "name",
	CallerKey:           "caller",
	FunctionKey:         "func",
	StacktraceKey:       "stacktrace",
	SkipLineEnding:      false,
	LineEnding:          "\n",
	EncodeLevel:         zapcore.LowercaseLevelEncoder,
	EncodeTime:          zapcore.RFC3339TimeEncoder,
	EncodeDuration:      zapcore.MillisDurationEncoder,
	EncodeCaller:        zapcore.ShortCallerEncoder,
	EncodeName:          nil,
	NewReflectedEncoder: nil,
}

var loggerConf = zap.Config{
	Development:       false,
	DisableCaller:     false,
	DisableStacktrace: true,
	Encoding:          string(encodingJSON),
	EncoderConfig:     encoderConfig,
	OutputPaths:       []string{"stdout"},
	ErrorOutputPaths:  []string{"stderr"},
}

// Configuration represents a Configuration object for a logger.
type Configuration struct {
	// ApplicationName holds the value for the "app" field in log
	// statements indicating the name of the current application.
	// If the value is set to "", the field will be omitted.
	ApplicationName string

	// Version holds the value for the "version" field in log
	// statements indicating the version of the current application.
	// If the value is set to "", the field will be omitted.
	Version string

	// MinimumLogLevel sets the minim level of logs that will get
	// logged by the respective logger. The DebugLevel is the lowest
	// while the FatalLevel is the highest. If set to Debug, everything
	// will be logged, while when set to Fatal, only Fatal statements
	// will be logged.
	MinimumLogLevel Level

	// PIIMode indicates how to the logger resolves PII fields in log
	// statements.
	PIIMode PIIMode
}

// The Logger struct resembles the actual loggers.
type Logger struct {
	logger  *zap.SugaredLogger
	piiMode PIIMode
}

// NewNOPLogger creates a new no-operation logger that does not write
// any log statements anywhere and is therefore tremendously helpful,
// when you need to fulfill the Interface, but you don't want to
// actually log anything.
func NewNOPLogger() *Logger {
	return &Logger{logger: zap.NewNop().Sugar()}
}

// MustNewLogger wraps NewLogger and panics, when an error is encountered.
func MustNewLogger(c Configuration) *Logger {
	l, e := NewLogger(c)
	if e != nil {
		panic(e)
	}

	return l
}

// NewLogger creates a new logger based on the configuration inputs and
// returns a pointer to it. If the validation of the input configuration
// fails an error will be issued.
func NewLogger(conf Configuration) (*Logger, error) {
	err := validateLoggerConf(conf)
	if err != nil {
		return nil, errors.Wrap(err, "received an error while validating the logger configuration")
	}

	loggerConf.Level = zap.NewAtomicLevelAt(zapcore.Level(conf.MinimumLogLevel))
	loggerConf.InitialFields = make(map[string]interface{}, 0)

	if conf.ApplicationName != "" {
		loggerConf.InitialFields["app"] = conf.ApplicationName
	}

	if conf.Version != "" {
		loggerConf.InitialFields["version"] = conf.Version
	}

	zapLogger, err := loggerConf.Build(
		zap.AddStacktrace(zapcore.WarnLevel),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create logger")
	}

	return &Logger{
		logger:  zapLogger.Sugar(),
		piiMode: conf.PIIMode,
	}, nil
}

// Debug logs all inputs on the debug level.
func (l *Logger) Debug(v ...any) {
	handleUninitialized(l)
	l.logger.Debug(v...)
}

// Debugf formats and logs all inputs on the debug level.
func (l *Logger) Debugf(format string, v ...any) {
	handleUninitialized(l)
	l.logger.Debugf(format, v...)
}

// Debugw logs all inputs and fields on the debug level.
func (l *Logger) Debugw(msg string, keyValuePairs ...any) {
	handleUninitialized(l)
	l.logger.Debugw(msg, resolvePIIFunctions(l.piiMode, keyValuePairs)...)
}

// Error logs all inputs on the error level.
func (l *Logger) Error(v ...any) {
	handleUninitialized(l)
	l.logger.Error(v...)
}

// Errorf formats and logs all inputs on the error level.
func (l *Logger) Errorf(format string, v ...any) {
	handleUninitialized(l)
	l.logger.Errorf(format, v...)
}

// Errorw logs all inputs and fields on the error level.
func (l *Logger) Errorw(msg string, keyValuePairs ...any) {
	handleUninitialized(l)
	l.logger.Errorw(msg, resolvePIIFunctions(l.piiMode, keyValuePairs)...)
}

// Fatal logs all inputs on the fatal level and runs os.exit(1) at
// the end.
func (l *Logger) Fatal(v ...any) {
	handleUninitialized(l)
	l.logger.Fatal(v...)
}

// Fatalf formats and logs all inputs on the fatal level and runs
// os.exit(1) at the end.
func (l *Logger) Fatalf(format string, v ...any) {
	handleUninitialized(l)
	l.logger.Fatalf(format, v...)
}

// Fatalw logs all inputs and fields on the fatal level and runs
// os.exit(1) at the end.
func (l *Logger) Fatalw(msg string, keyValuePairs ...any) {
	handleUninitialized(l)
	l.logger.Fatalw(msg, resolvePIIFunctions(l.piiMode, keyValuePairs)...)
}

// Info logs all inputs on the info level.
func (l *Logger) Info(v ...any) {
	handleUninitialized(l)
	l.logger.Info(v...)
}

// Infof formats and logs all inputs on the info level.
func (l *Logger) Infof(format string, v ...any) {
	handleUninitialized(l)
	l.logger.Infof(format, v...)
}

// Infow logs all inputs and fields on the info level.
func (l *Logger) Infow(msg string, keyValuePairs ...any) {
	handleUninitialized(l)
	fields := resolvePIIFunctions(l.piiMode, keyValuePairs)
	l.logger.Infow(msg, fields...)
}

func (l *Logger) Sync() error {
	handleUninitialized(l)

	return l.logger.Sync()
}

// Warn logs all inputs on the warn level.
func (l *Logger) Warn(v ...any) {
	handleUninitialized(l)
	l.logger.Warn(v...)
}

// Warnf formats and logs all inputs on the warn level.
func (l *Logger) Warnf(format string, v ...any) {
	handleUninitialized(l)
	l.logger.Warnf(format, v...)
}

// Warnw logs all inputs and fields on the warn level.
func (l *Logger) Warnw(msg string, keyValuePairs ...any) {
	handleUninitialized(l)
	l.logger.Warnw(msg, resolvePIIFunctions(l.piiMode, keyValuePairs)...)
}

// With returns a pointer to a new logger containing the added fields.
func (l *Logger) With(keyValuePairs ...any) *Logger {
	handleUninitialized(l)

	return &Logger{
		logger:  l.logger.With(resolvePIIFunctions(l.piiMode, keyValuePairs)...),
		piiMode: l.piiMode,
	}
}

func handleUninitialized(l *Logger) {
	if l == nil {
		ephemeralLogger := zap.Must(zap.NewProduction(zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.FatalLevel)))
		ephemeralLogger.Panic("logger has not been initialized - panicking")
	}
}

// The PIIResolver interface is what the logger checks against,
// when trying to resolve PII fields in log statements before writing
// the logs.
type PIIResolver interface {
	resolve(piiMode PIIMode) zap.Field
}

func resolvePIIFunctions(piiMode PIIMode, keyValuePairs []any) []any {
	out := make([]any, 0)

	for _, element := range keyValuePairs {
		if e, ok := element.(PIIResolver); ok {
			out = append(out, e.resolve(piiMode))

			continue
		}

		out = append(out, element)
	}

	return out
}

func validateLoggerConf(conf Configuration) error {
	if _, ok := logLevels[conf.MinimumLogLevel]; !ok {
		return errors.New("invalid minimum log level in logger configuration")
	}

	if _, ok := piiModes[conf.PIIMode]; !ok {
		return errors.New("invalid PII mode in logger configuration")
	}

	return nil
}

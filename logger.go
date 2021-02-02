package log

import (
	"io"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*RotationConfig

	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
	StdLogger     *stdlog.Logger

	Development  bool
	Format       Format
	Level        Level
	LogDirs      []string
	LogToConsole bool
	InfoOutput   io.Writer
	ErrOutput    io.Writer

	Caller bool
}

func New(opts ...Option) *Logger {
	rotationConfig := defaultRotationConfig

	l := &Logger{
		RotationConfig: &rotationConfig,
		Development:    false,
		Format:         FormatJSON,
		Level:          InfoLevel,
		LogToConsole:   true,
	}

	for _, o := range opts {
		o.apply(l)
	}

	l.NotifyOptionsChange()

	return l
}

func (l *Logger) initZapLogger() {
	var encoderCfg zapcore.EncoderConfig
	if l.Development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var encoder zapcore.Encoder
	switch l.Format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		lvll := fromZapLevel(lvl)
		return l.levelEnabled(lvll) && lvll >= ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		lvll := fromZapLevel(lvl)
		return l.levelEnabled(lvll) && lvll < ErrorLevel
	})

	cores := make([]zapcore.Core, 0)

	if l.LogToConsole {
		cores = append(cores,
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lowPriority),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), highPriority),
		)
	}

	if l.InfoOutput != nil {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(l.InfoOutput)), lowPriority))
	}

	if l.ErrOutput != nil {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(l.ErrOutput)), highPriority))
	}

	for _, dir := range l.LogDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		infoWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "info.log"),
			MaxSize:    l.MaxSize,
			MaxAge:     l.MaxAge,
			MaxBackups: l.MaxBackups,
			LocalTime:  l.LocalTime,
			Compress:   l.Compress,
		})

		errWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "error.log"),
			MaxSize:    l.MaxSize,
			MaxAge:     l.MaxAge,
			MaxBackups: l.MaxBackups,
			LocalTime:  l.LocalTime,
			Compress:   l.Compress,
		})

		cores = append(cores,
			zapcore.NewCore(encoder, infoWriter, lowPriority),
			zapcore.NewCore(encoder, errWriter, highPriority),
		)
	}

	zapLogger := zap.New(zapcore.NewTee(cores...))
	if l.Development {
		zapLogger.WithOptions(zap.Development())
	}
	if l.Caller {
		zapLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}

	l.logger = zapLogger
	l.sugaredLogger = zapLogger.Sugar()
	l.StdLogger, _ = zap.NewStdLogAt(zapLogger, toZapLevel(l.Level))
}

func (l *Logger) levelEnabled(lvl Level) bool {
	return l.Development || lvl >= l.Level
}

func (l *Logger) WithOptions(opts ...Option) {
	for _, o := range opts {
		o.apply(l)
	}

	l.NotifyOptionsChange()
}

func (l *Logger) NotifyOptionsChange() {
	l.initZapLogger()
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugaredLogger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Debug(args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugaredLogger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugaredLogger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugaredLogger.Warnf(template, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugaredLogger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugaredLogger.Errorf(template, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Error(args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.sugaredLogger.DPanic(args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.sugaredLogger.DPanicf(template, args...)
}

func (l *Logger) DPanicln(args ...interface{}) {
	l.DPanic(args...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sugaredLogger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugaredLogger.Panicf(template, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.Panic(args...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugaredLogger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugaredLogger.Fatalf(template, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugaredLogger.Fatalw(msg, keysAndValues...)
}

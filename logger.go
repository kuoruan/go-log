package log

import (
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*config
	logger *zap.SugaredLogger
}

func New(opts ...Option) *Logger {
	l := &Logger{
		config: &config{
			development:  false,
			format:       FormatConsole,
			level:        InfoLevel,
			logToConsole: true,
			addCaller:    false,
			callerSkip:   1,

			maxAge:     28,
			maxBackups: 7,
			maxSize:    500,
			localTime:  true,
			compress:   false,
		},
		logger: zap.NewNop().Sugar(),
	}

	for _, o := range opts {
		o.apply(l)
	}

	l.updateLogger()
	return l
}

// update zap logger based on the new config
func (l *Logger) updateLogger() {
	var encoderCfg zapcore.EncoderConfig
	if l.development {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var encoder zapcore.Encoder
	switch l.format {
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		lvll := fromZapLevel(lvl)
		return l.levelEnabled(lvll) && lvll < ErrorLevel
	})
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		lvll := fromZapLevel(lvl)
		return l.levelEnabled(lvll) && lvll >= ErrorLevel
	})

	cores := make([]zapcore.Core, 0)

	// add console log
	if l.logToConsole {
		cores = append(cores,
			zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), lowPriority),
			zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), highPriority),
		)
	}

	if l.infoOutput != nil {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(l.infoOutput)), lowPriority))
	}

	if l.errOutput != nil {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(l.errOutput)), highPriority))
	}

	for _, dir := range l.logDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			_ = os.MkdirAll(dir, 0755)
		}

		infoWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "info.log"),
			MaxSize:    l.maxSize,
			MaxAge:     l.maxAge,
			MaxBackups: l.maxBackups,
			LocalTime:  l.localTime,
			Compress:   l.compress,
		})

		errWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(dir, "error.log"),
			MaxSize:    l.maxSize,
			MaxAge:     l.maxAge,
			MaxBackups: l.maxBackups,
			LocalTime:  l.localTime,
			Compress:   l.compress,
		})

		cores = append(cores,
			zapcore.NewCore(encoder, infoWriter, lowPriority),
			zapcore.NewCore(encoder, errWriter, highPriority),
		)
	}

	zapLogger := l.logger.Desugar()

	zapOptions := []zap.Option{
		// set new zap cores
		zap.WrapCore(func(zapcore.Core) zapcore.Core {
			if len(cores) > 0 {
				return zapcore.NewTee(cores...)
			}
			return zapcore.NewNopCore()
		}),
	}

	if l.development {
		zapOptions = append(zapOptions, zap.Development())
	}

	if l.addCaller {
		zapOptions = append(zapOptions, zap.WithCaller(true), zap.AddCallerSkip(l.callerSkip))
	} else {
		zapOptions = append(zapOptions, zap.WithCaller(false))
	}

	l.logger = zapLogger.WithOptions(zapOptions...).Sugar()
}

func (l *Logger) levelEnabled(lvl Level) bool {
	return l.development || l.level.Enabled(lvl)
}

func (l *Logger) WithOptions(opts ...Option) *Logger {
	newLogger := &Logger{
		config: l.config.clone(),
		logger: zap.NewNop().Sugar(),
	}

	for _, o := range opts {
		o.apply(newLogger)
	}

	newLogger.updateLogger()
	return newLogger
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Debug(args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Error(args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.logger.DPanic(args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.logger.DPanicf(template, args...)
}

func (l *Logger) DPanicln(args ...interface{}) {
	l.DPanic(args...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.logger.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.Panic(args...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

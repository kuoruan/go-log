package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logFunc func(logger *zap.SugaredLogger, args ...interface{})
type logfFunc func(logger *zap.SugaredLogger, format string, args ...interface{})
type logwFunc func(logger *zap.SugaredLogger, msg string, keysAndValues ...interface{})

type Logger struct {
	*config
	logger *zap.SugaredLogger

	debug  logFunc
	debugf logfFunc
	debugw logwFunc

	info  logFunc
	infof logfFunc
	infow logwFunc

	warn  logFunc
	warnf logfFunc
	warnw logwFunc

	error  logFunc
	errorf logfFunc
	errorw logwFunc

	dpanic  logFunc
	dpanicf logfFunc
	dpanicw logwFunc

	panic  logFunc
	panicf logfFunc
	panicw logwFunc

	fatal  logFunc
	fatalf logfFunc
	fatalw logwFunc
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

		logger:  zap.NewNop().Sugar(),
		debug:   (*zap.SugaredLogger).Debug,
		debugf:  (*zap.SugaredLogger).Debugf,
		debugw:  (*zap.SugaredLogger).Debugw,
		info:    (*zap.SugaredLogger).Info,
		infof:   (*zap.SugaredLogger).Infof,
		infow:   (*zap.SugaredLogger).Infow,
		warn:    (*zap.SugaredLogger).Warn,
		warnf:   (*zap.SugaredLogger).Warnf,
		warnw:   (*zap.SugaredLogger).Warnw,
		error:   (*zap.SugaredLogger).Error,
		errorf:  (*zap.SugaredLogger).Errorf,
		errorw:  (*zap.SugaredLogger).Errorw,
		dpanic:  (*zap.SugaredLogger).DPanic,
		dpanicf: (*zap.SugaredLogger).DPanicf,
		dpanicw: (*zap.SugaredLogger).DPanicw,
		panic:   (*zap.SugaredLogger).Panic,
		panicf:  (*zap.SugaredLogger).Panicf,
		panicw:  (*zap.SugaredLogger).Panicw,
		fatal:   (*zap.SugaredLogger).Fatal,
		fatalf:  (*zap.SugaredLogger).Fatalf,
		fatalw:  (*zap.SugaredLogger).Fatalw,
	}

	for _, o := range opts {
		o.apply(l)
	}

	l.updateLogger()
	return l
}

func (l *Logger) newCopy() *Logger {
	return &Logger{
		config:  l.config.clone(),
		logger:  zap.NewNop().Sugar(),
		debug:   (*zap.SugaredLogger).Debug,
		debugf:  (*zap.SugaredLogger).Debugf,
		debugw:  (*zap.SugaredLogger).Debugw,
		info:    (*zap.SugaredLogger).Info,
		infof:   (*zap.SugaredLogger).Infof,
		infow:   (*zap.SugaredLogger).Infow,
		warn:    (*zap.SugaredLogger).Warn,
		warnf:   (*zap.SugaredLogger).Warnf,
		warnw:   (*zap.SugaredLogger).Warnw,
		error:   (*zap.SugaredLogger).Error,
		errorf:  (*zap.SugaredLogger).Errorf,
		errorw:  (*zap.SugaredLogger).Errorw,
		dpanic:  (*zap.SugaredLogger).DPanic,
		dpanicf: (*zap.SugaredLogger).DPanicf,
		dpanicw: (*zap.SugaredLogger).DPanicw,
		panic:   (*zap.SugaredLogger).Panic,
		panicf:  (*zap.SugaredLogger).Panicf,
		panicw:  (*zap.SugaredLogger).Panicw,
		fatal:   (*zap.SugaredLogger).Fatal,
		fatalf:  (*zap.SugaredLogger).Fatalf,
		fatalw:  (*zap.SugaredLogger).Fatalw,
	}
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
	c := l.newCopy()

	for _, o := range opts {
		o.apply(c)
	}

	c.updateLogger()
	return c
}

func (l *Logger) Debug(args ...interface{}) {
	l.debug(l.logger, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debugf(l.logger, format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.debug(l.logger, sprintln(args...))
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.debugw(l.logger, msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.info(l.logger, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.infof(l.logger, format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.info(l.logger, sprintln(args...))
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.infow(l.logger, msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.warn(l.logger, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.warnf(l.logger, template, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.warn(l.logger, sprintln(args...))
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.warnw(l.logger, msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.error(l.logger, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.errorf(l.logger, template, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.error(l.logger, sprintln(args...))
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.errorw(l.logger, msg, keysAndValues...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.dpanic(l.logger, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.dpanicf(l.logger, template, args...)
}

func (l *Logger) DPanicln(args ...interface{}) {
	l.dpanic(l.logger, sprintln(args...))
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.dpanicw(l.logger, msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.panic(l.logger, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.panicf(l.logger, template, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.panic(l.logger, sprintln(args...))
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.panicw(l.logger, msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.fatal(l.logger, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.fatalf(l.logger, template, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.fatal(l.logger, sprintln(args...))
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.fatalw(l.logger, msg, keysAndValues...)
}

func sprintln(args ...interface{}) string {
	return strings.TrimSuffix(fmt.Sprintln(args...), "\n")
}

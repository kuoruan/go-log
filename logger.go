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
	base    *zap.SugaredLogger
	writers []*lumberjack.Logger

	print  logFunc
	printf logfFunc
	printw logwFunc

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
			development: false,
			format:      FormatJSON,
			level:       InfoLevel,
			logToStdout: true,
			addCaller:   false,
			callerSkip:  1,

			maxAge:     28,
			maxBackups: 7,
			maxSize:    500,
			localTime:  true,
			compress:   false,
		},
		base: zap.NewNop().Sugar(),
	}

	for _, o := range opts {
		o.apply(l)
	}

	l.initLogFuncs()
	l.updateLogger()
	return l
}

func (l *Logger) initLogFuncs() {
	l.debug = (*zap.SugaredLogger).Debug
	l.debugf = (*zap.SugaredLogger).Debugf
	l.debugw = (*zap.SugaredLogger).Debugw
	l.info = (*zap.SugaredLogger).Info
	l.infof = (*zap.SugaredLogger).Infof
	l.infow = (*zap.SugaredLogger).Infow
	l.warn = (*zap.SugaredLogger).Warn
	l.warnf = (*zap.SugaredLogger).Warnf
	l.warnw = (*zap.SugaredLogger).Warnw
	l.error = (*zap.SugaredLogger).Error
	l.errorf = (*zap.SugaredLogger).Errorf
	l.errorw = (*zap.SugaredLogger).Errorw
	l.dpanic = (*zap.SugaredLogger).DPanic
	l.dpanicf = (*zap.SugaredLogger).DPanicf
	l.dpanicw = (*zap.SugaredLogger).DPanicw
	l.panic = (*zap.SugaredLogger).Panic
	l.panicf = (*zap.SugaredLogger).Panicf
	l.panicw = (*zap.SugaredLogger).Panicw
	l.fatal = (*zap.SugaredLogger).Fatal
	l.fatalf = (*zap.SugaredLogger).Fatalf
	l.fatalw = (*zap.SugaredLogger).Fatalw

	if l.development {
		l.print = (*zap.SugaredLogger).Debug
		l.printf = (*zap.SugaredLogger).Debugf
		l.printw = (*zap.SugaredLogger).Debugw
	} else {
		l.print = (*zap.SugaredLogger).Info
		l.printf = (*zap.SugaredLogger).Infof
		l.printw = (*zap.SugaredLogger).Infow
	}
}

// update zap base based on the new config
func (l *Logger) updateLogger() {
	var encoder zapcore.Encoder

	if l.encoder != nil {
		encoder = l.encoder
	} else {
		var encoderCfg zapcore.EncoderConfig

		if l.development {
			encoderCfg = zap.NewDevelopmentEncoderConfig()
		} else {
			encoderCfg = zap.NewProductionEncoderConfig()
		}

		switch l.format {
		case FormatJSON:
			encoder = zapcore.NewJSONEncoder(encoderCfg)
		default:
			encoder = zapcore.NewConsoleEncoder(encoderCfg)
		}
	}

	cores := make([]zapcore.Core, 0)
	writers := make([]*lumberjack.Logger, 0)

	// add stdout log
	if l.logToStdout {
		stdoutCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(l.zapLevelEnabled),
		)
		cores = append(cores, stdoutCore)
	}

	// add output core
	if l.output != nil {
		outputCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(zapcore.AddSync(l.output)),
			zap.LevelEnablerFunc(l.zapLevelEnabled),
		)
		cores = append(cores, outputCore)
	}

	// parse log dirs
	for _, dir := range l.logDirs {
		if dir == "" {
			continue
		}

		for _, level := range []zapcore.Level{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			zapcore.WarnLevel,
			zapcore.ErrorLevel,
			zapcore.DPanicLevel,
			zapcore.PanicLevel,
			zapcore.FatalLevel,
		} {
			if l.zapLevelEnabled(level) {
				lvl := level

				lvlWriter := &lumberjack.Logger{
					Filename:   filepath.Join(dir, fmt.Sprint(lvl.String(), ".log")),
					MaxSize:    l.maxSize,
					MaxAge:     l.maxAge,
					MaxBackups: l.maxBackups,
					LocalTime:  l.localTime,
					Compress:   l.compress,
				}

				lvlCore := zapcore.NewCore(
					encoder,
					zapcore.AddSync(lvlWriter),
					zap.LevelEnablerFunc(func(l zapcore.Level) bool {
						return l == lvl
					}),
				)

				cores = append(cores, lvlCore)
				writers = append(writers, lvlWriter)
			}
		}
	}

	// parse log files
	for _, file := range l.logFiles {
		if file == "" {
			continue
		}

		writer := &lumberjack.Logger{
			Filename:   file,
			MaxSize:    l.maxSize,
			MaxAge:     l.maxAge,
			MaxBackups: l.maxBackups,
			LocalTime:  l.localTime,
			Compress:   l.compress,
		}

		fileCore := zapcore.NewCore(encoder, zapcore.AddSync(writer), zap.LevelEnablerFunc(l.zapLevelEnabled))

		cores = append(cores, fileCore)
		writers = append(writers, writer)
	}

	zapLogger := l.base.Desugar()

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

	l.base = zapLogger.WithOptions(zapOptions...).Sugar()
	l.writers = writers
}

func (l *Logger) zapLevelEnabled(lvl zapcore.Level) bool {
	return l.development || l.level.Enabled(fromZapLevel(lvl))
}

func (l *Logger) WithOptions(opts ...Option) *Logger {
	c := &Logger{
		config: l.config.clone(),
		base:   zap.NewNop().Sugar(),
	}

	for _, o := range opts {
		o.apply(c)
	}

	c.initLogFuncs()
	c.updateLogger()
	return c
}

func (l *Logger) Print(args ...interface{}) {
	l.print(l.base, args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.printf(l.base, format, args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.print(l.base, sprintln(args...))
}

func (l *Logger) Printw(msg string, keysAndValues ...interface{}) {
	l.printw(l.base, msg, keysAndValues...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.debug(l.base, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.debugf(l.base, format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.debug(l.base, sprintln(args...))
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.debugw(l.base, msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.info(l.base, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.infof(l.base, format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.info(l.base, sprintln(args...))
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.infow(l.base, msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.warn(l.base, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.warnf(l.base, template, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.warn(l.base, sprintln(args...))
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.warnw(l.base, msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.error(l.base, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.errorf(l.base, template, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.error(l.base, sprintln(args...))
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.errorw(l.base, msg, keysAndValues...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.dpanic(l.base, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.dpanicf(l.base, template, args...)
}

func (l *Logger) DPanicln(args ...interface{}) {
	l.dpanic(l.base, sprintln(args...))
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.dpanicw(l.base, msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.panic(l.base, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.panicf(l.base, template, args...)
}

func (l *Logger) Panicln(args ...interface{}) {
	l.panic(l.base, sprintln(args...))
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.panicw(l.base, msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.fatal(l.base, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.fatalf(l.base, template, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.fatal(l.base, sprintln(args...))
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.fatalw(l.base, msg, keysAndValues...)
}

func (l *Logger) Rotate() error {
	for _, w := range l.writers {
		if err := w.Rotate(); err != nil {
			return err
		}
	}

	return nil
}

func sprintln(args ...interface{}) string {
	return strings.TrimSuffix(fmt.Sprintln(args...), "\n")
}

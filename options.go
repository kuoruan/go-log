package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WithLevel(lvl Level) Option {
	return optionFunc(func(l *Logger) {
		l.level = lvl
	})
}

func WithFormat(format Format) Option {
	return optionFunc(func(l *Logger) {
		l.format = format
	})
}

func Development() Option {
	return WithDevelopment(true)
}

func WithDevelopment(development bool) Option {
	return optionFunc(func(l *Logger) {
		l.development = development
	})
}

func WithEncoder(encoder zapcore.Encoder) Option {
	return optionFunc(func(l *Logger) {
		l.encoder = encoder
	})
}

func WithOutput(output io.Writer) Option {
	return optionFunc(func(l *Logger) {
		l.output = output
	})
}

func WithLogToStdout(logToStdout bool) Option {
	return optionFunc(func(l *Logger) {
		l.logToStdout = logToStdout
	})
}

func WithLogDirs(dirs ...string) Option {
	return optionFunc(func(l *Logger) {
		dst := make([]string, len(dirs))
		copy(dst, dirs)

		l.logDirs = dst
	})
}

func WithLogFiles(files ...string) Option {
	return optionFunc(func(l *Logger) {
		dst := make([]string, len(files))
		copy(dst, files)

		l.logFiles = dst
	})
}

func WithRotationConfig(config RotationConfig) Option {
	return optionFunc(func(l *Logger) {
		if config.MaxAge > 0 {
			l.maxAge = config.MaxAge
		}

		if config.MaxBackups > 0 {
			l.maxBackups = config.MaxBackups
		}

		if config.MaxSize > 0 {
			l.maxSize = config.MaxSize
		}

		l.compress = config.Compress
		l.localTime = config.LocalTime
	})
}

func WithCaller(caller bool) Option {
	return optionFunc(func(l *Logger) {
		l.addCaller = caller
	})
}

func AddCaller() Option {
	return WithCaller(true)
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *Logger) {
		l.callerSkip += skip
	})
}

type RotationConfig struct {
	MaxSize    int  `json:"maxSize"`    // megabytes
	MaxAge     int  `json:"maxAge"`     // days
	MaxBackups int  `json:"maxBackups"` // count
	LocalTime  bool `json:"localTime"`
	Compress   bool `json:"compress"`
}

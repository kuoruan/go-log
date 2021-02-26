package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type RotationConfig struct {
	MaxSize    int  `json:"maxSize"`    // megabytes
	MaxAge     int  `json:"maxAge"`     // days
	MaxBackups int  `json:"maxBackups"` // count
	LocalTime  bool `json:"localTime"`
	Compress   bool `json:"compress"`
}

func (c RotationConfig) apply(o *options) {
	if c.MaxAge > 0 {
		o.MaxAge = c.MaxAge
	}

	if c.MaxBackups > 0 {
		o.MaxBackups = c.MaxBackups
	}

	if c.MaxSize > 0 {
		o.MaxSize = c.MaxSize
	}

	o.Compress = c.Compress
	o.LocalTime = c.LocalTime
}

type options struct {
	RotationConfig

	Level   Level
	Format  Format
	Encoder zapcore.Encoder

	Development bool

	Output      io.Writer
	LogToStdout bool
	LogDirs     []string
	LogFiles    []string

	AddCaller  bool
	CallerSkip int
}

func (o options) Clone() options {
	c := o

	if c.Encoder != nil {
		c.Encoder = c.Encoder.Clone()
	}

	if len(c.LogDirs) > 0 {
		d := make([]string, len(c.LogDirs))

		copy(d, c.LogDirs)
		c.LogDirs = d
	}

	if len(c.LogFiles) > 0 {
		d := make([]string, len(c.LogFiles))

		copy(d, c.LogFiles)
		c.LogFiles = d
	}

	return c
}

func (o options) ZapLevelEnabled(lvl zapcore.Level) bool {
	return o.Development || o.Level.Enabled(fromZapLevel(lvl))
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(log *options) {
	f(log)
}

func WithLevel(lvl Level) Option {
	return optionFunc(func(l *options) {
		l.Level = lvl
	})
}

func WithFormat(format Format) Option {
	return optionFunc(func(l *options) {
		l.Format = format
	})
}

func Development() Option {
	return WithDevelopment(true)
}

func WithDevelopment(development bool) Option {
	return optionFunc(func(l *options) {
		l.Development = development
	})
}

func WithEncoder(encoder zapcore.Encoder) Option {
	return optionFunc(func(l *options) {
		l.Encoder = encoder
	})
}

func WithOutput(output io.Writer) Option {
	return optionFunc(func(l *options) {
		l.Output = output
	})
}

func LogToStdout() Option {
	return WithLogToStdout(true)
}

func WithLogToStdout(logToStdout bool) Option {
	return optionFunc(func(l *options) {
		l.LogToStdout = logToStdout
	})
}

func WithLogDirs(dirs ...string) Option {
	return optionFunc(func(l *options) {
		dst := make([]string, len(dirs))
		copy(dst, dirs)

		l.LogDirs = dst
	})
}

func WithLogFiles(files ...string) Option {
	return optionFunc(func(l *options) {
		dst := make([]string, len(files))
		copy(dst, files)

		l.LogFiles = dst
	})
}

func WithCaller(caller bool) Option {
	return optionFunc(func(l *options) {
		l.AddCaller = caller
	})
}

func AddCaller() Option {
	return WithCaller(true)
}

func AddCallerSkip(skip int) Option {
	return optionFunc(func(l *options) {
		l.CallerSkip += skip
	})
}

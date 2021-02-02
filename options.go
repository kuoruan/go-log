package log

import "io"

type Option interface {
	apply(*Logger)
}

type optionFunc func(*Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WithLevel(lvl Level) Option {
	return optionFunc(func(l *Logger) {
		l.Level = lvl
	})
}

func WithFormat(format Format) Option {
	return optionFunc(func(l *Logger) {
		l.Format = format
	})
}

func WithDevelopment(development bool) Option {
	return optionFunc(func(l *Logger) {
		l.Development = development
	})
}

func WithOutput(infoOut, errOut io.Writer) Option {
	return optionFunc(func(l *Logger) {
		l.InfoOutput = infoOut
		l.ErrOutput = errOut
	})
}

func WithLogToConsole(logToConsole bool) Option {
	return optionFunc(func(l *Logger) {
		l.LogToConsole = logToConsole
	})
}

func WithLogDirs(dirs ...string) Option {
	return optionFunc(func(l *Logger) {
		d := make([]string, len(dirs))
		copy(d, dirs)

		l.LogDirs = d
	})
}

func WithRotationConfig(config RotationConfig) Option {
	return optionFunc(func(l *Logger) {
		c := defaultRotationConfig

		if config.MaxAge > 0 {
			c.MaxAge = config.MaxAge
		}

		if config.MaxBackups > 0 {
			c.MaxBackups = config.MaxBackups
		}

		if config.MaxSize > 0 {
			c.MaxSize = config.MaxSize
		}

		l.RotationConfig = &c
	})
}

func WithCaller(caller bool) Option {
	return optionFunc(func(l *Logger) {
		l.Caller = caller
	})
}

type RotationConfig struct {
	MaxSize    int  `json:"maxSize"`    // megabytes
	MaxAge     int  `json:"maxAge"`     // days
	MaxBackups int  `json:"maxBackups"` // count
	LocalTime  bool `json:"localTime"`
	Compress   bool `json:"compress"`
}

var defaultRotationConfig = RotationConfig{
	MaxSize:    500,
	MaxAge:     28,
	MaxBackups: 3,
	LocalTime:  true,
	Compress:   true,
}

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

func WithOutput(infoOut, errOut io.Writer) Option {
	return optionFunc(func(l *Logger) {
		l.infoOutput = infoOut
		l.errOutput = errOut
	})
}

func WithLogToConsole(logToConsole bool) Option {
	return optionFunc(func(l *Logger) {
		l.logToConsole = logToConsole
	})
}

func WithLogDirs(dirs ...string) Option {
	return optionFunc(func(l *Logger) {
		d := make([]string, len(dirs))
		copy(d, dirs)

		l.logDirs = d
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

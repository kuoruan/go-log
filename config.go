package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

type config struct {
	development bool
	encoder     zapcore.Encoder
	format      Format
	output      io.Writer
	level       Level
	logDirs     []string
	logFiles    []string
	logToStdout bool
	addCaller   bool
	callerSkip  int

	maxAge     int
	maxBackups int
	maxSize    int
	compress   bool
	localTime  bool
}

func (c *config) clone() *config {
	cloned := *c
	return &cloned
}

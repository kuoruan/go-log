package log

import "io"

type config struct {
	development  bool
	format       Format
	errOutput    io.Writer
	infoOutput   io.Writer
	level        Level
	logDirs      []string
	logToConsole bool
	addCaller    bool
	callerSkip   int

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

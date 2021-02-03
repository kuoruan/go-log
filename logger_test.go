package log

import "testing"

func TestLogger_WithOptions(t *testing.T) {
	l := New(Development())
	n := l.WithOptions(AddCaller())

	l.Debug("debug original")
	n.Debug("debug new")
}

func TestLogger_Debug(t *testing.T) {
	l := New(Development())

	l.Debug("debug")
}

func TestLogger_Debugf(t *testing.T) {
	l := New(Development())

	l.Debugf("debugf %s", "zap")
}

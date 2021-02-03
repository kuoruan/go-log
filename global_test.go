package log

import "testing"

func TestSetOptions(t *testing.T) {
	SetOptions(Development(), WithCaller(true))

	Debug("debug")
}

package logger

import "testing"

func TestLogger(t *testing.T) {
	t.Run("TestLogger", func(t *testing.T) {
		Log.Infoln("This is an info message")
		Log.Warnln("This is a warning message")
		Log.Errorln("This is an error message")
		Log.Debugln("This is a debug message")
	})
}

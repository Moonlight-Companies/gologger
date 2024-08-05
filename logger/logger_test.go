package logger

import "testing"

func TestLogger(t *testing.T) {
	t.Run("TestLogger", func(t *testing.T) {
		Log.SetIncludeDeltaTime(true)
		Log.Infoln("This is an info message")
		Log.Warnln("This is a warning message")
		Log.Errorln("This is an error message")
		Log.Debugln("This is a debug message")

		Log.SetPrefix("TEST")
		if Log.GetPrefix() != "TEST" {
			t.Errorf("Expected prefix to be 'TEST', got %s", Log.GetPrefix())
		}

		Log.Infoln("This is an info message")
		Log.Warnln("This is a warning message")
		Log.Errorln("This is an error message")
		Log.Debugln("This is a debug message")
	})
}

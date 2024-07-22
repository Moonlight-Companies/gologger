package logger

import "testing"

func TestLogger(t *testing.T) {
	t.Run("TestLogger", func(t *testing.T) {
		Log.Infoln("This is an info message")
	})
}

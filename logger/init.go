package logger

import "github.com/Moonlight-Companies/gologger/coloransi"

var Log *Logger = NewLogger(coloransi.Color(coloransi.BrightWhite, coloransi.Blue, "global"), WithLevel(LogLevelDebug))

package sectionlogger

import (
	"fmt"
	"testing"
	"time"

	"github.com/Moonlight-Companies/gologger/coloransi"
)

func TestSectionLoggerTest(t *testing.T) {
	logger := New("9864cf59cf 22 secs p32::s95::kPLACE::cBAGS")

	// Add events
	logger.Event("PLC_DIVERT", "🚀🚀 OnItemPrepare 🚀🚀 Traffic3 --> Exit Robot 1 table b in 0 via SecondDivert")
	time.Sleep(25 * time.Millisecond)
	logger.Event("PLC_DIVERT", "🚀🚀 ZZzzZZ 🚀🚀 25 MS DELAY")
	time.Sleep(25 * time.Millisecond)
	logger.Event("PLC_DIVERT", "🚀🚀 ZZzzZZ 🚀🚀 25 MS DELAY")

	// Add labeled entries
	logger.Add("IDENTIFICATION", "Cookie", "9864cf59cf")
	logger.Add("IDENTIFICATION", "Code", "21098744")

	// Set colors if desired
	logger.SetBorderColor(coloransi.BrightWhite, coloransi.BrightCyan, 0)
	logger.SetPrefixColor(coloransi.BrightWhite, 0, 0)

	// Render and output
	lines := logger.Render()
	for _, line := range lines {
		fmt.Println(line)
	}
}

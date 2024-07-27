package coloransi

import (
	"fmt"
	"math/rand"
	"strings"
)

// ColorCode represents ANSI color codes and RGB colors as a 32-bit integer.
// The lower 8 bits represent ANSI color codes, and the upper 24 bits represent RGB values.
type ColorCode uint32

// ANSI color codes
const (
	Black   ColorCode = 30
	Red     ColorCode = 31
	Green   ColorCode = 32
	Yellow  ColorCode = 33
	Blue    ColorCode = 34
	Magenta ColorCode = 35
	Cyan    ColorCode = 36
	White   ColorCode = 37

	// For bright colors, add 60
	BrightBlack   ColorCode = Black + 60
	BrightRed     ColorCode = Red + 60
	BrightGreen   ColorCode = Green + 60
	BrightYellow  ColorCode = Yellow + 60
	BrightBlue    ColorCode = Blue + 60
	BrightMagenta ColorCode = Magenta + 60
	BrightCyan    ColorCode = Cyan + 60
	BrightWhite   ColorCode = White + 60

	// Background colors start at 40, bright background colors at 100
	BackgroundOffset       ColorCode = 10
	BrightBackgroundOffset ColorCode = 60

	// RGB color mask
	RGBMask ColorCode = 0xFFFFFF00
)

// Additional static RGB color definitions
func CreateRGB(r, g, b uint8) ColorCode {
	return ColorCode(uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8)
}

// Additional static RGB color definitions
var ColorOrange ColorCode = CreateRGB(255, 140, 0)
var ColorPink ColorCode = CreateRGB(255, 192, 203)
var ColorPurple ColorCode = CreateRGB(128, 0, 128)
var ColorTeal ColorCode = CreateRGB(0, 128, 128)
var ColorLimeGreen ColorCode = CreateRGB(50, 205, 50)
var ColorIndigo ColorCode = CreateRGB(75, 0, 130)

// RGB creates a ColorCode from RGB values
func RGB(r, g, b uint8) ColorCode {
	return ColorCode(uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8)
}

// IsRGB checks if the ColorCode represents an RGB color
func (c ColorCode) IsRGB() bool {
	return c&RGBMask != 0
}

// ColorChooseRandom returns a random color code (non-black, including RGB colors).
func ColorChooseRandom() ColorCode {
	colors := []ColorCode{
		Red, Green, Yellow, Blue, Magenta, Cyan, White,
		BrightRed, BrightGreen, BrightYellow, BrightBlue, BrightMagenta, BrightCyan, BrightWhite,
		ColorOrange, ColorPink, ColorPurple, ColorTeal, ColorLimeGreen, ColorIndigo,
	}
	return colors[rand.Intn(len(colors))]
}

func ColorFrom(item uint64) ColorCode {
	colors := []ColorCode{
		Red,
		Green,
		Yellow,
		Blue,
		Magenta,
		Cyan,
		White,
		BrightRed,
		BrightGreen,
		BrightYellow,
		BrightBlue,
		BrightMagenta,
		BrightCyan,
		BrightWhite,
	}

	// Use the item value to deterministically select a color
	index := uint64(item) % uint64(len(colors))
	return colors[index]
}

// Color formats the given text with the specified foreground and background colors.
func Color(fg, bg ColorCode, v ...interface{}) string {
	fgCode := OneForeground(fg)
	bgCode := OneBackground(bg)
	reset := Reset()
	args := make([]string, len(v))
	for i, arg := range v {
		args[i] = fmt.Sprint(arg)
	}
	text := strings.Join(args, " ")
	return fmt.Sprintf("%s%s%s%s", fgCode, bgCode, text, reset)
}

func Foreground(fg ColorCode, v ...interface{}) string {
	fgCode := OneForeground(fg)
	reset := Reset()
	args := make([]string, len(v))
	for i, arg := range v {
		args[i] = fmt.Sprint(arg)
	}
	text := strings.Join(args, " ")
	return fmt.Sprintf("%s%s%s", fgCode, text, reset)
}

// OneForeground returns the ANSI escape sequence for the given color code.
func OneForeground(code ColorCode) string {
	if code.IsRGB() {
		r := (code >> 24) & 0xFF
		g := (code >> 16) & 0xFF
		b := (code >> 8) & 0xFF
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\033[%dm", code)
}

// OneBackground returns the ANSI escape sequence for the given background color code.
func OneBackground(code ColorCode) string {
	if code.IsRGB() {
		r := (code >> 24) & 0xFF
		g := (code >> 16) & 0xFF
		b := (code >> 8) & 0xFF
		return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\033[%dm", code+BackgroundOffset)
}

// Reset returns the ANSI escape sequence to reset the text color.
func Reset() string {
	return "\033[0m"
}

package coloransi

import (
	"fmt"
	"testing"
)

func TestTextStyles(t *testing.T) {
	testCases := []struct {
		name     string
		style    TextStyle
		expected int
	}{
		{"Bold", Bold, 1},
		{"Dim", Dim, 2},
		{"Italic", Italic, 3},
		{"Underline", Underline, 4},
		{"Blink", Blink, 5},
		{"FastBlink", FastBlink, 6},
		{"Reverse", Reverse, 7},
		{"Hidden", Hidden, 8},
		{"Strike", Strike, 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if int(tc.style) != tc.expected {
				t.Errorf("Expected %s to be %d, but got %d", tc.name, tc.expected, tc.style)
			}
		})
	}
}

func TestStyle(t *testing.T) {
	testCases := []struct {
		name     string
		style    TextStyle
		text     string
		expected string
	}{
		{
			name:     "Bold Text",
			style:    Bold,
			text:     "Bold",
			expected: "\033[1mBold\033[0m",
		},
		{
			name:     "Underlined Text",
			style:    Underline,
			text:     "Underlined",
			expected: "\033[4mUnderlined\033[0m",
		},
		{
			name:     "Strike Through",
			style:    Strike,
			text:     "Strikethrough",
			expected: "\033[9mStrikethrough\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Style(tc.style, tc.text)
			if result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestColorAndStyle(t *testing.T) {
	testCases := []struct {
		name     string
		fg       ColorCode
		bg       ColorCode
		style    TextStyle
		text     string
		expected string
	}{
		{
			name:     "Bold Red on Black",
			fg:       Red,
			bg:       Black,
			style:    Bold,
			text:     "Test",
			expected: "\033[31m\033[40m\033[1mTest\033[0m",
		},
		{
			name:     "Underlined Green on White",
			fg:       Green,
			bg:       White,
			style:    Underline,
			text:     "Hello",
			expected: "\033[32m\033[47m\033[4mHello\033[0m",
		},
		{
			name:     "Italic RGB on Black",
			fg:       CreateRGB(100, 150, 200),
			bg:       Black,
			style:    Italic,
			text:     "RGB",
			expected: "\033[38;2;100;150;200m\033[40m\033[3mRGB\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ColorAndStyle(tc.fg, tc.bg, tc.style, tc.text)
			if result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestVisual(t *testing.T) {
	// This test is for visual inspection and doesn't automatically check the output
	fmt.Println("\nVisual Style Test (inspect manually):")
	fmt.Println(Style(Bold, "Bold Text"))
	fmt.Println(Style(Dim, "Dim Text"))
	fmt.Println(Style(Italic, "Italic Text"))
	fmt.Println(Style(Underline, "Underlined Text"))
	fmt.Println(Style(Blink, "Blinking Text"))
	fmt.Println(Style(Strike, "Strikethrough Text"))
	fmt.Println(Style(Reverse, "Reversed Text"))

	fmt.Println("\nCombined Color and Style Test:")
	fmt.Println(ColorAndStyle(Red, Black, Bold, "Bold Red on Black"))
	fmt.Println(ColorAndStyle(Green, White, Underline, "Underlined Green on White"))
	fmt.Println(ColorAndStyle(Blue, Yellow, Italic, "Italic Blue on Yellow"))
	fmt.Println(ColorAndStyle(CreateRGB(100, 150, 200), Black, Bold, "Bold RGB on Black"))
}

func TestColorCodes(t *testing.T) {
	testCases := []struct {
		name     string
		color    ColorCode
		expected int
	}{
		{"Black", Black, 30},
		{"Red", Red, 31},
		{"Green", Green, 32},
		{"Yellow", Yellow, 33},
		{"Blue", Blue, 34},
		{"Magenta", Magenta, 35},
		{"Cyan", Cyan, 36},
		{"White", White, 37},
		{"BrightBlack", BrightBlack, 90},
		{"BrightRed", BrightRed, 91},
		{"BrightGreen", BrightGreen, 92},
		{"BrightYellow", BrightYellow, 93},
		{"BrightBlue", BrightBlue, 94},
		{"BrightMagenta", BrightMagenta, 95},
		{"BrightCyan", BrightCyan, 96},
		{"BrightWhite", BrightWhite, 97},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if int(tc.color) != tc.expected {
				t.Errorf("Expected %s to be %d, but got %d", tc.name, tc.expected, tc.color)
			}
		})
	}
}

func TestBackgroundOffset(t *testing.T) {
	if BackgroundOffset != 10 {
		t.Errorf("Expected BackgroundOffset to be 10, but got %d", BackgroundOffset)
	}
}

func TestBrightBackgroundOffset(t *testing.T) {
	if BrightBackgroundOffset != 60 {
		t.Errorf("Expected BrightBackgroundOffset to be 60, but got %d", BrightBackgroundOffset)
	}
}

func TestColor(t *testing.T) {
	testCases := []struct {
		name     string
		fg       ColorCode
		bg       ColorCode
		text     string
		expected string
	}{
		{
			name:     "Red on Black",
			fg:       Red,
			bg:       Black,
			text:     "Test",
			expected: "\033[31m\033[40mTest\033[0m",
		},
		{
			name:     "Green on White",
			fg:       Green,
			bg:       White,
			text:     "Hello",
			expected: "\033[32m\033[47mHello\033[0m",
		},
		{
			name:     "BrightBlue on BrightYellow",
			fg:       BrightBlue,
			bg:       BrightYellow,
			text:     "Bright",
			expected: "\033[94m\033[103mBright\033[0m",
		},
		{
			name:     "RGB Foreground on Black",
			fg:       CreateRGB(100, 150, 200),
			bg:       Black,
			text:     "RGB",
			expected: "\033[38;2;100;150;200m\033[40mRGB\033[0m",
		},
		{
			name:     "Red on RGB Background",
			fg:       Red,
			bg:       CreateRGB(200, 225, 255),
			text:     "RGB BG",
			expected: "\033[31m\033[48;2;200;225;255mRGB BG\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Color(tc.fg, tc.bg, tc.text)
			if result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestForeground(t *testing.T) {
	testCases := []struct {
		name     string
		fg       ColorCode
		text     string
		expected string
	}{
		{
			name:     "Red Foreground",
			fg:       Red,
			text:     "Red Text",
			expected: "\033[31mRed Text\033[0m",
		},
		{
			name:     "RGB Foreground",
			fg:       CreateRGB(100, 150, 200),
			text:     "RGB Text",
			expected: "\033[38;2;100;150;200mRGB Text\033[0m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := Foreground(tc.fg, tc.text)
			if result != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestCreateRGB(t *testing.T) {
	testCases := []struct {
		name     string
		r, g, b  uint8
		expected ColorCode
	}{
		{"Red", 255, 0, 0, ColorCode(0xFF000000)},
		{"Green", 0, 255, 0, ColorCode(0x00FF0000)},
		{"Blue", 0, 0, 255, ColorCode(0x0000FF00)},
		{"Mixed", 100, 150, 200, ColorCode(0x6496C800)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CreateRGB(tc.r, tc.g, tc.b)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestIsRGB(t *testing.T) {
	testCases := []struct {
		name     string
		color    ColorCode
		expected bool
	}{
		{"ANSI Red", Red, false},
		{"RGB Color", CreateRGB(100, 150, 200), true},
		{"ANSI BrightBlue", BrightBlue, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.color.IsRGB()
			if result != tc.expected {
				t.Errorf("Expected IsRGB() to return %v for %s, but got %v", tc.expected, tc.name, result)
			}
		})
	}
}

func TestColorVisual(t *testing.T) {
	// This test is for visual inspection and doesn't automatically check the output
	fmt.Println("Visual Color Test (inspect manually):")
	fmt.Println(Color(Red, Black, "Red on Black"))
	fmt.Println(Color(Green, White, "Green on White"))
	fmt.Println(Color(Blue, Yellow, "Blue on Yellow"))
	fmt.Println(Color(BrightCyan, BrightMagenta, "Bright Cyan on Bright Magenta"))
	fmt.Println(Color(CreateRGB(100, 150, 200), Black, "RGB (100, 150, 200) on Black"))
	fmt.Println(Color(White, CreateRGB(200, 100, 50), "White on RGB (200, 100, 50)"))
	fmt.Println(Foreground(ColorOrange, "Orange Text"))
	fmt.Println(Foreground(ColorPink, "Pink Text"))
	fmt.Println(Foreground(ColorPurple, "Purple Text"))
	fmt.Println(Foreground(ColorTeal, "Teal Text"))
	fmt.Println(Foreground(ColorLimeGreen, "Lime Green Text"))
	fmt.Println(Foreground(ColorIndigo, "Indigo Text"))
}

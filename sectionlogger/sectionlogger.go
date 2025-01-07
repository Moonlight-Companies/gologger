package sectionlogger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Moonlight-Companies/gologger/coloransi"
	. "github.com/Moonlight-Companies/gologger/coloransi"
)

// logEntry defines the interface for different types of log entries
type logEntry interface {
	Width() int          // Returns the width needed for the label/timestamp part
	Section() string     // Returns the section this entry belongs to
	RenderLabel() string // Returns the formatted label part
	RenderValue() string // Returns the formatted value part
}

// labelEntry implements logEntry for label-based entries
type labelEntry struct {
	section string
	label   string
	value   string
}

func (e *labelEntry) Width() int {
	return len(e.label)
}

func (e *labelEntry) Section() string {
	return e.section
}

func (e *labelEntry) RenderLabel() string {
	return e.label
}

func (e *labelEntry) RenderValue() string {
	return e.value
}

// eventEntry implements logEntry for timestamp-based events
type eventEntry struct {
	section   string
	value     string
	timestamp time.Time
	startTime time.Time
	prevTime  time.Time
}

func (e *eventEntry) Width() int {
	return 13 // Width for "mmmmmm mmmmmm" milliseconds
}

func (e *eventEntry) Section() string {
	return strings.ToUpper(e.section)
}

func (e *eventEntry) RenderLabel() string {
	sinceStart := e.timestamp.Sub(e.startTime).Milliseconds()
	sincePrev := e.timestamp.Sub(e.prevTime).Milliseconds()

	return fmt.Sprintf("%-6d %-6d", sinceStart, sincePrev)
}

func (e *eventEntry) RenderValue() string {
	return e.value
}

// SectionLogger manages the logging of sections and entries
type SectionLogger struct {
	mu          sync.RWMutex
	prefix      string
	entries     []logEntry
	startTime   time.Time
	lastTime    time.Time
	width       int
	borderFg    ColorCode
	borderBg    ColorCode
	borderStyle TextStyle
	prefixFg    ColorCode
	prefixBg    ColorCode
	prefixStyle TextStyle
	complete    bool
}

func New(prefix string) *SectionLogger {
	now := time.Now()
	return &SectionLogger{
		prefix:      prefix,
		entries:     make([]logEntry, 0),
		startTime:   now,
		lastTime:    now,
		width:       80,
		borderFg:    Black,
		borderBg:    BrightGreen,
		borderStyle: 0,
		prefixFg:    Black,
		prefixBg:    Cyan,
		prefixStyle: 0,
		complete:    false,
	}
}

func (l *SectionLogger) TimeSinceLast() time.Duration {
	return time.Since(l.lastTime)
}

func (l *SectionLogger) SetWidth(width int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.width = width
}

func (l *SectionLogger) renderEntry(entry logEntry, maxKeyWidth int) string {
	labelPadded := fmt.Sprintf("%-*s", maxKeyWidth, entry.Section()+"/"+entry.RenderLabel())

	return fmt.Sprintf("%s %s : %s",
		ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
		Foreground(BrightWhite, labelPadded),
		entry.RenderValue())
}

func (l *SectionLogger) Add(section, label, value string) {
	section = strings.ToUpper(section)

	entry := &labelEntry{
		section: section,
		label:   label,
		value:   value,
	}

	l.mu.Lock()
	l.entries = append(l.entries, entry)
	l.lastTime = time.Now()
	complete := l.complete
	l.mu.Unlock()

	if complete {
		l.mu.RLock()
		// For late messages, print a mini-section if it's the first late message for this section
		isNewSection := true
		for i := len(l.entries) - 2; i >= 0; i-- {
			if l.entries[i].Section() == section {
				isNewSection = false
				break
			}
		}
		l.mu.RUnlock()

		if isNewSection {
			// Print a continuation header for the new section
			fmt.Printf("%s %s\n",
				ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
				l.createBorderLine("┣", section))
		}

		// Print the entry directly
		fmt.Println(l.renderEntry(entry, l.calculateMaxWidth()))
	}
}

func (l *SectionLogger) Clear(section string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	section = strings.ToUpper(section)

	var newEntries []logEntry
	for _, entry := range l.entries {
		if entry.Section() != section {
			newEntries = append(newEntries, entry)
		}
	}
	l.entries = newEntries
}

func (l *SectionLogger) Event(section, value string) {
	now := time.Now()

	l.mu.Lock()
	entry := &eventEntry{
		section:   section,
		value:     value,
		timestamp: now,
		startTime: l.startTime,
		prevTime:  l.lastTime,
	}
	l.entries = append(l.entries, entry)
	l.lastTime = now
	complete := l.complete
	l.mu.Unlock()

	if complete {
		l.mu.RLock()
		// For late messages, print a mini-section if it's the first late message for this section
		isNewSection := true
		for i := len(l.entries) - 2; i >= 0; i-- {
			if l.entries[i].Section() == section {
				isNewSection = false
				break
			}
		}
		l.mu.RUnlock()

		if isNewSection {
			// Print a continuation header for the new section
			fmt.Printf("%s %s\n",
				ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
				l.createBorderLine("┣", section))
		}

		// Print the entry directly
		fmt.Println(l.renderEntry(entry, l.calculateMaxWidth()))
	}
}

func (l *SectionLogger) SetBorderColor(fg, bg ColorCode, style TextStyle) {
	l.borderFg = fg
	l.borderBg = bg
	l.borderStyle = style
}

func (l *SectionLogger) SetPrefixColor(fg, bg ColorCode, style TextStyle) {
	l.prefixFg = fg
	l.prefixBg = bg
	l.prefixStyle = style
}

func (l *SectionLogger) calculateMaxWidth() int {
	l.mu.RLock()
	defer l.mu.RUnlock()

	maxWidth := 0
	for _, entry := range l.entries {
		if width := entry.Width(); width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func (l *SectionLogger) createBorderLine(leftChar string, text string) string {
	l.mu.RLock()
	width := l.width
	borderFg := l.borderFg
	borderBg := l.borderBg
	borderStyle := l.borderStyle
	l.mu.RUnlock()

	parts := []string{
		ColorAndStyle(borderFg, borderBg, borderStyle, leftChar),
		ColorAndStyle(borderFg, borderBg, borderStyle, strings.Repeat("━", 4)),
	}

	borderWidth := width - 5 // Account for the left char and 4 dashes already added

	if text != "" {
		textColored := ColorAndStyle(BrightWhite, Black, borderStyle|coloransi.Bold, " "+text+" ")
		textSegment := fmt.Sprintf("%s%s%s",
			ColorAndStyle(borderFg, borderBg, borderStyle, "┥"),
			textColored,
			ColorAndStyle(borderFg, borderBg, borderStyle, "┝"))
		parts = append(parts, textSegment)

		borderWidth -= len(text) + 4
	}

	parts = append(parts, ColorAndStyle(borderFg, borderBg, borderStyle, strings.Repeat("━", borderWidth)))
	parts = append(parts, ColorAndStyle(borderFg, borderBg, borderStyle, "┅"))

	return strings.Join(parts, "")
}

func (l *SectionLogger) Render() []string {
	var results []string

	sections := make(map[string][]logEntry)
	var sectionOrder []string
	l.mu.RLock()
	l.complete = true
	for _, entry := range l.entries {
		section := entry.Section()
		if len(sections[section]) == 0 {
			sectionOrder = append(sectionOrder, section)
		}
		sections[section] = append(sections[section], entry)
	}
	l.mu.RUnlock()

	maxKeyWidth := l.calculateMaxWidth()

	for i, section := range sectionOrder {
		entries := sections[section]

		borderChar := "┣"
		if i == 0 {
			borderChar = "┏"
		}

		headerLine := fmt.Sprintf("%s %s",
			ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
			l.createBorderLine(borderChar, section))
		results = append(results, headerLine)

		for _, entry := range entries {
			labelPadded := fmt.Sprintf("%-*s", maxKeyWidth, entry.RenderLabel())

			line := fmt.Sprintf("%s %s %s : %s",
				ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
				ColorAndStyle(l.borderFg, l.borderBg, l.borderStyle, "┃"),
				Foreground(BrightWhite, labelPadded),
				entry.RenderValue())

			results = append(results, line)
		}
	}

	if len(results) > 0 {
		results = append(results, fmt.Sprintf("%s %s",
			ColorAndStyle(l.prefixFg, l.prefixBg, l.prefixStyle, l.prefix),
			l.createBorderLine("┗", "")))
	}

	return results
}

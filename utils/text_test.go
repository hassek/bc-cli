package utils

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected string
	}{
		{
			name:     "Short text - no wrapping needed",
			text:     "Hello world",
			width:    20,
			expected: "Hello world",
		},
		{
			name:     "Long text - simple wrapping",
			text:     "This is a long line that needs to be wrapped",
			width:    20,
			expected: "This is a long line\nthat needs to be\nwrapped",
		},
		{
			name:     "Text with existing line breaks",
			text:     "Line one\nLine two\nLine three",
			width:    20,
			expected: "Line one\nLine two\nLine three",
		},
		{
			name:     "Empty text",
			text:     "",
			width:    20,
			expected: "",
		},
		{
			name:     "Single long word - break it up",
			text:     "Supercalifragilisticexpialidocious",
			width:    10,
			expected: "Supercalif\nragilistic\nexpialidoc\nious",
		},
		{
			name:     "Mixed long and short words",
			text:     "This verylongwordthatneedstobebrokenuphere is text",
			width:    15,
			expected: "This\nverylongwordtha\ntneedstobebroke\nnuphere is text",
		},
		{
			name:     "Preserve empty lines",
			text:     "Paragraph one\n\nParagraph two",
			width:    20,
			expected: "Paragraph one\n\nParagraph two",
		},
		{
			name:     "Width of zero defaults to 80",
			text:     "Short text",
			width:    0,
			expected: "Short text",
		},
		{
			name:     "Multiple spaces between words",
			text:     "Word1    Word2     Word3",
			width:    20,
			expected: "Word1 Word2 Word3",
		},
		{
			name:     "Long description example",
			text:     "This is a very long product description that would overflow the terminal width if not wrapped properly, containing detailed information about the coffee beans.",
			width:    40,
			expected: "This is a very long product description\nthat would overflow the terminal width\nif not wrapped properly, containing\ndetailed information about the coffee\nbeans.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapText(tt.text, tt.width)
			if result != tt.expected {
				t.Errorf("WrapText() failed\nInput: %q\nWidth: %d\nExpected:\n%s\nGot:\n%s",
					tt.text, tt.width, tt.expected, result)
			}

			// Verify no line exceeds width (except for words that are longer than width)
			// Use effective width (default to DefaultTerminalWidth if width is 0)
			effectiveWidth := tt.width
			if effectiveWidth <= 0 {
				effectiveWidth = DefaultTerminalWidth
			}
			lines := strings.Split(result, "\n")
			for i, line := range lines {
				if len(line) > effectiveWidth {
					// Check if it's a single long word (no spaces)
					if strings.Contains(line, " ") {
						t.Errorf("Line %d exceeds width %d: %q (len=%d)", i, effectiveWidth, line, len(line))
					}
				}
			}
		})
	}
}

func TestWrapTextWithIndent(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		indent   string
		expected string
	}{
		{
			name:     "Simple indent",
			text:     "This is a long line that needs wrapping",
			width:    30,
			indent:   "  ",
			expected: "This is a long line that\n  needs wrapping",
		},
		{
			name:     "Multiple lines with indent",
			text:     "First line here and it is long enough to wrap\nSecond line",
			width:    30,
			indent:   "    ",
			expected: "First line here and it is\n    long enough to wrap\n    Second line",
		},
		{
			name:     "Empty text",
			text:     "",
			width:    20,
			indent:   "  ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapTextWithIndent(tt.text, tt.width, tt.indent)
			if result != tt.expected {
				t.Errorf("WrapTextWithIndent() failed\nExpected:\n%s\nGot:\n%s",
					tt.expected, result)
			}
		})
	}
}

func TestGetTerminalWidth(t *testing.T) {
	// This test just ensures the function doesn't panic
	// Actual width depends on the terminal running the test
	width := GetTerminalWidth()
	if width < MinTerminalWidth {
		t.Errorf("GetTerminalWidth() returned %d, which is less than MinTerminalWidth %d",
			width, MinTerminalWidth)
	}
}

func TestBreakWord(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		width    int
		expected []string
	}{
		{
			name:     "Short word",
			word:     "hello",
			width:    10,
			expected: []string{"hello"},
		},
		{
			name:     "Exact width",
			word:     "helloworld",
			width:    10,
			expected: []string{"helloworld"},
		},
		{
			name:     "Break into chunks",
			word:     "verylongword",
			width:    5,
			expected: []string{"veryl", "ongwo", "rd"},
		},
		{
			name:     "Single character width",
			word:     "hello",
			width:    1,
			expected: []string{"h", "e", "l", "l", "o"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := breakWord(tt.word, tt.width)
			if len(result) != len(tt.expected) {
				t.Errorf("breakWord() returned %d chunks, expected %d", len(result), len(tt.expected))
				return
			}
			for i, chunk := range result {
				if chunk != tt.expected[i] {
					t.Errorf("breakWord() chunk %d = %q, expected %q", i, chunk, tt.expected[i])
				}
			}
		})
	}
}

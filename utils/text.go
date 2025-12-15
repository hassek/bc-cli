package utils

import (
	"os"
	"strings"
	"unicode"

	"golang.org/x/term"
)

const (
	// DefaultTerminalWidth is used when terminal width cannot be detected
	DefaultTerminalWidth = 80
	// MinTerminalWidth is the minimum width we'll use for wrapping
	MinTerminalWidth = 40
)

// GetTerminalWidth returns the current terminal width
// Falls back to DefaultTerminalWidth if detection fails
func GetTerminalWidth() int {
	fd := int(os.Stdout.Fd())
	width, _, err := term.GetSize(fd)
	if err != nil || width < MinTerminalWidth {
		return DefaultTerminalWidth
	}
	return width
}

// WrapText wraps text to fit within the specified width
// Preserves existing line breaks and handles word boundaries
func WrapText(text string, width int) string {
	if width <= 0 {
		width = DefaultTerminalWidth
	}

	// Handle empty text
	if text == "" {
		return ""
	}

	// Split by existing line breaks
	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		// If line is empty or just whitespace, preserve it
		if strings.TrimSpace(line) == "" {
			wrappedLines = append(wrappedLines, "")
			continue
		}

		// Wrap this line
		wrapped := wrapLine(line, width)
		wrappedLines = append(wrappedLines, wrapped...)
	}

	return strings.Join(wrappedLines, "\n")
}

// wrapLine wraps a single line of text to the specified width
func wrapLine(line string, width int) []string {
	// Trim leading/trailing whitespace from the line
	line = strings.TrimSpace(line)

	if len(line) <= width {
		return []string{line}
	}

	var result []string
	var currentLine strings.Builder
	currentWidth := 0

	words := strings.FieldsFunc(line, unicode.IsSpace)

	for i, word := range words {
		wordLen := len(word)

		// If this is the first word on the line
		if currentWidth == 0 {
			// If word is longer than width, break it up
			if wordLen > width {
				chunks := breakWord(word, width)
				result = append(result, chunks[:len(chunks)-1]...)
				currentLine.WriteString(chunks[len(chunks)-1])
				currentWidth = len(chunks[len(chunks)-1])
			} else {
				currentLine.WriteString(word)
				currentWidth = wordLen
			}
		} else {
			// Check if adding this word (with a space) would exceed width
			spaceNeeded := 1 + wordLen
			if currentWidth+spaceNeeded > width {
				// Save current line and start new one
				result = append(result, currentLine.String())
				currentLine.Reset()

				// Handle long words that need breaking
				if wordLen > width {
					chunks := breakWord(word, width)
					result = append(result, chunks[:len(chunks)-1]...)
					currentLine.WriteString(chunks[len(chunks)-1])
					currentWidth = len(chunks[len(chunks)-1])
				} else {
					currentLine.WriteString(word)
					currentWidth = wordLen
				}
			} else {
				// Add word to current line
				currentLine.WriteString(" ")
				currentLine.WriteString(word)
				currentWidth += spaceNeeded
			}
		}

		// Last word - save the line
		if i == len(words)-1 && currentLine.Len() > 0 {
			result = append(result, currentLine.String())
		}
	}

	return result
}

// breakWord breaks a long word into chunks of the specified width
func breakWord(word string, width int) []string {
	if len(word) <= width {
		return []string{word}
	}

	var chunks []string
	for len(word) > width {
		chunks = append(chunks, word[:width])
		word = word[width:]
	}
	if word != "" {
		chunks = append(chunks, word)
	}
	return chunks
}

// WrapTextWithIndent wraps text and indents all lines except the first
func WrapTextWithIndent(text string, width int, indent string) string {
	wrapped := WrapText(text, width-len(indent))
	lines := strings.Split(wrapped, "\n")

	for i := 1; i < len(lines); i++ {
		lines[i] = indent + lines[i]
	}

	return strings.Join(lines, "\n")
}

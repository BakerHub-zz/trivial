package md

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Escape string start with ASCII punctuation character using backslash (\).
// Strings not start with ASCII punctuation character are returned without escape.
func Escape(txt string) string {
	r, _ := utf8.DecodeRuneInString(txt)
	if r < 127 && unicode.IsPunct(r) {
		return "\\" + txt
	}
	return txt
}

const (
	backtick = "`"
	asterisk = "*"
	underscore = "_"
)

var backtickRegex = regexp.MustCompile("`+")

func longestBacktickLength(s string) int {
	results := backtickRegex.FindAllString(s, -1)
	max := 0
	for _, backticks := range results {
		l := len(backticks)
		if l > max {
			max = l
		}
	}
	return max
}

// Code spans, wraps the input string with backtick (`).
func Code(txt string) string {
	n := longestBacktickLength(txt) + 1
	delimiters := strings.Repeat(backtick, n)
	sep := ""
	if strings.HasPrefix(txt, backtick) || strings.HasSuffix(txt, backtick) {
		sep = " "
	}
	return delimiters + sep + txt + sep + delimiters
}

func emphasis(txt string, n int) string {
	if len(txt) == 0 {
		return ""
	}
	delimiters := asterisk
	if strings.Contains(txt, asterisk) {
		delimiters = underscore
	}
	delimiters = strings.Repeat(delimiters, n)
	return delimiters + txt + delimiters
}

// Em emphasis by wrap text with * or _.
func Em(txt string) string {
	return emphasis(txt, 1)
}

// Strong emphasis by wrap text with ** or __.
func Strong(txt string) string {
	return emphasis(txt, 2)
}

// Link generates link with text and optional url and title.
func Link(txt string, urlAndTitle ...string) string {
	switch len(urlAndTitle) {
	case 0:
		return fmt.Sprintf("[%s]", txt)
	case 1:
		return fmt.Sprintf("[%s](%s)", txt, urlAndTitle[0])
	default:
		return fmt.Sprintf(`[%s](%s "%s")`, txt, urlAndTitle[0], strings.Join(urlAndTitle[1:], " "))
	}
}

// Link generates image with caption and optional url and title.
func Image(caption string, urlAndTitle ...string) string {
	return "!" + Link(caption, urlAndTitle...)
}

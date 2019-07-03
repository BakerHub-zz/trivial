package text

import (
	"strings"
	"unicode"
)

func IsChineseRune(c rune) bool {
	return unicode.Is(unicode.Han, c)
}

func ContainsChinese(s string) bool {
	for _, r := range s {
		if IsChineseRune(r) {
			return true
		}
	}
	return false
}

func textWidth(text string) int {
	w := 0
	for _, c := range text {
		if IsChineseRune(c) {
			w += 2
		} else {
			w += 1
		}
	}
	return w
}

type WordWrapper struct {
	limit              int
	used               int
	builder            *strings.Builder
	lastNeedsSeparator bool
}

func NewWordWrapper(width int) *WordWrapper {
	return &WordWrapper{limit: width, used: 0, builder: &strings.Builder{}, lastNeedsSeparator: false}
}

func (ww *WordWrapper) append(lines []string, text string, needSeparator bool) []string {
	tw := textWidth(text)
	sw := 1
	sep := " "
	if ww.used == 0 || !ww.lastNeedsSeparator && !needSeparator {
		sw = 0
		sep = ""
	}

	if ww.used+tw+sw <= ww.limit {
		ww.put(sep, sw)
		ww.put(text, tw)
	} else {
		lines = append(lines, ww.commit())
		ww.put(text, tw)
	}
	ww.lastNeedsSeparator = needSeparator
	return lines
}

func (ww *WordWrapper) put(text string, w int) {
	if w > 0 {
		ww.builder.WriteString(text)
		ww.used += w
	}
}

func (ww *WordWrapper) reset() {
	ww.builder.Reset()
	ww.used = 0
	ww.lastNeedsSeparator = false
}

func (ww *WordWrapper) commit() string {
	rv := ww.builder.String()
	ww.reset()
	return rv
}

func (ww *WordWrapper) Lines(text string) (lines []string) {
	ww.reset()

	for _, field := range strings.Fields(text) {
		if ContainsChinese(field) {
			for _, c := range field {
				lines = ww.append(lines, string(c), false)
			}
		} else {
			lines = ww.append(lines, field, true)
		}
	}

	if ww.used > 0 {
		lines = append(lines, ww.commit())
	}

	return lines
}

func Wordwrap(text string, limit int) []string {
	ww := NewWordWrapper(limit)
	return ww.Lines(text)
}

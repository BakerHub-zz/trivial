package md

import (
	"errors"
	"fmt"
	"github.com/BakerHub/trivial/text"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/BakerHub/trivial/check"
)

var ErrInvalidHeadingLevel = errors.New("head level should between [1..6]")

type List struct {
	ordered bool
	bullet  string
	index   int
}

func (l *List) NextBullet() string {
	l.index += 1
	if l.ordered {
		return strconv.Itoa(l.index) + "."
	}
	return l.bullet
}

type Writer struct {
	underlay   io.Writer
	indent     string
	lineLength int

	lists []*List

	lastLineEmpty       bool
	hasPendingEmptyLine bool
}

type Option func(*Writer)

func NewWriter(writer io.Writer, options ...Option) *Writer {
	wr := &Writer{
		underlay:   writer,
		indent:     "    ",
		lineLength: 80,
	}
	for _, opt := range options {
		opt(wr)
	}

	return wr
}

func Indent(indent string) Option {
	return func(w *Writer) {
		w.indent = indent
	}
}

func LineLength(limit int) Option {
	return func(w *Writer) {
		w.lineLength = limit
	}
}

func (wr *Writer) ensureEmptyLine() {
	if !wr.lastLineEmpty {
		wr.writeNewLine()
	}
}

func (wr *Writer) pendingEmptyLine() {
	wr.hasPendingEmptyLine = true
	wr.lastLineEmpty = false
}

func (wr *Writer) flushPendingEmptyLine() {
	if wr.hasPendingEmptyLine {
		wr.ensureEmptyLine()
		wr.hasPendingEmptyLine = false
	}
}

func (wr *Writer) mustWriteHeading(level int, text string) {
	if level < 1 || level > 6 {
		log.Panic(ErrInvalidHeadingLevel)
	}

	wr.flushPendingEmptyLine()

	_, err := fmt.Fprintf(wr.underlay, "%s %s\n", strings.Repeat("#", level), text)
	check.Check(err)

	wr.pendingEmptyLine()
}

func (wr *Writer) WriteH1(text string) {
	wr.mustWriteHeading(1, text)
}

func (wr *Writer) WriteH2(text string) {
	wr.mustWriteHeading(2, text)
}

func (wr *Writer) WriteH3(text string) {
	wr.mustWriteHeading(3, text)
}

func (wr *Writer) WriteH4(text string) {
	wr.mustWriteHeading(4, text)
}

func (wr *Writer) WriteH5(text string) {
	wr.mustWriteHeading(5, text)
}

func (wr *Writer) WriteH6(text string) {
	wr.mustWriteHeading(6, text)
}

func (wr *Writer) WriteParagraph(txt string) {
	lines := text.Wordwrap(txt, wr.lineLength)
	if len(lines) == 0 {
		return
	}

	wr.flushPendingEmptyLine()

	for _, line := range lines {
		_, err := fmt.Fprintln(wr.underlay, line)
		check.Check(err)
	}

	wr.pendingEmptyLine()
}

func (wr *Writer) writeIndented(txt string, firstLineIndent string, followingLineIndent string) {
	wrapWidth := wr.lineLength - len(firstLineIndent)
	lines := text.Wordwrap(txt, wrapWidth)

	for i, line := range lines {
		indent := followingLineIndent
		if i == 0 {
			indent = firstLineIndent
		}
		_, err := fmt.Fprintf(wr.underlay, "%s%s\n", indent, line)
		check.Check(err)
	}
	wr.lastLineEmpty = false
}

func (wr *Writer) BeginUnorderedList() {
	list := &List{false, "-", 1}
	wr.beginList(list)
}

func (wr *Writer) beginList(list *List) {
	wr.lists = append(wr.lists, list)
	if wr.listLevel() == 0 {
		wr.flushPendingEmptyLine()
	}
}

func (wr *Writer) EndList() {
	level := wr.listLevel()
	if level < 0 {
		log.Panic("no list to end")
	} else if level == 0 {
		wr.pendingEmptyLine()
	}

	wr.lists = wr.lists[:level]
}

func (wr *Writer) BeginOrderedList() {
	list := &List{true, "", 0}
	wr.beginList(list)
}

func (wr *Writer) listLevel() int {
	return len(wr.lists) - 1
}

func (wr *Writer) currentList() *List {
	if len(wr.lists) == 0 {
		return nil
	}

	return wr.lists[len(wr.lists)-1]
}

func (wr *Writer) listIndent() string {
	return strings.Repeat(wr.indent, wr.listLevel())
}

func (wr *Writer) WriteListItem(text string) {
	list := wr.currentList()
	if list == nil {
		log.Panic("list not begin")
	}

	indent := fmt.Sprintf("%s%s ", wr.listIndent(), list.NextBullet())
	followingIndent := strings.Repeat(" ", len(indent))
	wr.writeIndented(text, indent, followingIndent)
}

func (wr *Writer) writeNewLine() {
	_, err := fmt.Fprintln(wr.underlay)
	check.Check(err)
	wr.lastLineEmpty = true
}

func (wr *Writer) WriteHorizontalRule() {
	wr.flushPendingEmptyLine()

	_, err := fmt.Fprintln(wr.underlay, "* * *")
	check.Check(err)

	wr.pendingEmptyLine()
}

func (wr *Writer) WriteFencedCodeBlock(code string, language string) {
	wr.flushPendingEmptyLine()
	_, err := fmt.Fprintf(wr.underlay, "```%s\n%s\n```\n", language, code)
	check.Check(err)
	wr.pendingEmptyLine()
}

func (wr *Writer) WriteIndentedCodeBlock(code string) {
	wr.flushPendingEmptyLine()
	wr.writeIndented(code, wr.indent, wr.indent)
	wr.pendingEmptyLine()
}

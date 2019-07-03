package md

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

type MockedWriter struct {
	mock.Mock
}

func (m *MockedWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func TestNewWriter(t *testing.T) {
	mw := &MockedWriter{}
	w := NewWriter(mw)
	assert.Equal(t, mw, w.underlay)
}

func TestNewWriterWithLineLength(t *testing.T) {
	mw := &MockedWriter{}
	limit := 10
	lineLength := LineLength(limit)
	w := NewWriter(mw, lineLength)
	assert.Equal(t, limit, w.lineLength)
}

func TestNewWriterWithIndent(t *testing.T) {
	mw := &MockedWriter{}
	indent := "  "
	indentOption := Indent(indent)
	w := NewWriter(mw, indentOption)
	assert.Equal(t, indent, w.indent)
}

func TestWriter_WriteSuccess(t *testing.T) {
	cases := map[string]struct {
		lineLength int
		indent     string
		use        func(wr *Writer)
		want       string
	}{
		"H1": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH1("This is the title") },
			"# This is the title\n",
		},
		"H1 exceed line length": {
			5,
			"  ",
			func(wr *Writer) { wr.WriteH1("This is the title") },
			"# This is the title\n",
		},
		"H2": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH2("This is the title") },
			"## This is the title\n",
		},
		"H3": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH3("This is the title") },
			"### This is the title\n",
		},
		"H4": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH4("This is the title") },
			"#### This is the title\n",
		},
		"H5": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH5("This is the title") },
			"##### This is the title\n",
		},
		"H6": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteH6("This is the title") },
			"###### This is the title\n",
		},
		"Paragraph no wrap": {
			80,
			"  ",
			func(wr *Writer) { wr.WriteParagraph("This is the paragraph.") },
			"This is the paragraph.\n",
		},
		"Paragraph wrap": {
			10,
			"  ",
			func(wr *Writer) { wr.WriteParagraph("This is the paragraph.") },
			"This is\nthe\nparagraph.\n",
		},
		"HorizontalRule": {
			10,
			"  ",
			func(wr *Writer) { wr.WriteHorizontalRule() },
			"* * *\n",
		},
		"FencedCodeBlock": { // TODO: check for indent rule for code block
			10,
			"  ",
			func(wr *Writer) { wr.WriteFencedCodeBlock("const a = 'a';", "javascript") },
			"```javascript\nconst a = 'a';\n```\n",
		},
		"IndentedCodeBlock 2 space": { // TODO: check for indent rule for code block
			10,
			"  ",
			func(wr *Writer) { wr.WriteIndentedCodeBlock("const a = 'a';") },
			"  const a\n  = 'a';\n",
		},
		"IndentedCodeBlock 4 space": { // TODO: check for indent rule for code block
			10,
			"    ",
			func(wr *Writer) { wr.WriteIndentedCodeBlock("const a = 'a';") },
			"    const\n    a =\n    'a';\n",
		},
		"Ordered list": {
			lineLength: 80,
			indent:     "    ",
			use: func(wr *Writer) {
				wr.BeginOrderedList()
				wr.WriteListItem("first")
				wr.WriteListItem("second")
				wr.EndList()
			},
			want:
			`1. first
2. second
`,
		},
		"Unordered list": {
			lineLength: 80,
			indent:     "    ",
			use: func(wr *Writer) {
				wr.BeginUnorderedList()
				wr.WriteListItem("first")
				wr.WriteListItem("second")
				wr.EndList()
			},
			want:
			`- first
- second
`,
		},
		"Nesting list": {
			lineLength: 80,
			indent:     "    ",
			use: func(wr *Writer) {
				wr.BeginUnorderedList()
				wr.WriteListItem("first")
				wr.BeginOrderedList()
				wr.WriteListItem("one")
				wr.WriteListItem("two")
				wr.EndList()
				wr.WriteListItem("second")
				wr.EndList()
			},
			want:
			`- first
    1. one
    2. two
- second
`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			builder := &strings.Builder{}

			writer := NewWriter(builder, Indent(tc.indent), LineLength(tc.lineLength))
			tc.use(writer)

			assert.Equal(t, tc.want, builder.String())
		})
	}
}

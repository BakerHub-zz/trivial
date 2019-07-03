package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordwrap(t *testing.T) {
	cases := map[string]struct {
		limit int
		text  string
		want  []string
	}{
		"en":                  {10, "This is a text line.", []string{"This is a", "text line."}},
		"en punctuation wrap": {10, "This is an, text line.", []string{"This is", "an, text", "line."}},
		"cn":                  {10, "这是一个，测试", []string{"这是一个，", "测试"}},
		"cn punctuation":      {10, "这真是一个，测试", []string{"这真是一个", "，测试"}}, // Currently not support chinese punctuation.
		"mixed":               {10, "good 这是一个，测试", []string{"good 这是", "一个，测试"}},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Wordwrap(tc.text, tc.limit)
			assert.Equal(t, tc.want, got)
		})
	}
}

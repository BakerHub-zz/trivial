package md

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCode(t *testing.T) {
	cases := map[string]struct{
		input string
		want string
	}{
		"simple": {"some text", "`some text`"},
		"extra backtick": {"some `` test ``` text", "````some `` test ``` text````"},
		"leading backtick": {"``some", "``` ``some ```"},
		"tailing backtick": {"some`", "`` some` ``"},
		"empty": {"", "``"},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Code(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEscape(t *testing.T) {
	cases := map[string]struct{
		input string
		want string
	}{
		"escapable": {"[]", "\\[]"},
		"not escapable": {"simple", "simple"},
		"not escapable cn": {"「」", "「」"},
		"empty": {"", ""},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Escape(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestEm(t *testing.T) {
	cases := map[string]struct{
		input string
		want string
	}{
		"simple": {"text", "*text*"},
		"asterisk": {"simple * or not", "_simple * or not_"},
		"underscore": {"simple_or_not", "*simple_or_not*"},
		"empty": {"", ""},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Em(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestStrong(t *testing.T) {
	cases := map[string]struct{
		input string
		want string
	}{
		"simple": {"text", "**text**"},
		"asterisk": {"simple * or not", "__simple * or not__"},
		"underscore": {"simple_or_not", "**simple_or_not**"},
		"empty": {"", ""},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Strong(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}


func TestLink(t *testing.T) {
	cases := map[string]struct{
		input []string
		want string
	}{
		"link ref": {[]string{"text"}, "[text]"},
		"link with url": {[]string{"text", "url"}, "[text](url)"},
		"link with title": {[]string{"text", "url", "title"}, `[text](url "title")`},
		"link with more titles": {[]string{"text", "url", "title", "and", "more"}, `[text](url "title and more")`},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Link(tc.input[0], tc.input[1:]...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestImage(t *testing.T) {
	cases := map[string]struct{
		input []string
		want string
	}{
		"image ref": {[]string{"text"}, "![text]"},
		"image with url": {[]string{"text", "url"}, "![text](url)"},
		"image with title": {[]string{"text", "url", "title"}, `![text](url "title")`},
		"image with more titles": {[]string{"text", "url", "title", "and", "more"}, `![text](url "title and more")`},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := Image(tc.input[0], tc.input[1:]...)
			assert.Equal(t, tc.want, got)
		})
	}
}

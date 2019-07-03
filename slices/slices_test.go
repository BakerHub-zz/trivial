package slices

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnique(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  []string
	}{
		"no duplication": {
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		"duplication": {
			[]string{"a", "b", "c", "a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		"only same elements": {
			[]string{"a", "a", "a"},
			[]string{"a"},
		},
		"empty": {
			[]string{},
			nil,
		},
		"nil": {
			nil,
			nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Unique(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCompact(t *testing.T) {
	cases := map[string]struct {
		input []string
		want  []string
	}{
		"no empty": {
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		"has empty": {
			[]string{"a", "", "", "b", "", "c"},
			[]string{"a", "b", "c"},
		},
		"all empty": {
			[]string{"", "", "", "", "", "", "", ""},
			nil,
		},
		"nil": {
			nil,
			nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Compact(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMap(t *testing.T) {
	cases := map[string]struct {
		input []string
		fn    func(string) string
		want  []string
	}{
		"repeat": {
			[]string{"a", "b", "c"},
			func(s string) string { return strings.Repeat(s, 2) },
			[]string{"aa", "bb", "cc"},
		},
		"empty": {
			[]string{},
			func(s string) string { return strings.Repeat(s, 2) },
			nil,
		},
		"nil": {
			nil,
			nil,
			nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Map(tc.input, tc.fn)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFilter(t *testing.T) {
	cases := map[string]struct {
		input []string
		fn    func(string) bool
		want  []string
	}{
		"keep long": {
			[]string{"a", "aa", "b", "bb", "c", "cc"},
			func(s string) bool { return len(s) == 2 },
			[]string{"aa", "bb", "cc"},
		},
		"empty": {
			[]string{},
			func(s string) bool { return len(s) == 2 },
			nil,
		},
		"nil": {
			nil,
			nil,
			nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Filter(tc.input, tc.fn)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestContains(t *testing.T) {
	cases := map[string]struct {
		input []string
		value string
		want  bool
	}{
		"keep long": {
			[]string{"a", "aa", "b", "bb", "c", "cc"},
			"a",
			true,
		},
		"empty": {
			[]string{},
			"",
			false,
		},
		"nil": {
			nil,
			"",
			false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Contains(tc.input, tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}

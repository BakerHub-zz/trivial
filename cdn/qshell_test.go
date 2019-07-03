package cdn

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithIgnoreSuffixes(t *testing.T) {
	cases := map[string]struct {
		args []string

		want string
	}{
		"empty": {
			[]string{},
			".DS_Store,Thumbs.db",
		},
		"append": {
			[]string{"a", "b"},
			"a,b,.DS_Store,Thumbs.db",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			var q QShellUploader
			opt := IgnoreSuffixes(tc.args...)
			opt(&q)
			assert.Equal(t, tc.want, q.skipSuffixes)
		})
	}
}

func TestNewQShellUploader(t *testing.T) {
	NewQShellUploader("test")
	NewQShellUploader("test", IgnoreSuffixes(".exe"))
}

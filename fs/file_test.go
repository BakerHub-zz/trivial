package fs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExists(t *testing.T) {
	cases := map[string]struct {
		path string
		want bool
	}{
		"exist file": {"file_test.go", true},
		"exist dir":  {".", true},
		"not exist":  {"file-not-exist.go", false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := Exists(tc.path)
			assert.Equal(t, tc.want, got)
		})
	}
}

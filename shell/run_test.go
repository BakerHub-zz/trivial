package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunSuccess(t *testing.T) {
	Run("ls", ".")
}

func TestRunPanic(t *testing.T) {
	assert.Panics(t, func() {
		Run("false")
	}, "Run did not panic")
}

func TestShell_Run(t *testing.T) {
	s := &Shell{}
	s.Run("ls", ".")
}

package check

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckNoError(t *testing.T) {
	Check(nil)
}

func TestCheckShouldPanicWhenPassedError(t *testing.T) {
	assert.Panics(t, func() {
		Check(errors.New("text"))
	}, "Check did not panic")
}

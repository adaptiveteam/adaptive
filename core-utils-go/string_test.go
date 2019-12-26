package core_utils_go

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsString(t *testing.T) {
	list := []string{"hello", "world"}
	assert.True(t, ListContainsString(list, "hello"))

	assert.False(t, ListContainsString(list, "qwerty"))
}

func TestUuid(t *testing.T) {
	assert.True(t, len(Uuid()) == 36)
}

func TestTrimLower(t *testing.T) {
	str := " Hello, world!  "
	assert.True(t, TrimLower(str) == "hello, world!")
}

func TestTextWrapEmptyString(t *testing.T) {
	assert.True(t, TextWrap(EmptyString, Underscore) == EmptyString)
}

func TestTextWrapOneWrapper(t *testing.T) {
	str := "hello"
	assert.True(t, TextWrap(str, Underscore) == "_hello_")
}

func TestTextWrapMultipleWrappers(t *testing.T) {
	str := "hello"
	assert.True(t, TextWrap(str, Underscore, Asterisk) == "*_hello_*")
}


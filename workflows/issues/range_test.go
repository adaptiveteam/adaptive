package issues

import (
	"testing"
)

type S struct {
	a string
}
func TestRangeLoop(t *testing.T) {
	structs := []S{{a: "a"}}
	for _, s := range structs {
		s.a = "b"
	}
	var res string
	for _, s := range structs {
		res = res + s.a
	}

	if res == "b" {t.Fatalf("Unexpected res %s", res)}
}

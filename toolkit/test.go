package toolkit

import (
	"fmt"
	"testing"
)

func AssertEqual(x, y int, t *testing.T) {
	if x != y {
		t.Error(fmt.Sprintf("%d != %d", x, y))
	}
}

func AssertStrEqual(x, y string, t *testing.T) {
	if x != y {
		t.Error(fmt.Sprintf("%s != %s", x, y))
	}
}

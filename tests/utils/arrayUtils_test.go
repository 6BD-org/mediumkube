package utils

import (
	"mediumkube/utils"
	"testing"
)

func testContains(t *testing.T) {
	a := []string{"a", "B"}
	if !utils.Contains(a, "a") {
		t.Fail()
	}
}

package utils

import (
	"mediumkube/utils"
	"testing"
)

type A struct {
	Key string
	Val string
}

func testContains(t *testing.T) {
	a := []string{"a", "B"}
	if !utils.Contains(a, "a") {
		t.Fail()
	}

	if !utils.ContainsT(a, "a") {
		t.Fail()
	}

	objLst := []A{{Key: "A", Val: "B"}, {Key: "C", Val: "D"}}
	if !utils.ContainsT(objLst, A{Key: "A", Val: "B"}) {
		t.Fail()
	}
}

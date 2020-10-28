package utils

import (
	"mediumkube/utils"
	"testing"
)

func TestReadByte(t *testing.T) {
	bytes := utils.ReadByte("./test.txt")
	if bytes[len(bytes)-1] != '\n' {
		t.Fail()
	}
}

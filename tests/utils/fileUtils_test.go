package utils

import (
	"fmt"
	"mediumkube/utils"
	"os"
	"testing"
)

func TestReadByte(t *testing.T) {
	bytes := utils.ReadByte("./test.txt")
	if bytes[len(bytes)-1] != '\n' {
		t.Fail()
	}
}

func TestGetFileDir(t *testing.T) {
	absPath := "/abc/def/g.jpg"
	dir := utils.GetFileDir(absPath)
	if dir != "/abc/def" {
		t.Fail()
	}
}

func TestFileMode(t *testing.T) {
	fmt.Println(os.FileMode(0666).String())

}

package utils

import (
	"fmt"
	"mediumkube/utils"
	"os"
	"path/filepath"
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
	file := utils.GetFileName(absPath)
	if dir != "/abc/def" {
		t.Fail()
	}

	if file != "g.jpg" {
		t.Fail()
	}

	absPath = "/abc/a/b/"
	dir = utils.GetFileDir(absPath)
	file = utils.GetFileName(absPath)
	if dir != "/abc/a/b" {
		t.Fail()
	}

	if file != "" {
		t.Fail()
	}

	absPath = ""
	dir = utils.GetFileDir(absPath)
	file = utils.GetFileName(absPath)

	if file != "" {
		t.Fail()
	}

	if dir != "" {
		t.Fail()
	}

}

func TestFileMode(t *testing.T) {
	fmt.Println(os.FileMode(0666).String())
}

func TestWalkDir(t *testing.T) {
	files := utils.WalkDir("./walk_root")
	if len(files) != 2 {
		t.Fail()
	}

	files = utils.WalkDir("..")
	fmt.Println(files)

}

func TestGetDirName(t *testing.T) {
	dirName := utils.GetDirName("a/b/c/")
	if dirName != "c" {
		t.Fail()
	}

	dirName = utils.GetDirName("c")
	if dirName != "c" {
		t.Fail()
	}
	dirName = utils.GetDirName("/c")
	if dirName != "c" {
		t.Fail()
	}
	dirName = utils.GetDirName("a/b/c")
	if dirName != "c" {
		t.Fail()
	}

}

func TestPath(t *testing.T) {
	fmt.Println("On Unix:")
	fmt.Println(filepath.Base("/foo/bar/baz.js"))
	fmt.Println(filepath.Base("/foo/bar/baz"))
	fmt.Println(filepath.Base("/foo/bar/baz/"))
	fmt.Println(filepath.Base("dev.txt"))
	fmt.Println(filepath.Base("../todo.txt"))
	fmt.Println(filepath.Base(".."))
	fmt.Println(filepath.Base("."))
	fmt.Println(filepath.Base("/"))
	fmt.Println(filepath.Base(""))
}

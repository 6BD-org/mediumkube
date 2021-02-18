package utils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ReadStr read content as a string
func ReadStr(path string) string {
	data, err := ioutil.ReadFile(path)
	CheckErr(err)
	return string(data)
}

// ReadByte read a file as byte array
func ReadByte(path string) []byte {
	data, err := ioutil.ReadFile(path)
	CheckErr(err)
	return data
}

// WriteStr to file or die
func WriteStr(path string, content string, perm os.FileMode) {
	err := ioutil.WriteFile(path, []byte(content), perm)
	CheckErr(err)
}

// GetFileName get file name from its full path
func GetFileName(fullPath string) string {
	splitted := strings.Split(fullPath, "/")
	if len(splitted) == 0 {
		return ""
	}
	return splitted[len(splitted)-1]
}

// GetFileDir get dir of file given full path
func GetFileDir(fullPath string) string {
	lastSlash := strings.LastIndex(fullPath, "/")
	if lastSlash < 0 {
		return ""
	}
	return fullPath[:lastSlash]
}

// GetDirName get the name of a directory
// for example:
// the dir name of /a/b/c/ is c
// the dir name of /a/b/c is also c.
// NOTE: the argument must be a directory. If you use a file, you might get unexpected result
func GetDirName(fullPath string) string {

	if fullPath[len(fullPath)-1] == '/' {
		fullPath = fullPath[:len(fullPath)-1]
	}

	splitted := strings.Split(fullPath, "/")
	return splitted[len(splitted)-1]
}

// WalkDir list all files in directory
func WalkDir(path string) []string {
	files := make([]string, 0)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil && err != filepath.SkipDir {
			log.Panic(err)
		} else {
			var fi os.FileInfo
			fi, err = os.Stat(path)
			CheckErr(err)
			if !fi.IsDir() {
				files = append(files, path)
			}
		}

		return nil
	})

	return files
}

// Copy a file or die
func Copy(src string, tgt string) {
	srcFile, err := os.Open(src)
	CheckErr(err)
	tgtFile, err := os.Create(tgt)
	CheckErr(err)
	defer srcFile.Close()
	defer tgtFile.Close()

	_, err = io.Copy(tgtFile, srcFile)
	CheckErr(err)
}

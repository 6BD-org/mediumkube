package utils

import (
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

// GetFileName get file name from its full path
func GetFileName(fullPath string) string {
	splitted := strings.Split(fullPath, "/")
	return splitted[len(splitted)-1]
}

// GetFileDir get dir of file given full path
func GetFileDir(fullPath string) string {
	lastSlash := strings.LastIndex(fullPath, "/")
	return fullPath[:lastSlash]
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

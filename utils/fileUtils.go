package utils

import (
	"io/ioutil"
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

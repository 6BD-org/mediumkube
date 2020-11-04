package flogger

import (
	"bytes"
	"io"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// posMap records current end position
// of watched files
var posMap = make(map[string]int64)

// Keep track of file
var fileMap = make(map[string]*os.File)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

// initPosition gets init end position
// This is done before
func initPosition(f *os.File) int64 {
	// Print log from beginning of the file
	return 0
}

func handleEvent(e fsnotify.Event) {
	if e.Op&fsnotify.Write != fsnotify.Write {
		return
	}
	name := e.Name
	file, ok := fileMap[name]

	var err error
	if !ok {

		file, err = os.Open(name)
		panicOnErr(err)
		fileMap[name] = file
		posMap[name] = initPosition(file)
	}

	buffer := make([]byte, 512)
	var size int
	var builder bytes.Buffer
	// Read from init position until eof is encountered

	for {
		size, err = file.ReadAt(buffer, posMap[name])
		posMap[name] += int64(size)
		builder.Write(buffer[:size])
		if err == io.EOF || size == 0 {
			break
		}
	}
	// fmt.Println(size, string(buffer[:size]))

	logStr := builder.String()
	logStrs := strings.Split(logStr, "\n")

	for _, logStr = range logStrs {
		logStr = strings.TrimSpace(logStr)
		if len(logStr) > 0 {
			flog(name, logStr)
		}
	}

}

// FLog watch a list of file changes and log them on stdio
func FLog(filelist []string) {

	watcher, err := fsnotify.NewWatcher()
	panicOnErr(err)

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				handleEvent(event)
				// Process event here
			case _, ok := <-watcher.Errors:
				if !ok {
					return
				}
				// Process err
			}

		}

	}()

	for _, f := range filelist {
		watcher.Add(f)
	}

	for _, fptr := range fileMap {
		fptr.Close()
	}
	<-done
}

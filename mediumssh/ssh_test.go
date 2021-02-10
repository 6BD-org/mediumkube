package mediumssh

import (
	"log"
	"os"
	"testing"
)

func TestArrow(t *testing.T) {
	stdin := os.Stdin

	buf := make([]byte, 1024)
	for {

		n, _ := stdin.Read(buf)
		log.Println(buf[:n])
	}
}

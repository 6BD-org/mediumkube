package utils

import (
	"log"
	"testing"
)

func TestGen(t *testing.T) {
	m := GenerateMac()
	log.Println(m)
}

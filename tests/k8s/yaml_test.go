package k8s

import (
	"log"
	"mediumkube/pkg/k8s"
	"reflect"
	"testing"
)

func TestParseYaml(t *testing.T) {
	ch := make(chan interface{})
	go k8s.ParseResources("./test.yaml", ch)
	for v := range ch {
		log.Println(reflect.TypeOf(v))
	}
}

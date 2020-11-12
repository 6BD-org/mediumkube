package k8s

import (
	"log"
	"mediumkube/k8s"
	"reflect"
	"testing"
)

func TestParseYaml(t *testing.T) {
	resMap := k8s.ParseResources("./test.yaml")
	for k, v := range resMap {
		log.Println(k, v)
		log.Println(reflect.TypeOf(v))
	}
}

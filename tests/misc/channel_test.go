package misc

import (
	"fmt"
	"testing"
)

// func TestChannel(t *testing.T) {
// 	ch := make(chan int)
// 	ch <- 1
// 	fmt.Println(<-ch)
// 	go close(ch)
// 	for i := range ch {
// 		fmt.Println(i)
// 	}
// }

func TestType(t *testing.T) {
	var a interface{} = 1
	switch t := a.(type) {
	default:
		fmt.Println(t)
	}
}

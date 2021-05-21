package utils

import "fmt"

func EtcdEp(host string, port int) string {
	return fmt.Sprintf("http://%v:%v", host, port)
}

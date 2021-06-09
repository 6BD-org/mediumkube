package profiling

import (
	"net/http"
	_ "net/http/pprof"
)

type MediumkubeProfile struct{}

func Profile(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte{})
}

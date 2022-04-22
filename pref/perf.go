package perf

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/tkeel-io/kit/log"
)

// 使用方法
// perf.Init([]string{"127.0.0.1:7001"}, perf.PerfHandles()...)

type HandleFunc struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

// Init pprof
// endpoints 127.0.0.1:7001,10.0.0.1:7001
// endpoints 0.0.0.0:7001
func Init(endpoints []string, handles ...*HandleFunc) {
	pprofServeMux := http.NewServeMux()
	for _, handle := range handles {
		prefix := "/debug/" + handle.Pattern
		prefix = strings.ReplaceAll(prefix, "//", "/")
		pprofServeMux.HandleFunc(prefix, handle.Handler)
		log.Info("handle", "prefix", prefix)
	}
	if len(endpoints) == 0 {
		endpoints = append(endpoints, ":0")
	}

	for _, addr := range endpoints {
		go func() {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
			log.Info(fmt.Sprintf("debug listen on %s", listener.Addr()))
			err = http.Serve(listener, pprofServeMux)
			log.Error(fmt.Sprintf("http.ListenAndServe(%s, pprofServeMux) error(%v)", addr, err))
			panic(err)
		}()
	}
}

func PerfHandles() []*HandleFunc {
	return []*HandleFunc{
		{"/pprof/", pprof.Index},
		{"/pprof/cmdline", pprof.Cmdline},
		{"/pprof/profile", pprof.Profile},
		{"/pprof/symbol", pprof.Symbol},
	}
}

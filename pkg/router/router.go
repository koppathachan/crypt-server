package router

import (
	"net/http"
	"strings"
)

type Router interface {
	http.Handler
	SetHandlerFunc(method, path string, fn http.HandlerFunc)
}

type router struct {
	mux      map[string]map[string]http.Handler
	notFound http.Handler
}

func (rt *router) SetHandlerFunc(method, path string, fn http.HandlerFunc) {
	rt.mux[strings.ToLower(method)][path] = http.Handler(fn)
}

//ServeHTTP function to pass to myHandler
func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToLower(r.Method)
	if hm, ok := rt.mux[method]; ok {
		if h, ok := hm[r.URL.String()]; ok {
			h.ServeHTTP(w, r)
			return
		}
	}
	rt.notFound.ServeHTTP(w, r)
}

func NewRouter() Router {
	var mux map[string]map[string]http.Handler
	mux = make(map[string]map[string]http.Handler)
	mux["get"] = make(map[string]http.Handler)
	mux["post"] = make(map[string]http.Handler)
	return &router{
		mux:      mux,
		notFound: http.NotFoundHandler(),
	}
}

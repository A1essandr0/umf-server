package dserver

import (
	"net/http"
	"regexp"
)

type RoutePattern struct {
    Method  string
    Regex   *regexp.Regexp
    Handler http.HandlerFunc
}

func NewRoute(method string, pattern string, handler http.HandlerFunc) RoutePattern {
	return RoutePattern{method, regexp.MustCompile("^" + pattern + "$"), handler}
}
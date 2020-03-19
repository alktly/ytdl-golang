package main

import "net/http"

// Route describe a HTTP route
type route struct {
	method  string
	path    string
	handler func(http.ResponseWriter, *http.Request)
}

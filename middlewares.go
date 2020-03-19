package main

import "net/http"

func methodHandler(
	method string,
	f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(nil)
			return
		}

		f(w, r)
	}
}

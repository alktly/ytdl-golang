package main

import (
	"fmt"
	"net/http"
	"strings"
)

// TODO: Improve logging by including exec.Command output in response with debug flag
func badRequest(err error, w http.ResponseWriter) {
	fmt.Printf("Response: %v\n", err)
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("400 - Bad request!\n%v", err)))
}

func notFound(err error, w http.ResponseWriter) {
	fmt.Printf("Response: %v\n", err)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("404 - Not found!\n%v", err)))
}

func serverError(err error, w http.ResponseWriter) {
	fmt.Printf("Response: %v\n", err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("500 - Something bad happened\n%v", err)))
}

func success(msg string, w http.ResponseWriter) {
	fmt.Printf("Response: %v\n", msg)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}

func jsonResponse(data []byte, w http.ResponseWriter) {
	fmt.Printf("Response: %v\n", string(data))
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func writeErrorStatus(status int, err error, w http.ResponseWriter) {

	var msg string

	switch status {
	case http.StatusNotFound:
		msg = fmt.Sprintf("404 - Not found!\n%v", err)
	case http.StatusInternalServerError:
		msg = fmt.Sprintf("500 - Something bad happened\n%v", err)
	}

	fmt.Printf("Response: %v\n", msg)
	w.WriteHeader(status)
	w.Write([]byte(msg))
}

func catch404(out []byte) bool {
	if strings.Contains(string(out), "ERROR: Incomplete YouTube ID") ||
		strings.Contains(string(out), "ERROR: This video has been removed") ||
		strings.Contains(string(out), "ERROR: This video is unavailable.") {
		return true
	}

	return false
}

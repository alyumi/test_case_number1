package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	Queue   string
	Item    string
	Timeout time.Duration
}

func NewResponse(r *http.Request) *Response {
	var (
		queue   = r.URL.Path[1:]
		item    = r.URL.Query().Get("v")
		timeout = r.URL.Query().Get("timeout")
	)

	if queue[:len(queue)-1] == "/" {
		queue = queue[:len(queue)-2]
	}

	t, err := strconv.Atoi(timeout)
	if err != nil {
		log.Println("No timeout")
	}

	return &Response{
		Queue:   queue,
		Item:    item,
		Timeout: time.Duration(t) * time.Second,
	}
}

func main() {

	go http.HandleFunc("/", method)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func method(w http.ResponseWriter, r *http.Request) {
	if (r.URL.Path != "/") && (r.URL.Path != "") {
		method := r.Method
		d := NewResponse(r)
		switch method {
		case "GET":
			get(*d, w)
		case "PUT":
			put(*d, w)
		default:
			unknownMethod(w)
		}
	}
}

func unknownMethod(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func put(d Response, w http.ResponseWriter) {
	if d.Item == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}

func get(d Response, w http.ResponseWriter) {
	panic("unimplemented")
}

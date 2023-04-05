package main

import (
	"errors"
	"log"
	"net/http"
)

type Data struct {
	Queue   string
	Item    string
	Timeout int
}

func NewData() *Data {
	return &Data{
		Queue:   "",
		Item:    "",
		Timeout: 0,
	}
}

func (d Data) HasTimeout() bool {
	return d.Timeout == 0
}

func (d *Data) AddQueueName(q string) (string, error) {
	if q != "" {
		d.Queue = q
		return q, nil
	}
	return "", errors.New("no queue name")
}

func (d *Data) AddItemName(i string) error {
	if i != "" {
		d.Item = i
		return nil
	}
	return errors.New("no item name")
}

func (d *Data) AddTimeout(t int) error {
	if d.HasTimeout() {
		d.Timeout = t
		return nil
	}
	return errors.New("request does not have timeout")
}

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func copyFile(f *os.File, start int) {
	var (
		buf = make([]byte, 1024)
		fl  = false

		n       int
		err     error
		newText string
	)

	for {
		n, err = f.Read(buf)

		if err == io.EOF {
			break
		}

		if n > 0 {
			text := string(buf[:n])
			if fl == false {
				newText = text[start+1 : n]
				fl = true
			} else {
				newText = newText + text
			}
		}

		f.WriteString(newText)
	}
}

func getQueue(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[1:] + ".txt"
	absDir, err := filepath.Abs("./main.go")
	if err != nil {
		fmt.Println("hhuu")
	}
	baseDir := filepath.Dir(absDir)
	path := filepath.Join(baseDir, "\\queues\\", file)

	f, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		log.Panic("Can't open file")
	}
	defer f.Close()

	var (
		buf = make([]byte, 1024)
		n   int
		fl  = false
	)

	for {
		n, err = f.Read(buf)
		if err == io.EOF {
			break
		}

		if n > 0 {
			text := string(buf[:n])
			for i, v := range text {
				if string(v) == "\n" {
					fmt.Println(text[:i])
					w.Write([]byte(text[:i]))
					fl = true
					copyFile(f, i)
					break
				}
			}
		}

		if fl {
			break
		}
	}

}

func putQueue(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path[1:] + ".txt"
	absDir, err := filepath.Abs("./main.go")
	if err != nil {
		fmt.Println("hhuu")
	}
	baseDir := filepath.Dir(absDir)
	path := filepath.Join(baseDir, "\\queues\\", file)

	if r.URL.Query().Has("v") {
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Panic("Can't open file")
		}
		defer f.Close()
		f.WriteString(r.URL.Query().Get("v") + "\n")
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func getQueuesNames() ([]string, error) {
	files, err := os.ReadDir("queues")
	if err != nil {
		log.Println(err)
	}

	var queues []string

	for _, file := range files {
		name := file.Name()
		queues = append(queues, name[:len(name)-4])
	}

	if len(queues) == 0 {
		return nil, errors.New("No queues")
	}

	return queues, nil
}

func containsQueueFile(path string, queues []string) {

	for _, item := range queues {
		if path[1:] == item {
			return
		}
	}

	absDir, err := filepath.Abs("./main.go")
	if err != nil {
		fmt.Println("hhuu")
	}

	file, err := os.Create(filepath.Join(filepath.Dir(absDir), "\\queues\\", path+".txt"))
	if err != nil {
		log.Println(err)
	}

	defer file.Close()
}

func queueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		putQueue(w, r)
	} else if r.Method == "GET" {
		getQueue(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		queueHandler(w, r)
	}

}

func main() {
	http.HandleFunc("/", homeHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))

}

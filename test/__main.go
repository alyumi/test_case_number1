package test

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
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

func (d *Data) AddQueueName(q string) string {
	if q != "" {
		if q[len(q)-1] == '/' {
			d.Queue = q[1 : len(q)-1]
		} else {
			d.Queue = q[1:]
		}
		return d.Queue
	}
	return ""
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

func Must(err error) {
	_, filename, line, _ := runtime.Caller(1)

	if err != nil {
		log.Fatalf("[ERROR] %s:%d %v", filename, line, err)
	}
}

func MustQueue(pathName string) {

	_, filename, line, _ := runtime.Caller(1)
	absPath, e := filepath.Abs("./main.go")
	if e != nil {
		log.Printf("[ERROR] %s:%d %v", filename, line, e)
	}

	var (
		mainDir   = filepath.Dir(absPath)
		queuesDir = filepath.Join(mainDir + "\\queues")
	)

	_, filename, line, _ = runtime.Caller(1)
	f, e := os.Create(queuesDir + "\\" + pathName + ".txt")
	if e != nil {
		log.Printf("[ERROR] %s:%d %v", filename, line, e)
	}
	defer f.Close()
}

func openFile(path string, flags int, permitions fs.FileMode) *os.File {
	f, err := os.OpenFile(path, flags, permitions)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("file: ", filepath.Base(path), "openned")
	}

	return f
}

func getFirstItemInQueue(fileScanner *bufio.Scanner) (string, []string) {

	var (
		response string
		text     []string

		i = 0
	)

	for fileScanner.Scan() {
		if i == 0 {
			response = fileScanner.Text()
			i++
		} else {
			text = append(text, fileScanner.Text())
		}
	}

	return response, text
}

func copyFilesAndRemove(path, pathM, pathTemp string) {

	//Renaming part
	err := os.Rename(path, pathM)
	if err == nil {
		log.Println("Cannot rename file")
	}

	err = os.Rename(pathTemp, path)
	if err == nil {
		log.Println("Cannot rename file")
	}

	//Removing part
	err = os.Remove(pathM)
	if err == nil {
		log.Println("Cannot remove file", pathM)
	} else {
		log.Println("Removed:", pathM)
	}
}

func getFromQueue(d Data, ctx context.Context, w http.ResponseWriter) {

	var (
		path     = filepath.Dir("./main.go") + "\\queues\\" + d.Queue + ".txt"
		pathTemp = filepath.Dir("./main.go") + "\\queues\\temp.txt"
	)

	f := openFile(path, os.O_RDWR|os.O_APPEND, 0600)

	temp := openFile(pathTemp, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	defer func() {
		os.Remove(pathTemp)
	}()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var (
		response   string
		text       []string
		tempWriter = bufio.NewWriter(temp)
	)

	select {
	case <-ctx.Done():
		w.WriteHeader(http.StatusNotFound)
	default:
	}

	response, text = getFirstItemInQueue(fileScanner)

	if response == "" {
		w.WriteHeader(http.StatusNotFound)
	} else {
		ans, err := json.Marshal(map[string]string{"ans": response})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write(ans)
	}

	for _, v := range text {
		_, err := tempWriter.WriteString(v + "\n")
		if err != nil {
			log.Println("cannot write in file")
		}
	}

	if err := tempWriter.Flush(); err != nil {
		log.Println(err)
	}
	//Renaming part

	f.Close()
	temp.Close()
	err := os.Rename(pathTemp, path)
	if err != nil {
		log.Println("Cannot rename file", err)
	}

}

func putToQueue(d Data, w http.ResponseWriter) {

	path := filepath.Dir("./main.go") + "\\queues\\" + d.Queue + ".txt"
	log.Println("apth:", path)

	f := openFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0600)
	defer f.Close()

	item := d.Item
	n, err := f.WriteString(item + "\n")

	log.Println(n, " bytes written in file")

	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "" {
		d := NewData()

		d.AddQueueName(r.URL.Path)

		switch r.Method {
		case "GET":
			var (
				timeout = r.URL.Query().Get("timeout")

				t   int
				err error
			)
			if timeout != "" {
				t, err = strconv.Atoi(timeout)
				if err != nil {
					log.Println("Cannot atoi")
				}
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
			defer cancel()
			getFromQueue(*d, ctx, w)
		case "PUT":
			Must(d.AddItemName(r.URL.Query().Get("v")))
			putToQueue(*d, w)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		log.Printf("\npath: %s\nitem: %s", d.Queue, d.Item)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

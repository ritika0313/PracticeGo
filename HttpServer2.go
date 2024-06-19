package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"log"
	"os"

)

func getBase(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("received root request\n")
	io.WriteString(res, "This is home page\n")
}

/* http://127.0.0.1:3333/hello?name=ritika&lastname=chopra */
func getHi(res http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	fmt.Println("Entered hi")
	defer fmt.Println("Exiting Hi")

	select {
	//gets activated after x seconds
	case <-time.After(10 * time.Second):
		fmt.Print("Job Compeleted")
		// gets activated upon ctrl C from client
	case <-ctx.Done():
		fmt.Println("Job terminated", ctx.Err())
	}

	fmt.Printf("Received Hi Request\n")

	name := "guest"
	lastname := ""

	keys, ok := req.URL.Query()["name"]
	if ok {
		name = keys[0]
	}
	keys, ok = req.URL.Query()["lastname"]
	if ok {
		lastname = keys[0]
	}
	fmt.Fprintf(res, "Hi, Welcome! %s %s\n", name, lastname)	
	res.Write([]byte("you are here\n"))
	fmt.Fprintf(res, "enjoy here %s", req.URL.Path[1:])
}

type P struct {
	c rune
}

func (ph *P) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ph.c++

	switch r.Method {
	case "GET":
		//fmt.Println("processing GET")
		http.ServeFile(w, r, "form1.html")

	case "POST":
		fmt.Println("processing POST")
		body, err := io.ReadAll(r.Body)

		if err != nil {
			fmt.Println("error")
			return
		}
		defer r.Body.Close()
		w.Write(body)

	default:
		fmt.Fprintf(w, "Sorry only GET is supported")
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		fmt.Printf("Processing %v request\n", r.Method)
		http.ServeFile(w, r, "form.html")

	case "POST":
		fmt.Printf("Processing %v request\n", r.Method)

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Parseform error = %v", err)
			return
		}
		name := r.FormValue("name")
		occupation := r.FormValue("occupation")

		fmt.Fprintf(w, "Name = %s Occupation = %s\n", name, occupation)

	default:
		fmt.Fprintf(w, "Sorry only GET and POST is supported. Method = %v", r.Method)
	}
}

func imageHandler(w http.ResponseWriter, _ *http.Request) {
	buf, err := os.ReadFile("1.jpg")

	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "image/jpg")
	w.Write(buf)
}

func handleShutdown(wg *sync.WaitGroup) {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer wg.Done()
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	select {
	case <-shutdownCh:
		fmt.Println("Exiting gracefully")
		os.Exit(0)
	case <-time.After(2000 * time.Second):
		fmt.Println("Timeout, terminating")
		os.Exit(0)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go handleShutdown(&wg)

	ph := &P{c: 'A'}
	mux := http.NewServeMux()

	http.Handle("/file", fileserver)*/
	mux.HandleFunc("/", getBase)
	mux.HandleFunc("/hello", getHi)
	mux.HandleFunc("/form", formHandler)
	mux.Handle("/print", ph)

	mux.HandleFunc("/image", imageHandler)
	
	err := http.ListenAndServe("127.0.0.1:3333", mux)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	wg.Wait()
}

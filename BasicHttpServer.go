package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type userInfo struct {
	occupation string
}

var users = make(map[string]userInfo)

func handleRootFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling Root")
	http.ServeFile(w, r, "index.html")
}

func handleCreateFunc(w http.ResponseWriter, r *http.Request) {

	fmt.Println("received Create request")

	switch r.Method {
	case "GET":
		fmt.Println("received CREATE request")
		http.ServeFile(w, r, "form1.html")

	case "POST":
		fmt.Println("received POST request")
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Error in parsing form")
			return
		}
		Name := r.FormValue("name")
		occ := r.FormValue("occupation")
		fmt.Fprintf(w, "Name = %s Occupation = %s\n", Name, occ)

	default:
		fmt.Println("Invalid request method")
	}
}

func handleGetFunc(w http.ResponseWriter, r *http.Request) {
	var name string
	name = r.URL.Query()["name"][0]
	fmt.Println("Received GET req for USER: ", name, r.URL.Query()["name"])
	
	retUser, ok := users[name]
	if !ok {
		fmt.Println("user not found for name", name)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "%s %s\n", name, retUser.occupation)
	fmt.Printf("found %s %s\n", name, retUser.occupation)
}

var srv *http.Server

func startServer() (cancelFunc context.CancelFunc) {
	_, cancelFunc = context.WithCancel(context.Background())

	// prepare to accept client connections
	fmt.Println("Starting Server Goroutine")
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println("server could not be started OR Shutting Down")
		}
	}()
	return
}

func waitForInterrupt() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT)
	interrupt := <-sigC
	fmt.Println("Received shutdown request interrupt [%v]", interrupt)
}

func handleGracefulShutdown(cancelFunc context.CancelFunc) {
	fmt.Println("SHUTDOWN:Cancelling the context of all the server goroutines")
	cancelFunc()

	ctx, cancelTimeout := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancelTimeout()

	shutDoneChan := make(chan struct{})
	go func() {
		defer close(shutDoneChan)
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Println("Error in shutting down server")
		} else {
			fmt.Println("Server Shutdown Successfully")
		}
	}()

	select {
	case <-shutDoneChan:
		fmt.Println("Successfully completed Graceful Shutdown")
	case <-ctx.Done():
		fmt.Println("Graceful Shutdown Timeout.. Terminating")
	}
}
func main() {

	// Create a mux handler and define its routers
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRootFunc)
	mux.HandleFunc("/createUser", handleCreateFunc)
	mux.HandleFunc("/getUser", handleGetFunc)
	
	//create a server
	srv = &http.Server{
		Addr:    "127.0.0.1:3333",
		Handler: mux,
	}

	// Start a server to run in a separate go routine
	cancelFunc := startServer()

	// Once and interrupt is received, Do graceful shutdown before closing out
	defer handleGracefulShutdown(cancelFunc)

	// keep main waiting on a signal interrupt, with server service running in parallel
	waitForInterrupt()
}

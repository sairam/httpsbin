package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// InitRouter initialises the routes to handle
func InitRouter() {

	r := mux.NewRouter()

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	r.HandleFunc("/bin", createBinHandler).Methods("POST")
	r.HandleFunc("/new", createBinHandler).Methods("GET") // unexposed feature
	r.HandleFunc("/echo", echoHandler)

	r.HandleFunc("/{bin}", binViewHandler).Methods("GET").Queries("inspect", "")
	r.HandleFunc("/{bin}", binPersistHandler)

	r.HandleFunc("/", homeHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         Config.LocalServer,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("echo")

	r.Write(w)
}

func createBinHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createBin")

	dir, err := createNewDir(RandString(7))
	if err != nil {
		fmt.Println(err)
		dir, err = createNewDir(RandString(7))
		if err != nil {
			http.NotFound(w, r)
			return
			// well, we tried
		}
	}
	http.Redirect(w, r, dir, 302)
	// create a name, make a directory and redirect to the link
}
func binPersistHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("binPersist")

	bin := mux.Vars(r)["bin"]
	// ignore if there is a . in the request
	if len(bin) <= 5 || len(strings.Split(bin, ".")) > 1 {
		return
	}

	ok, err := ifDirExists(bin)
	if !ok || err != nil {
		http.NotFound(w, r)
	}

	ir := newRequest(r)
	fileName := fmt.Sprintf("%s/%d", bin, time.Now().Unix())
	ir.Save(fileName)
}

func binViewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("binView")
	bin := mux.Vars(r)["bin"]
	ok, err := ifDirExists(bin)
	if !ok || err != nil {
		http.NotFound(w, r)
	}
	irds := RetrieveLatestFromBin(bin, 10)
	fmt.Printf("%+v", irds)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home Page")
}

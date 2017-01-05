package httpsbin

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

const inspectQueryString = "inspect"

// InitRouter initialises the routes to handle
func InitRouter() {

	r := mux.NewRouter()

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "favicon.ico")
	})

	r.HandleFunc("/bin", createBinHandler).Methods("POST")
	r.HandleFunc("/new", createBinHandler).Methods("GET") // unexposed feature
	r.HandleFunc("/echo", echoHandler)

	r.HandleFunc("/{bin}", binViewHandler).Methods("GET").Queries(inspectQueryString, "")
	r.HandleFunc("/{bin}", binPersistHandler)

	r.HandleFunc("/", homeHandler).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         Config.LocalServer,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	log.Println("server is running on", Config.LocalServer)

	log.Fatal(srv.ListenAndServe())
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("op: echo")

	r.Write(w)
}

func createBinHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("op: createBin")

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
	http.Redirect(w, r, dir+"?"+inspectQueryString, 302)
}
func binPersistHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("op: binPersist")

	bin := mux.Vars(r)["bin"]
	// ignore if there is a . in the request
	if len(bin) <= 5 || len(strings.Split(bin, ".")) > 1 {
		return
	}

	ok, err := ifDirExists(bin)
	if !ok || err != nil {
		http.NotFound(w, r)
		return
	}

	ir := newRequest(r)
	fi := fmt.Sprintf("%d", time.Now().Unix())
	fileName := MergeOSPath(bin, fi)
	ir.Save(fileName)
	w.Write([]byte("ok"))
	go CleanUpMaxItemsInDir(bin)
}

func binViewHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("op: binView")
	bin := mux.Vars(r)["bin"]
	ok, err := ifDirExists(bin)
	if !ok || err != nil {
		http.NotFound(w, r)
		return
	}
	irds := RetrieveLatestFromBin(bin, 10)
	pageURL := Config.ServerProto + "://" + Config.ServerHost + "/" + bin
	DisplayPage(w, "bin", &struct {
		ThisPageURL string
		BinData     []IncomingRequestDisplay
	}{pageURL, irds})
}

// BinItem displayed on home page
type BinItem struct {
	Path  string
	Count int
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("op: Home Page")
	var rv []BinItem
	rv = append(rv, BinItem{"gQfOuZf", 2})
	rv = append(rv, BinItem{"aaaaaaa", 0})
	DisplayPage(w, "home", &struct{ RecentlyViewed []BinItem }{rv})
}

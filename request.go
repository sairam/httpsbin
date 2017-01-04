package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// IncomingRequest is the incoming request on the wire
type IncomingRequest struct {
	FromIP     string `yaml:"fromip"`
	URL        string `yaml:"url"`
	Method     string `yaml:"method"`
	ReceivedAt int64  `yaml:"received_at"`
	Request    []byte `yaml:"request,flow"`
	Body       []byte `yaml:"body,flow"`
}

func newRequest(r *http.Request) *IncomingRequest {
	var err error
	ir := &IncomingRequest{
		FromIP:     r.RemoteAddr,
		URL:        r.RequestURI,
		Method:     r.Method,
		ReceivedAt: time.Now().Unix(),
	}

	ir.Request, err = httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}

	return ir
}

func readRequest(ir *IncomingRequest) {
	a, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(ir.Request)))
	// body.(io.ReadCloser)
	if ir.Body != nil {
		a.Body = ioutil.NopCloser(bytes.NewReader(ir.Body))
	}
	fmt.Printf("%+v\n", a)
}

// Load data from the file
// Usage: (&IncomingRequest{}).Load(filename)
func (ir *IncomingRequest) Load(fileName string) {
	file, err := fsutil.OpenFile(fileName, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	out, err := readCompressedFileIO(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = yaml.Unmarshal(out, &ir); err != nil {
		fmt.Println("Error loading data ", err)
	}
}

// Save data from structure into file
func (ir *IncomingRequest) Save(fileName string) {
	out, err := yaml.Marshal(ir)
	if err != nil {
		fmt.Println("Error saving diff ", err)
		return
	}

	file, err := fsutil.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
	}

	writeCompressedFileIO(file, out)
}

type IncomingRequestDisplay struct {
	Display   string
	Reference int64
	*IncomingRequest
}

// ByInt ..
type ByInt []*IncomingRequestDisplay

func (a ByInt) Len() int           { return len(a) }
func (a ByInt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInt) Less(i, j int) bool { return a[i].Reference < a[j].Reference }

// RetrieveLatestFromBin ..
func RetrieveLatestFromBin(bin string, count int) []*IncomingRequestDisplay {
	fis, err := fsutil.ReadDir(bin)
	if err != nil {
		log.Print(err)
		return []*IncomingRequestDisplay{}
	}
	files := make([]*IncomingRequestDisplay, 0, len(fis))
	for _, fi := range fis {
		fileName := fi.Name()
		intFilename, _ := strconv.ParseInt(fileName, 10, 64)
		reference := time.Unix(intFilename, 0).Format(time.RFC1123)
		ir := &IncomingRequest{}
		ir.Load(strings.Join([]string{bin, fi.Name()}, string(os.PathSeparator)))
		files = append(files, &IncomingRequestDisplay{reference, intFilename, ir})
	}
	sort.Sort(sort.Reverse(ByInt(files)))

	if len(fis) > Config.MaxFilesToDisplay {
		go CleanUpMaxItemsInDir(bin)
	}

	return files
}

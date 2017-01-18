package httpsbin

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

	"github.com/blevesearch/bleve"

	yaml "gopkg.in/yaml.v2"
)

// IncomingRequestView ..
type IncomingRequestView struct {
	FromIP     string
	URL        string
	Method     string
	ReceivedAt int64
	Headers    []string
	Body       string
	Length     int64
	// Scheme     string
}

// IncomingRequest is the incoming request on the wire
type IncomingRequest struct {
	FromIP     string `yaml:"fromip"`
	URL        string `yaml:"url"`
	Method     string `yaml:"method"`
	ReceivedAt int64  `yaml:"received_at"`
	Request    []byte `yaml:"request,flow"`
	RequestStr string `yaml:"requeststr,flow"`
}

// IncomingRequestDisplay ..
type IncomingRequestDisplay struct {
	Display   string
	Reference int64
	IncomingRequestView
}

// ByInt ..
type ByInt []IncomingRequestDisplay

func (a ByInt) Len() int           { return len(a) }
func (a ByInt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByInt) Less(i, j int) bool { return a[i].Reference < a[j].Reference }

func newRequest(r *http.Request) *IncomingRequest {
	var err error
	ir := &IncomingRequest{
		FromIP:     r.RemoteAddr,
		URL:        r.RequestURI,
		Method:     r.Method,
		ReceivedAt: time.Now().Unix(),
	}

	ir.Request, err = httputil.DumpRequest(r, true)
	ir.RequestStr = string(ir.Request)
	if err != nil {
		fmt.Println(err)
	}

	return ir
}

func readRequest(ir *IncomingRequest) {
	a, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(ir.Request)))
	t, _ := ioutil.ReadAll(a.Body)
	fmt.Println(a)
	fmt.Println(t)
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

// Config.DataDir should be considered only when persistence is set to disk. see persist.go
func bleveIndex(dir string) (bleve.Index, error) {
	var index bleve.Index
	var err error

	if Config.Persistence == persistenceFs {
		path := MergeOSPath(Config.DataDir, dir, ".index.db")
		if index, err = bleve.Open(path); err != nil {
			mapping := bleve.NewIndexMapping()
			index, err = bleve.New(path, mapping)
		}
	} else {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.NewMemOnly(mapping)
		// maintain a hash in memory for tests
	}
	return index, err
}

// searchFromBin within the dir
func searchFromBin(dir string, queryStr string) (*bleve.SearchResult, error) {
	index, _ := bleveIndex(dir)
	query := bleve.NewQueryStringQuery(queryStr)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := index.Search(searchRequest)
	return searchResult, err
}

// Index indexes data ithe directory name
func (ir *IncomingRequest) Index(dir string, fileName string) {
	index, _ := bleveIndex(dir)
	irv := convertToView(ir)
	index.Index(fileName, irv)
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

func convertToView(ir *IncomingRequest) IncomingRequestView {
	irv := IncomingRequestView{
		FromIP:     ir.FromIP,
		URL:        ir.URL,
		Method:     ir.Method,
		ReceivedAt: ir.ReceivedAt,
	}
	var headers []string

	request, _ := http.ReadRequest(bufio.NewReader(bytes.NewReader(ir.Request)))
	for k, v := range request.Header {
		t := fmt.Sprintf("%s: %s", k, strings.Join(v, "; "))
		headers = append(headers, t)
	}
	irv.Headers = headers

	if request.Body != nil {
		body, _ := ioutil.ReadAll(request.Body)
		irv.Body = string(body)
	}
	irv.Length = request.ContentLength

	return irv
}

// RetrieveLatestFromBin ..
// TODO - split this up into two methods.
// 1. taking the files to display
// 2. loading the data from the files
func RetrieveLatestFromBin(bin string, count int) []IncomingRequestDisplay {
	fis, err := fsutil.ReadDir(bin)
	if err != nil {
		log.Print(err)
		return []IncomingRequestDisplay{}
	}
	files := make([]IncomingRequestDisplay, 0, len(fis))
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		fileName := fi.Name()

		intFilename, _ := strconv.ParseInt(fileName, 10, 64)
		reference := time.Unix(intFilename, 0).Format(time.RFC1123)
		ir := &IncomingRequest{}
		ir.Load(MergeOSPath(bin, fileName))
		ird := convertToView(ir)
		files = append(files, IncomingRequestDisplay{reference, intFilename, ird})
	}
	sort.Sort(sort.Reverse(ByInt(files)))

	if len(fis) > Config.MaxFilesToDisplay {
		go CleanUpMaxItemsInDir(bin)
	}

	return files
}

// SearchInBin ..
func SearchInBin(bin, query string) ([]IncomingRequestDisplay, int) {
	result, err := searchFromBin(bin, query)

	fis, _ := fsutil.ReadDir(bin)
	total := len(fis) - 1
	if total < 0 {
		total = 0
	}

	if err != nil || result.Hits.Len() == 0 {
		return []IncomingRequestDisplay{}, total
	}

	files := make([]IncomingRequestDisplay, 0, result.Hits.Len())
	for _, hit := range result.Hits {
		fi := hit.ID
		fileName := MergeOSPath(fi)

		// copied from RetrieveLatestFromBin
		intFilename, _ := strconv.ParseInt(fileName, 10, 64)
		reference := time.Unix(intFilename, 0).Format(time.RFC1123)
		ir := &IncomingRequest{}
		ir.Load(MergeOSPath(bin, fileName))
		ird := convertToView(ir)
		files = append(files, IncomingRequestDisplay{reference, intFilename, ird})
	}

	return files, total
}

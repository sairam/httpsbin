package httpsbin

import (
	"html/template"
	"io"

	"github.com/sairam/kinli"
)

// InitView ..
func InitView() {
	kinli.ViewFuncs = template.FuncMap{}
	kinli.CacheMode = Config.CacheView
	kinli.InitTmpl()
}

// DisplayPage ..
func DisplayPage(w io.Writer, path string, ctx interface{}) {
	page := &kinli.Page{
		Title:   "Home Page",
		Context: ctx,
		ClientConfig: map[string]string{
			"WebsiteName":     "HttpBin",
			"EmailContact":    "sairam.kunala@gmail.com",
			"GoogleAnalytics": "",
		},
	}
	kinli.DisplayPage(w, path, page)
}

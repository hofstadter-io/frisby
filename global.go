package frisby

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"bitbucket.org/_metalogic_/request"
)

type globalData struct {
	Req  *request.Request
	Errs map[string][]error

	NumRequest int
	NumAsserts int
	NumErrored int

	PrintProgressName bool
	PrintProgressDot  bool

	PathSeparator string
}

// DefaultPathSeparator is the default path separator
const DefaultPathSeparator = "."

func init() {
	Global.Req = request.NewRequest(new(http.Client))
	Global.Errs = make(map[string][]error, 0)
	Global.PrintProgressDot = true
	Global.PathSeparator = DefaultPathSeparator
}

// Set BasicAuth values for the coming request
func (G *globalData) BasicAuth(user, passwd string) *globalData {
	G.Req.BasicAuth = request.BasicAuth{Username: user, Password: passwd}
	return G
}

// Set Proxy URL for the coming request
func (G *globalData) SetProxy(url string) *globalData {
	G.Req.Proxy = url
	return G
}

// Set a Header value for the coming request
func (G *globalData) SetHeader(key, value string) *globalData {
	if G.Req.Headers == nil {
		G.Req.Headers = make(map[string]string)
	}
	G.Req.Headers[key] = value
	return G
}

// Set several Headers for the coming request
func (G *globalData) SetHeaders(headers map[string]string) *globalData {
	if G.Req.Headers == nil {
		G.Req.Headers = make(map[string]string)
	}
	for key, value := range headers {
		G.Req.Headers[key] = value
	}
	return G
}

// Set a Cookie value for the coming request
func (G *globalData) SetCookie(key, value string) *globalData {
	if G.Req.Cookies == nil {
		G.Req.Cookies = make(map[string]string)
	}
	G.Req.Cookies[key] = value
	return G
}

// Set several Cookie values for the coming request
func (G *globalData) SetCookies(cookies map[string]string) *globalData {
	if G.Req.Cookies == nil {
		G.Req.Cookies = make(map[string]string)
	}
	for key, value := range cookies {
		G.Req.Cookies[key] = value
	}
	return G
}

// Set a Gorm data for the coming request
func (G *globalData) SetData(key, value string) *globalData {
	if G.Req.Data == nil {
		G.Req.Data = make(map[string]string)
	}
	G.Req.Data[key] = value
	return G
}

// Set several Gorm data for the coming request
func (G *globalData) SetDatas(datas map[string]string) *globalData {
	if G.Req.Data == nil {
		G.Req.Data = make(map[string]string)
	}
	for key, value := range datas {
		G.Req.Data[key] = value
	}
	return G
}

// Set a url Param for the coming request
func (G *globalData) SetParam(key, value string) *globalData {
	if G.Req.Params == nil {
		G.Req.Params = make(map[string]string)
	}
	G.Req.Params[key] = value
	return G
}

// Set several url Param for the coming request
func (G *globalData) SetParams(params map[string]string) *globalData {
	if G.Req.Params == nil {
		G.Req.Params = make(map[string]string)
	}
	for key, value := range params {
		G.Req.Params[key] = value
	}
	return G
}

// Set the JSON body for the coming request
func (G *globalData) SetJSON(json interface{}) *globalData {
	G.Req.JSON = json
	return G
}

// Add a file to the Gorm data for the coming request
func (G *globalData) AddFile(filename string) *globalData {
	file, err := os.Open(filename)
	if err != nil {
		G.AddError("Global", err.Error())
		fmt.Println("Error adding file to global")
	} else {
		fileField := request.FileField{FieldName: "file", FileName: filename, File: file}
		G.Req.Files = append(G.Req.Files, fileField)
	}
	return G
}

// Manually add an error, if you need to
func (G *globalData) AddError(name, errStr string) *globalData {
	G.NumErrored++
	err := errors.New(errStr)
	G.Errs[name] = append(G.Errs[name], err)
	return G
}

// Get all errors for the globalData object
//
// This function should be called last
func (G *globalData) Errors() map[string][]error {
	return G.Errs
}

// Prints a report for the FrisbyGlobal Object
//
// If there are any errors, they will all be printed as well
func (G *globalData) PrintReport() *globalData {
	fmt.Printf("\nFor %d requests made\n", G.NumRequest)
	if len(G.Errs) == 0 {
		fmt.Printf("  All tests passed\n")
	} else {
		fmt.Printf("  FAILED  [%d/%d]\n", G.NumErrored, G.NumAsserts)
		for key, val := range G.Errs {
			fmt.Printf("      [%s]\n", key)
			for _, e := range val {
				fmt.Println("        - ", e)
			}
		}
	}

	return G
}

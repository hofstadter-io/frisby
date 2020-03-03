package frisby

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/_metalogic_/request"
)

// Global stores information used in multiple requests
var Global globalData

const defaultFileKey = "file"

// Frisby holds request/response
type Frisby struct {
	Name   string
	URL    string
	Method string

	Req           *request.Request
	Resp          *request.Response
	Errs          []error
	ExecutionTime float64
}

// Create creates a new Frisby object with the given name.
//
// The given name will be used if you call PrintReport()
func Create(name string) *Frisby {
	F := new(Frisby)
	F.Name = name
	F.Req = request.NewRequest(new(http.Client))
	F.Errs = make([]error, 0)

	// copy in global settings
	F.Req.BasicAuth = Global.Req.BasicAuth
	F.Req.BearerAuth = Global.Req.BearerAuth
	F.Req.Proxy = Global.Req.Proxy
	F.SetHeaders(Global.Req.Headers)
	F.SetCookies(Global.Req.Cookies)
	F.SetDatas(Global.Req.Data)
	F.SetParams(Global.Req.Params)
	F.Req.JSON = Global.Req.JSON
	F.Req.Files = append(F.Req.Files, Global.Req.Files...)

	// initialize request
	F.Req.Params = make(map[string]string)

	return F
}

// Get sets the HTTP method to GET for the given URL
func (F *Frisby) Get(url string) *Frisby {
	F.Method = "GET"
	F.URL = url
	return F
}

// Post sets the HTTP method to POST for the given URL
func (F *Frisby) Post(url string) *Frisby {
	F.Method = "POST"
	F.URL = url
	return F
}

// Put sets the HTTP method to PUT for the given URL
func (F *Frisby) Put(url string) *Frisby {
	F.Method = "PUT"
	F.URL = url
	return F
}

// Patch sets the HTTP method to PATCH for the given URL
func (F *Frisby) Patch(url string) *Frisby {
	F.Method = "PATCH"
	F.URL = url
	return F
}

// Delete sets the HTTP method to DELETE for the given URL
func (F *Frisby) Delete(url string) *Frisby {
	F.Method = "DELETE"
	F.URL = url
	return F
}

// Head sets the HTTP method to HEAD for the given URL
func (F *Frisby) Head(url string) *Frisby {
	F.Method = "HEAD"
	F.URL = url
	return F
}

// Options sets the HTTP method to OPTIONS for the given URL
func (F *Frisby) Options(url string) *Frisby {
	F.Method = "OPTIONS"
	F.URL = url
	return F
}

// BasicAuth sets BasicAuth values for the coming request
func (F *Frisby) BasicAuth(user, passwd string) *Frisby {
	F.Req.BasicAuth = request.BasicAuth{Username: user, Password: passwd}
	return F
}

// SetProxy sets proxy URL for the coming request
func (F *Frisby) SetProxy(url string) *Frisby {
	F.Req.Proxy = url
	return F
}

// SetHeader sets a Header value for the coming request
func (F *Frisby) SetHeader(key, value string) *Frisby {
	if F.Req.Headers == nil {
		F.Req.Headers = make(map[string]string)
	}
	F.Req.Headers[key] = value
	return F
}

// SetHeaders sets several Headers for the coming request
func (F *Frisby) SetHeaders(headers map[string]string) *Frisby {
	if F.Req.Headers == nil {
		F.Req.Headers = make(map[string]string)
	}
	for key, value := range headers {
		F.Req.Headers[key] = value
	}
	return F
}

// SetCookie sets a Cookie value for the coming request
func (F *Frisby) SetCookie(key, value string) *Frisby {
	if F.Req.Cookies == nil {
		F.Req.Cookies = make(map[string]string)
	}
	F.Req.Cookies[key] = value
	return F
}

// SetCookies sets several Cookie values for the coming request
func (F *Frisby) SetCookies(cookies map[string]string) *Frisby {
	if F.Req.Cookies == nil {
		F.Req.Cookies = make(map[string]string)
	}
	for key, value := range cookies {
		F.Req.Cookies[key] = value
	}
	return F
}

// SetData sets a Form data for the coming request
func (F *Frisby) SetData(key, value string) *Frisby {
	if F.Req.Data == nil {
		F.Req.Data = make(map[string]string)
	}
	F.Req.Data[key] = value
	return F
}

// SetDatas sets several Form data for the coming request
func (F *Frisby) SetDatas(datas map[string]string) *Frisby {
	if F.Req.Data == nil {
		F.Req.Data = make(map[string]string)
	}
	for key, value := range datas {
		F.Req.Data[key] = value
	}
	return F
}

// SetParam sets a url Param for the coming request
func (F *Frisby) SetParam(key, value string) *Frisby {
	if F.Req.Params == nil {
		F.Req.Params = make(map[string]string)
	}
	F.Req.Params[key] = value
	return F
}

// SetParams sets several url Param for the coming request
func (F *Frisby) SetParams(params map[string]string) *Frisby {
	if F.Req.Params == nil {
		F.Req.Params = make(map[string]string)
	}
	for key, value := range params {
		F.Req.Params[key] = value
	}
	return F
}

// SetJSON sets the JSON body for the coming request
func (F *Frisby) SetJSON(json interface{}) *Frisby {
	F.Req.JSON = json
	return F
}

// AddFile adds a file to the Form data for the coming request
func (F *Frisby) AddFile(filename string) *Frisby {
	file, err := os.Open(filename)
	if err != nil {
		F.Errs = append(F.Errs, err)
	} else {
		fileField := request.FileField{
			FieldName: defaultFileKey,
			FileName:  filepath.Base(filename),
			File:      file}
		F.Req.Files = append(F.Req.Files, fileField)
	}
	return F
}

// AddFileByKey adds a file to the Form data for the coming request
func (F *Frisby) AddFileByKey(key, filename string) *Frisby {
	file, err := os.Open(filename)
	if err != nil {
		F.Errs = append(F.Errs, err)
	} else {
		if len(key) == 0 {
			key = defaultFileKey
		}
		fileField := request.FileField{
			FieldName: key,
			FileName:  filepath.Base(filename),
			File:      file}
		F.Req.Files = append(F.Req.Files, fileField)
	}
	return F
}

// Send the actual request to the URL
func (F *Frisby) Send() *Frisby {
	Global.NumRequest++
	if Global.PrintProgressName {
		fmt.Println(F.Name)
	} else if Global.PrintProgressDot {
		fmt.Printf("")
	}

	start := time.Now()

	var err error
	switch F.Method {
	case "GET":
		F.Resp, err = F.Req.Get(F.URL)
	case "POST":
		F.Resp, err = F.Req.Post(F.URL)
	case "PUT":
		F.Resp, err = F.Req.Put(F.URL)
	case "PATCH":
		F.Resp, err = F.Req.Patch(F.URL)
	case "DELETE":
		F.Resp, err = F.Req.Delete(F.URL)
	case "HEAD":
		F.Resp, err = F.Req.Head(F.URL)
	case "OPTIONS":
		F.Resp, err = F.Req.Options(F.URL)
	}

	F.ExecutionTime = time.Since(start).Seconds()

	if err != nil {
		F.Errs = append(F.Errs, err)
	}

	return F
}

// AddError manually adds an error, if you need to
func (F *Frisby) AddError(errStr string) *Frisby {
	err := errors.New(errStr)
	F.Errs = append(F.Errs, err)
	Global.AddError(F.Name, errStr)
	return F
}

// Error gets the most recent error for the Frisby object
//
// This function should be called last
func (F *Frisby) Error() error {
	if len(F.Errs) > 0 {
		return F.Errs[len(F.Errs)-1]
	}
	return nil
}

// Errors gets all errors for the Frisby object
//
// This function should be called last
func (F *Frisby) Errors() []error {
	return F.Errs
}

// PauseTest pauses your testrun for a defined amount of seconds
func (F *Frisby) PauseTest(t time.Duration) *Frisby {
	time.Sleep(t * time.Second)
	return nil
}

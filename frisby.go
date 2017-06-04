package frisby

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mozillazg/request"
)

var Global global_data

const defaultFileKey = "file"

type Frisby struct {
	Name   string
	Url    string
	Method string

	Req           *request.Request
	Resp          *request.Response
	Errs          []error
	ExecutionTime float64
}

// Creates a new Frisby object with the given name.
//
// The given name will be used if you call PrintReport()
func Create(name string) *Frisby {
	F := new(Frisby)
	F.Name = name
	F.Req = request.NewRequest(new(http.Client))
	F.Errs = make([]error, 0)

	// copy in global settings
	F.Req.BasicAuth = Global.Req.BasicAuth
	F.Req.Proxy = Global.Req.Proxy
	F.SetHeaders(Global.Req.Headers)
	F.SetCookies(Global.Req.Cookies)
	F.SetDatas(Global.Req.Data)
	F.SetParams(Global.Req.Params)
	F.Req.Json = Global.Req.Json
	F.Req.Files = append(F.Req.Files, Global.Req.Files...)

	// initialize request
	F.Req.Params = make(map[string]string)

	return F
}

// Set the HTTP method to GET for the given URL
func (F *Frisby) Get(url string) *Frisby {
	F.Method = "GET"
	F.Url = url
	return F
}

// Set the HTTP method to POST for the given URL
func (F *Frisby) Post(url string) *Frisby {
	F.Method = "POST"
	F.Url = url
	return F
}

// Set the HTTP method to PUT for the given URL
func (F *Frisby) Put(url string) *Frisby {
	F.Method = "PUT"
	F.Url = url
	return F
}

// Set the HTTP method to PATCH for the given URL
func (F *Frisby) Patch(url string) *Frisby {
	F.Method = "PATCH"
	F.Url = url
	return F
}

// Set the HTTP method to DELETE for the given URL
func (F *Frisby) Delete(url string) *Frisby {
	F.Method = "DELETE"
	F.Url = url
	return F
}

// Set the HTTP method to HEAD for the given URL
func (F *Frisby) Head(url string) *Frisby {
	F.Method = "HEAD"
	F.Url = url
	return F
}

// Set the HTTP method to OPTIONS for the given URL
func (F *Frisby) Options(url string) *Frisby {
	F.Method = "OPTIONS"
	F.Url = url
	return F
}

// Set BasicAuth values for the coming request
func (F *Frisby) BasicAuth(user, passwd string) *Frisby {
	F.Req.BasicAuth = request.BasicAuth{user, passwd}
	return F
}

// Set Proxy URL for the coming request
func (F *Frisby) SetProxy(url string) *Frisby {
	F.Req.Proxy = url
	return F
}

// Set a Header value for the coming request
func (F *Frisby) SetHeader(key, value string) *Frisby {
	if F.Req.Headers == nil {
		F.Req.Headers = make(map[string]string)
	}
	F.Req.Headers[key] = value
	return F
}

// Set several Headers for the coming request
func (F *Frisby) SetHeaders(headers map[string]string) *Frisby {
	if F.Req.Headers == nil {
		F.Req.Headers = make(map[string]string)
	}
	for key, value := range headers {
		F.Req.Headers[key] = value
	}
	return F
}

// Set a Cookie value for the coming request
func (F *Frisby) SetCookie(key, value string) *Frisby {
	if F.Req.Cookies == nil {
		F.Req.Cookies = make(map[string]string)
	}
	F.Req.Cookies[key] = value
	return F
}

// Set several Cookie values for the coming request
func (F *Frisby) SetCookies(cookies map[string]string) *Frisby {
	if F.Req.Cookies == nil {
		F.Req.Cookies = make(map[string]string)
	}
	for key, value := range cookies {
		F.Req.Cookies[key] = value
	}
	return F
}

// Set a Form data for the coming request
func (F *Frisby) SetData(key, value string) *Frisby {
	if F.Req.Data == nil {
		F.Req.Data = make(map[string]string)
	}
	F.Req.Data[key] = value
	return F
}

// Set several Form data for the coming request
func (F *Frisby) SetDatas(datas map[string]string) *Frisby {
	if F.Req.Data == nil {
		F.Req.Data = make(map[string]string)
	}
	for key, value := range datas {
		F.Req.Data[key] = value
	}
	return F
}

// Set a url Param for the coming request
func (F *Frisby) SetParam(key, value string) *Frisby {
	if F.Req.Params == nil {
		F.Req.Params = make(map[string]string)
	}
	F.Req.Params[key] = value
	return F
}

// Set several url Param for the coming request
func (F *Frisby) SetParams(params map[string]string) *Frisby {
	if F.Req.Params == nil {
		F.Req.Params = make(map[string]string)
	}
	for key, value := range params {
		F.Req.Params[key] = value
	}
	return F
}

// Set the JSON body for the coming request
func (F *Frisby) SetJson(json interface{}) *Frisby {
	F.Req.Json = json
	return F
}

// Add a file to the Form data for the coming request
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

// Add a file to the Form data for the coming request
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
		F.Resp, err = F.Req.Get(F.Url)
	case "POST":
		F.Resp, err = F.Req.Post(F.Url)
	case "PUT":
		F.Resp, err = F.Req.Put(F.Url)
	case "PATCH":
		F.Resp, err = F.Req.Patch(F.Url)
	case "DELETE":
		F.Resp, err = F.Req.Delete(F.Url)
	case "HEAD":
		F.Resp, err = F.Req.Head(F.Url)
	case "OPTIONS":
		F.Resp, err = F.Req.Options(F.Url)
	}

	F.ExecutionTime = time.Since(start).Seconds()

	if err != nil {
		F.Errs = append(F.Errs, err)
	}

	return F
}

// Manually add an error, if you need to
func (F *Frisby) AddError(err_str string) *Frisby {
	err := errors.New(err_str)
	F.Errs = append(F.Errs, err)
	Global.AddError(F.Name, err_str)
	return F
}

// Get the most recent error for the Frisby object
//
// This function should be called last
func (F *Frisby) Error() error {
	if len(F.Errs) > 0 {
		return F.Errs[len(F.Errs)-1]
	}
	return nil
}

// Get all errors for the Frisby object
//
// This function should be called last
func (F *Frisby) Errors() []error {
	return F.Errs
}

// Pause your testrun for a defined amount of seconds
func (F *Frisby) PauseTest(t time.Duration) *Frisby {
	time.Sleep(t * time.Second)
	return nil
}

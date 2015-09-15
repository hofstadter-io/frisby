package frisby

import (
	"net/http"
	"os"

	"github.com/mozillazg/request"
)

type Frisby struct {
	Name   string
	Url    string
	Method string

	req  *request.Request
	resp *request.Response
	errs []error
}

// Creates a new Frisby object with the given name.
//
// The given name will be used if you call PrintReport()
func Create(name string) *Frisby {
	F := new(Frisby)
	F.Name = name
	F.req = request.NewRequest(new(http.Client))

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
	F.req.BasicAuth = request.BasicAuth{user, passwd}
	return F
}

// Set Proxy URL for the coming request
func (F *Frisby) Proxy(url string) *Frisby {
	F.req.Proxy = url
	return F
}

// Set a Header value for the coming request
func (F *Frisby) SetHeader(key, value string) *Frisby {
	F.req.Headers[key] = value
	return F
}

// Set several Headers for the coming request
func (F *Frisby) SetHeaders(headers map[string]string) *Frisby {
	for key, value := range headers {
		F.req.Headers[key] = value
	}
	return F
}

// Set a Cookie value for the coming request
func (F *Frisby) SetCookie(key, value string) *Frisby {
	F.req.Cookies[key] = value
	return F
}

// Set several Cookie values for the coming request
func (F *Frisby) SetCookies(cookies map[string]string) *Frisby {
	for key, value := range cookies {
		F.req.Cookies[key] = value
	}
	return F
}

// Set a Form data for the coming request
func (F *Frisby) SetData(key, value string) *Frisby {
	if F.req.Data == nil {
		F.req.Data = make(map[string]string)
	}
	F.req.Data[key] = value
	return F
}

// Set several Form data for the coming request
func (F *Frisby) SetDatas(datas map[string]string) *Frisby {
	for key, value := range datas {
		F.req.Data[key] = value
	}
	return F
}

// Set a url Param for the coming request
func (F *Frisby) SetParam(key, value string) *Frisby {
	F.req.Params[key] = value
	return F
}

// Set several url Param for the coming request
func (F *Frisby) SetParams(params map[string]string) *Frisby {
	for key, value := range params {
		F.req.Params[key] = value
	}
	return F
}

// Set the JSON body for the coming request
func (F *Frisby) SetJson(json interface{}) *Frisby {
	F.req.Json = json
	return F
}

// Add a file to the Form data for the coming request
func (F *Frisby) AddFile(filename string) *Frisby {
	file, err := os.Open("test.txt")
	if err != nil {
		F.errs = append(F.errs, err)
	} else {
		fileField := request.FileField{"file", "test.txt", file}
		F.req.Files = append(F.req.Files, fileField)
	}
	return F
}

// Send the actual request to the URL
func (F *Frisby) Send() *Frisby {
	var err error
	switch F.Method {
	case "GET":
		F.resp, err = F.req.Get(F.Url)
	case "POST":
		F.resp, err = F.req.Post(F.Url)
	case "PUT":
		F.resp, err = F.req.Put(F.Url)
	case "PATCH":
		F.resp, err = F.req.Patch(F.Url)
	case "DELETE":
		F.resp, err = F.req.Delete(F.Url)
	case "HEAD":
		F.resp, err = F.req.Head(F.Url)
	case "OPTIONS":
		F.resp, err = F.req.Options(F.Url)
	}

	if err != nil {
		F.errs = append(F.errs, err)
	}
	return F
}

// Get the most recent error for the Frisby object
//
// This function should be called last
func (F *Frisby) Error() error {
	return F.errs[len(F.errs)-1]
}

// Get all errors for the Frisby object
//
// This function should be called last
func (F *Frisby) Errors() []error {
	return F.errs
}

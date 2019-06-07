package aidi

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)
var Global = global_data{}
type Aidi struct {
	Name string
	Url string
	Req *http.Request
	Resp *http.Response
	Errs  []error
	ExecutionTime float64
	body  io.Reader
}

func CreateCase(name string) *Aidi {
	a := new(Aidi)
	a.Name = name
	a.Req = &http.Request{}
	a.Resp = &http.Response{}
	a.Errs = make([]error, 0)

	return a
}

func(a *Aidi) Get(url string) *Aidi {
	a.Req.Method = "GET"
	a.Url = url
	return a
}

func(a *Aidi) Post(url string) *Aidi {
	a.Req.Method = "POST"
	a.Url = url
	return a
}

func(a *Aidi) Put(url string) *Aidi {
	a.Req.Method = "PUT"
	a.Url = url
	return a
}

// Set the HTTP method to PATCH for the given URL
func (a *Aidi) Patch(url string) *Aidi {
	a.Req.Method = "PATCH"
	a.Url = url
	return a
}

// Set the HTTP method to DELETE for the given URL
func (a *Aidi) Delete(url string) *Aidi {
	a.Req.Method = "DELETE"
	a.Url = url
	return a
}

// Set a Header value for the coming request
func (a *Aidi) SetHeader(key, value string) *Aidi {
	a.Req.Header.Set(key, value)
	return a
}

// Set a Headers value for the coming request
func (a *Aidi) SetHeaders(headers map[string]string) *Aidi {
	for key, value := range headers {
		a.Req.Header.Set(key, value)
	}
	return a
}

// Set a Cookie value for the coming request
func (a *Aidi) SetCookie(cookie *http.Cookie) *Aidi {
	a.Req.AddCookie(cookie)
	return a
}


// Add a file to the Form data for the coming request
func (a *Aidi) AddFile(filename string) *Aidi {
	file, err := os.Open(filename)
	if err != nil {
		a.Errs = append(a.Errs, err)
		return a
	}
	defer file.Close()

	a.body = file
	return a
}

// Add a file to the Form data for the coming request
func (a *Aidi) SetBody(body io.Reader) *Aidi {
	a.body = body
	return a
}


// Send the actual request to the URL
func (a *Aidi) Send() *Aidi {
	Global.NumRequest++
	if Global.PrintProgressName {
		logrus.Println(a.Name)
	} else if Global.PrintProgressDot {
		logrus.Println("")
	}

	// set reuqest body
	rc, ok := a.body.(io.ReadCloser)
	if !ok && a.body != nil {
		rc = ioutil.NopCloser(a.body)
	}
	a.Req.Body = rc

	// set request url
	url, err := url.Parse(a.Url)
	if err != nil {
		a.Errs = append(a.Errs, err)
		return a
	}
	a.Req.URL = url
	client := &http.Client{}

	a.Resp, err = client.Do(a.Req)
	if err != nil {
		a.Errs = append(a.Errs, err)
	}

	return a
}

// Add a file to the Form data for the coming request
func (a *Aidi) AddError(message string) *Aidi {
	a.Errs = append(a.Errs, errors.New(message))
	return a
}
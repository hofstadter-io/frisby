package frisby

import (
	"net/http"
	"os"

	"github.com/mozillazg/request"
)

/*
type Args struct {
    Client    *http.Client
    Headers   map[string]string
    Cookies   map[string]string
    Data      map[string]string
    Params    map[string]string
    Files     []FileField
    Json      interface{}
    Proxy     string
    BasicAuth BasicAuth
}
*/

type Frisby struct {
	Name   string
	Url    string
	Method string

	req  *request.Request
	resp *request.Response
	errs []error
}

func Create(name string) *Frisby {
	F := new(Frisby)
	F.Name = name
	F.req = request.NewRequest(new(http.Client))

	return F
}

func (F *Frisby) Get(url string) *Frisby {
	F.Method = "GET"
	F.Url = url
	return F
}

func (F *Frisby) Post(url string) *Frisby {
	F.Method = "POST"
	F.Url = url
	return F
}

func (F *Frisby) Put(url string) *Frisby {
	F.Method = "PUT"
	F.Url = url
	return F
}

func (F *Frisby) Patch(url string) *Frisby {
	F.Method = "PATCH"
	F.Url = url
	return F
}

func (F *Frisby) Delete(url string) *Frisby {
	F.Method = "DELETE"
	F.Url = url
	return F
}

func (F *Frisby) Head(url string) *Frisby {
	F.Method = "HEAD"
	F.Url = url
	return F
}

func (F *Frisby) Options(url string) *Frisby {
	F.Method = "OPTIONS"
	F.Url = url
	return F
}

// func (F *Frisby) Send() (*Response, error) {
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

func (F *Frisby) Error() error {
	return F.errs[len(F.errs)-1]
}

func (F *Frisby) Errors() []error {
	return F.errs
}

func (F *Frisby) BasicAuth(user, passwd string) *Frisby {
	F.req.BasicAuth = request.BasicAuth{user, passwd}
	return F
}

func (F *Frisby) Proxy(url string) *Frisby {
	F.req.Proxy = url
	return F
}

func (F *Frisby) SetHeader(key, value string) *Frisby {
	F.req.Headers[key] = value
	return F
}

func (F *Frisby) SetHeaders(headers map[string]string) *Frisby {
	for key, value := range headers {
		F.req.Headers[key] = value
	}
	return F
}

func (F *Frisby) SetCookie(key, value string) *Frisby {
	F.req.Cookies[key] = value
	return F
}

func (F *Frisby) SetCookies(cookies map[string]string) *Frisby {
	for key, value := range cookies {
		F.req.Cookies[key] = value
	}
	return F
}

func (F *Frisby) SetData(key, value string) *Frisby {
	F.req.Data[key] = value
	return F
}

func (F *Frisby) SetDatas(datas map[string]string) *Frisby {
	for key, value := range datas {
		F.req.Data[key] = value
	}
	return F
}

func (F *Frisby) SetParam(key, value string) *Frisby {
	F.req.Params[key] = value
	return F
}

func (F *Frisby) SetParams(params map[string]string) *Frisby {
	for key, value := range params {
		F.req.Params[key] = value
	}
	return F
}

func (F *Frisby) SetJson(json interface{}) *Frisby {
	F.req.Json = json
	return F
}

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

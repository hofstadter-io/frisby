# frisby

[![Build Status](https://travis-ci.org/verdverm/frisby.svg?branch=master)](https://travis-ci.org/verdverm/frisby)
[![GoDoc](https://godoc.org/github.com/verdverm/frisby?status.svg)](https://godoc.org/github.com/verdverm/frisby)

REST API testing framework inspired by frisby-js, written in Go

### Installation

```
go get -u github.com/verdverm/frisby
```

### Basic Usage

First create a Frisby object:

```
// create an object with a given name (used in the report)
F := frisby.Create("Test successful user login").
    Get("https://golang.org")
```

Next perform the operation:

```
F.Send()
```

Then inspect the response:

```
F.ExpectStatus(200).
    ExpectContent("The Go Programming Language").
    PrintReport()
```

Check any error(s):

`err := F.Error()`

```
errs := F.Errors()
for _,e := range errs {
	fmt.Println("Error: ", e)
}
```

There is also a Global object for setting repeated Pre-flight options.

```
frisby.Global.BasicAuth("username", "password")
```

### HTTP Method functions

Your basic HTTP verbs:

* Get(url string)
* Post(url string)
* Put(url string)
* Patch(url string)
* Delete(url string)
* Head(url string)
* Options(url string)

### Pre-flight functions

Functions called before `Send()`

You can also set theses on the `frisby.Global` object for persisting state over multiple requests.

( Most of these come from [github.com/mozillazg/request](https://github.com/mozillazg/request))

* BasicAuth(username,password string)
* Proxy(url string)
* SetHeader(key,value string)
* SetHeaders(map[string]string)
* SetCookies(key,value string)
* SetCookiess(map[string]string)
* SetDate(key,value string)
* SetDates(map[string]string)
* SetParam(key,value string)
* SetParams(map[string]string)
* SetJson(interface{})
* SetFile(filename string)


### Post-flight functions

Functions called after `Send()`

* ExpectStatus(code int)
* ExpectHeader(key, value string)
* ExpectContent(content string)
* ExpectJson(path string, value interface{})
* ExpectJsonLength(path string, length int)
* ExpectJsonType(path string, value_type reflect.Kind)
* AfterContent( func(Frisby,[]byte,error) )
* AfterText( func(Frisby,string,error) )
* AfterJson( func(Frisby,simplejson.Json,error) )
* PrintBody()
* PrintReport()


### More examples

```
package main

import (
	"fmt"
	"reflect"

	"github.com/bitly/go-simplejson"
	"github.com/verdverm/frisby"
)

func main() {
	fmt.Println("Frisby!\n")

	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language")

	frisby.Create("Test GET Go homepage (which fails)").
		Get("http://golang.org").
		Send().
		ExpectStatus(400).
		ExpectContent("A string which won't be found")

	frisby.Create("Test POST").
		Post("http://httpbin.org/post").
		SetData("test_key", "test_value").
		Send().
		ExpectStatus(200)

	frisby.Create("Test ExpectJsonType").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJsonType("url", reflect.String)

	frisby.Create("Test ExpectJson").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJson("url", "http://httpbin.org/post").
		ExpectJson("headers.Accept", "*/*")

	frisby.Create("Test ExpectJsonLength (which fails)").
		Post("http://httpbin.org/post").
		SetJson([]string{"item1", "item2", "item3"}).
		Send().
		ExpectStatus(200).
		ExpectJsonLength("json", 4)

	frisby.Create("Test AfterJson").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
		val, _ := json.Get("url").String()
		frisby.Global.SetProxy(val)
	})

	frisby.Global.PrintReport()
}

```

Sample Output

```
Frisby!

.......
For 7 requests made
  FAILED  [3/13]
      [Test ExpectJsonLength]
        -  Expect length to be 4, but got 3
      [Test GET Go homepage (which fails)]
        -  Expected Status 400, but got 200: "200 OK"
        -  Expected Body to contain "A string which won't be found", but it was missing
```

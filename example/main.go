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

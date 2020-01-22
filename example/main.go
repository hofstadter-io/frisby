package main

import (
	"fmt"
	"reflect"

	"github.com/EducationPlannerBC/frisby"
	"github.com/bitly/go-simplejson"
)

func oldmain() {
	fmt.Println("Frisby!")

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

	frisby.Create("Test ExpectJSONType").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJSONType("url", reflect.String)

	frisby.Create("Test ExpectJSON").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJSON("url", "http://httpbin.org/post").
		ExpectJSON("headers.Accept", "*/*")

	frisby.Create("Test ExpectJSONLength").
		Post("http://httpbin.org/post").
		SetJSON([]string{"item1", "item2", "item3"}).
		Send().
		ExpectStatus(200).
		ExpectJSONLength("json", 3)

	frisby.Create("Test AfterJSON").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		AfterJSON(func(F *frisby.Frisby, json *simplejson.Json, err error) {
			val, _ := json.Get("url").String()
			frisby.Global.SetProxy(val)
		})

	frisby.Global.PrintReport()
}

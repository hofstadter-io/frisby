package frisby_test

import (
	"fmt"
	"reflect"

	"github.com/EducationPlannerBC/frisby"
	"github.com/bitly/go-simplejson"
)

func init() {
	frisby.Global.PrintProgressDot = false
}

func ExampleFrisby_Get() {
	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language").
		PrintReport()

	// Output: Pass  [Test GET Go homepage]
}

func ExampleFrisby_Post() {
	frisby.Create("Test POST").
		Post("http://httpbin.org/post").
		SetData("test_key", "test_value").
		Send().
		ExpectStatus(200).
		PrintReport()

	// Output: Pass  [Test POST]
}

func ExampleFrisby_PrintReport() {
	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(400).
		ExpectContent("A string which won't be found").
		AddError("Manually Added Error").
		PrintReport()

	// Output: FAIL  [Test GET Go homepage]
	//         -  Expected Status 400, but got 200: "200 OK"
	//         -  Expected Body to contain "A string which won't be found", but it was missing
	//         -  Manually Added Error
}

func ExampleFrisby_ExpectJSONType() {
	frisby.Create("Test ExpectJSONType").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJSONType("url", reflect.String).
		PrintReport()

	// Output: Pass  [Test ExpectJSONType]
}

func ExampleFrisby_ExpectJSON() {
	frisby.Create("Test ExpectJSON").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJSON("url", "http://httpbin.org/post").
		ExpectJSON("headers.Accept", "*/*").
		PrintReport()

	// Output: Pass  [Test ExpectJSON]
}

func ExampleFrisby_AfterJSON() {
	frisby.Create("Test AfterJSON").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		AfterJSON(func(F *frisby.Frisby, json *simplejson.Json, err error) {
			val, _ := json.Get("url").String()
			fmt.Println("url =", val)
		})

	// Output: url = http://httpbin.org/post
}

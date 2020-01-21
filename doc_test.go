package frisby_test

import (
	"fmt"
	"reflect"

	"github.com/bitly/go-simplejson"
	"github.com/EducationPlannerBC/frisby"
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

func ExampleFrisby_ExpectJsonType() {
	frisby.Create("Test ExpectJsonType").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJsonType("url", reflect.String).
		PrintReport()

	// Output: Pass  [Test ExpectJsonType]
}

func ExampleFrisby_ExpectJson() {
	frisby.Create("Test ExpectJson").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		ExpectJson("url", "http://httpbin.org/post").
		ExpectJson("headers.Accept", "*/*").
		PrintReport()

	// Output: Pass  [Test ExpectJson]
}

func ExampleFrisby_AfterJson() {
	frisby.Create("Test AfterJson").
		Post("http://httpbin.org/post").
		Send().
		ExpectStatus(200).
		AfterJson(func(F *frisby.Frisby, json *simplejson.Json, err error) {
			val, _ := json.Get("url").String()
			fmt.Println("url =", val)
		})

	// Output: url = http://httpbin.org/post
}

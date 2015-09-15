package frisby

import (
	"fmt"
)

func ExampleGet() {
	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language").
		PrintReport()
}

func ExamplePost() {
	frisby.Create("Test POST").
		Post("http://httpbin.org/post").
		SetData("test_key", "test_value").
		Send().
		ExpectStatus(200).
		PrintReport()

	// Pass  [Test POST]
}

func ExamplePrintReport_Pass() {
	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language").
		PrintReport()

	// Pass  [Test GET Go homepage]
}

func ExamplePrintReport_Fail() {
	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(400).
		ExpectContent("A string which won't be found").
		PrintReport()

	// FAIL  [Test GET Go homepage]
	//         -  Expected Status 400, but got 200: "200 OK"
	//         -  Expected Body to contain "A string which won't be found", but it was missing
}

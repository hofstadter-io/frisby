package main

import (
	"fmt"

	"github.com/verdverm/frisby"
)

func main() {
	fmt.Println("Frisby!\n")

	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(200).
		ExpectContent("The Go Programming Language").
		PrintReport()

	frisby.Create("Test GET Go homepage").
		Get("http://golang.org").
		Send().
		ExpectStatus(400).
		ExpectContent("A string which won't be found").
		PrintReport()

	frisby.Create("Test POST").
		Post("http://httpbin.org/post").
		SetData("test_key", "test_value").
		Send().
		ExpectStatus(200).
		// PrintBody().
		PrintReport()

}

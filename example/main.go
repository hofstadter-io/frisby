package main

import (
	"fmt"
	"aidi"
	"strings"
)
func main() {
	fmt.Println("Aidi!\n")

	//aidi.CreateCase("Test GET Go homepage").
	//	Get("http://golang.org").
	//	Send().
	//	ExpectStatus(200)
	payload := strings.NewReader(`{
	"current_pwd":"123456",
	"new_pwd":"123456"
}`)
	aidi.CreateCase("Test GET Go homepage").
		Put("http://localhost:3000/1/api/v1/account").
		SetBody(payload).
		Send().
		ExpectStatus(200).
		PrintReport()
	aidi.Global.PrintReport()
}

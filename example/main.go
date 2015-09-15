package main

import (
	"fmt"
	"reflect"

	"github.com/verdverm/frisby"
)

func main() {
	fmt.Println("Frisby!\n")

	creds := map[string]string{
		"username": "test",
		"password": "test",
	}
	frisby.Create("Test successful user login").
		Post("http://localhost:8080/auth/login").
		SetJson(creds).
		Send().
		ExpectStatus(200).
		ExpectJsonType("uid", reflect.String).
		ExpectJson("uid", "4fa54e0-1edb-4051-a175-0076603cde7").
		PrintReport()

	bad_username := map[string]string{
		"username": "test2",
		"password": "test",
	}
	frisby.Create("Test bad username login").
		Post("http://localhost:8080/auth/login").
		SetJson(bad_username).
		Send().
		ExpectStatus(400).
		ExpectJson("Error", "login failure").
		PrintReport()

	bad_password := map[string]string{
		"username": "test",
		"password": "test2",
	}
	frisby.Create("Test bad password login").
		Post("http://localhost:8080/auth/login").
		SetJson(bad_password).
		Send().
		ExpectStatus(400).
		ExpectJson("Error", "login failure").
		PrintReport()

}

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EducationPlannerBC/frisby"
	"github.com/bitly/go-simplejson"
)

func getKeyVal(header string) (string, string) {
	var key, val string
	if header != "" {
		res := strings.Split(header, ":")
		if len(res) == 2 {
			key = res[0]
			val = res[1]
		}
	}
	return key, val
}

// Adapted from Jacob Lambert's answer at https://stackoverflow.com/questions/48465575/easy-way-to-split-string-into-map-in-go/48465690
func getKeyValMap(jsondata string) map[string]string {
	entries := strings.Split(jsondata, ",")
	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, ":")
		m[parts[0]] = parts[1]
	}
	return m
}

// From https://golangcookbook.com/chapters/strings/reverse/
func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// For load testing
func getKeyReverseValMap(jsondata string) map[string]string {
	entries := strings.Split(jsondata, ",")
	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, ":")
		m[parts[0]] = reverse(parts[1])
	}
	return m
}

func getMethod(testname string, url string, header1 string, header2 string, header3 string, header4 string, status string, expected string) {
	F := frisby.Create(testname).Get(url)
	var key, val string
	key, val = getKeyVal(header1)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header3)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header4)
	if key != "" {
		F.SetHeader(key, val)
	}
	i, err := strconv.Atoi(status)
	if err != nil {
		log.Fatal(err)
	}
	F.Send().ExpectStatus(i)
	if expected != "" {
		F.ExpectContent(expected)
	}
	// fmt.Println("getMethod response:")
	// F.PrintBody()
	// F.PrintGoTestReport()
}

func postMethod(testname string, url string, header1 string, header2 string, header3 string, jsondata string, status string, expected string) {
	F := frisby.Create(testname).Post(url)
	var key, val string
	key, val = getKeyVal(header1)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header3)
	if key != "" {
		F.SetHeader(key, val)
	}
	if jsondata != "" {
		F.SetJSON(jsondata)
	}
	i, err := strconv.Atoi(status)
	if err != nil {
		log.Fatal(err)
	}
	F.Send().ExpectStatus(i)
	if expected != "" {
		F.ExpectContent(expected)
	}
	// fmt.Println("postMethod response:")
	// F.PrintBody()
}

func putMethod(testname string, url string, header1 string, header2 string, header3 string, jsondata string, status string, expected string) {
	F := frisby.Create(testname).Put(url)
	var key, val string
	key, val = getKeyVal(header1)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header3)
	if key != "" {
		F.SetHeader(key, val)
	}
	if jsondata != "" {
		F.SetJSON(jsondata)
	}
	i, err := strconv.Atoi(status)
	if err != nil {
		log.Fatal(err)
	}
	F.Send().ExpectStatus(i)
	if expected != "" {
		F.ExpectContent(expected)
	}
	fmt.Println("putMethod response:")
	F.PrintBody()
}

func deleteMethod(testname string, url string, header1 string, header2 string, header3 string, jsondata string, status string, expected string) {
	F := frisby.Create(testname).Delete(url)
	var key, val string
	key, val = getKeyVal(header1)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2)
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header3)
	if key != "" {
		F.SetHeader(key, val)
	}
	if jsondata != "" {
		F.SetJSON(jsondata)
	}
	i, err := strconv.Atoi(status)
	if err != nil {
		log.Fatal(err)
	}
	F.Send().ExpectStatus(i)
	if expected != "" {
		F.ExpectContent(expected)
	}
	// fmt.Println("deleteMothod response:")
	// F.PrintBody()
}

func authByGUIDMethod(testname string, url string, bearerToken string, header2 string, header3 string, jsondata string, status string, expected string) {
	var key, val string
	var F *frisby.Frisby
	var GUID string

	var username = "test_" + strconv.FormatInt(time.Now().UnixNano(), 16) + "@example.com"
	account := map[string]string{
		"username":  username,
		"password":  "Abc123",
		"firstName": "FirstName",
		"lastName":  "LastName",
		// "isActive":  true,
		// "root":      true,
	}

	// Create user account
	F = frisby.Create(testname + " Create user account").Post(url)
	key, val = getKeyVal(bearerToken) // bearer token
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2) // content type for POST
	if key != "" {
		F.SetHeader(key, val)
	}
	F.SetJSON(account)
	F.Send().ExpectStatus(200)
	// TODO - handle exception when 'guid' is not in response!
	F.AfterJSON(func(F *frisby.Frisby, json *simplejson.Json, err error) {
		GUID, _ = json.Get("guid").String()
	})
	//F.PrintBody()

	F = frisby.Create(testname + " Get user account by guid").Get(url + "/" + GUID)
	key, val = getKeyVal(bearerToken) // bearer token
	if key != "" {
		F.SetHeader(key, val)
	}
	F.Send().ExpectStatus(200).ExpectContent("FirstName")
	//F.PrintBody()

	// Change firstName and password
	account["firstName"] = "FirstName2"
	account["password"] = "Def456"
	F = frisby.Create(testname + " Change password for guid").Put(url + "/" + GUID)
	F.SetJSON(account)
	key, val = getKeyVal(bearerToken) // bearer token
	if key != "" {
		F.SetHeader(key, val)
	}
	F.Send().ExpectStatus(200).ExpectContent("FirstName2")
	//F.PrintBody()

	F = frisby.Create(testname + " Delete user with guid").Delete(url + "/" + GUID)
	key, val = getKeyVal(bearerToken) // bearer token
	if key != "" {
		F.SetHeader(key, val)
	}
	F.Send().ExpectStatus(200).ExpectContent("deleted")
	//F.PrintBody()
}

func loadMethod(testname string, url string, bearerToken string, header2 string, repsAndSleeps string, jsondata string, status string, expected string) {
	var key, val string
	var F *frisby.Frisby

	config := getKeyValMap(jsondata)
	// Do an update
	F = frisby.Create(testname + " Update").Put(url)
	key, val = getKeyVal(bearerToken) // bearer token
	if key != "" {
		F.SetHeader(key, val)
	}
	key, val = getKeyVal(header2) // content type
	if key != "" {
		F.SetHeader(key, val)
	}
	F.SetJSON(config)
	F.Send().ExpectStatus(200).ExpectContent(expected)
	//F.PrintBody()

	var reps = 100
	var sleep = 0
	if repsAndSleeps != "" {
		res := strings.Split(repsAndSleeps, ",")
		if len(res) == 2 {
			reps, _ = strconv.Atoi(res[0])
			sleep, _ = strconv.Atoi(res[1])
		}
	}

	for i := 1; i <= reps; i++ {
		fmt.Print(".")
		// Get
		F = frisby.Create(testname + " Get").Get(url)
		key, val = getKeyVal(bearerToken) // bearer token
		if key != "" {
			F.SetHeader(key, val)
		}
		key, val = getKeyVal(header2) // content type
		if key != "" {
			F.SetHeader(key, val)
		}
		F.Send().ExpectStatus(200)
		//F.PrintBody()

		revValConfig := getKeyReverseValMap(jsondata)
		// Do another update
		F = frisby.Create(testname + " Update").Put(url)
		key, val = getKeyVal(bearerToken) // bearer token
		if key != "" {
			F.SetHeader(key, val)
		}
		key, val = getKeyVal(header2) // content type
		if key != "" {
			F.SetHeader(key, val)
		}
		F.SetJSON(revValConfig)
		F.Send().ExpectStatus(200)
		//F.PrintBody()

		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

}

func checkRoute(record []string) {
	if len(record) != 9 {
		fmt.Println("The format of the csv file is invalid!")
		os.Exit(1)
	}
	testname := record[0]
	method := record[1]
	url := record[2]
	header1 := record[3]
	header2 := record[4]
	header3 := record[5]
	status := record[7]
	expected := record[8]
	if method == "GET" {
		header4 := record[6]
		getMethod(testname, url, header1, header2, header3, header4, status, expected)
	}
	if method == "POST" {
		jsondata := record[6]
		postMethod(testname, url, header1, header2, header3, jsondata, status, expected)
	}
	if method == "PUT" {
		jsondata := record[6]
		putMethod(testname, url, header1, header2, header3, jsondata, status, expected)
	}
	if method == "DELETE" {
		jsondata := record[6]
		deleteMethod(testname, url, header1, header2, header3, jsondata, status, expected)
	}
	if method == "AUTHGUID" {
		jsondata := record[6]
		authByGUIDMethod(testname, url, header1, header2, header3, jsondata, status, expected)
	}
	if method == "LOAD" {
		jsondata := record[6]
		loadMethod(testname, url, header1, header2, header3, jsondata, status, expected)
	}
}

func iterateRoutes(routesfile string) {

	// Open file describing expected routes
	csvfile, err := os.Open(routesfile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file "+routesfile, err)
	}

	// Parse the file using '|' as default separator and ignoring double-quotes
	r := csv.NewReader(csvfile)
	r.Comma = '|'
	r.LazyQuotes = true

	// Iterate through the routes
	for {
		// Read each route from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(record)
		checkRoute(record)
	}

}

func main() {

	args := os.Args
	argsLen := len(args)
	if argsLen != 2 {
		fmt.Println("Usage: " + args[0] + " <routes file>.csv")
		fmt.Println("\nThe routes file is a delimited text file that uses | to separate values")
		fmt.Println("\nThe format is:")
		fmt.Println("testName|GET|url|header1|header2|header3|header4|responseStatus|expectedOutput")
		fmt.Println("testName|POST|url|header1|header2|header3|jsondata|responseStatus|expectedOutput")
		fmt.Println("testName|PUT|url|header1|header2|header3|jsondata|responseStatus|expectedOutput")
		fmt.Println("testName|DELETE|url|header1|header2|header3|jsondata|responseStatus|expectedOutput")
		fmt.Println("auth by GUID|AUTHGUID|url|bearerToken|header2|header3|jsondata|responseStatus|expectedOutput")
		fmt.Println("\nThe file route_file.csv illustrates the correct format.")
		os.Exit(1)
	}

	routesFile := args[1]

	fmt.Println("Route Testing!")

	iterateRoutes(routesFile)

	frisby.Global.PrintReport()
}

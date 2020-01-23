package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/EducationPlannerBC/frisby"
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
	//F.PrintBody()
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
		fmt.Println("\nThe file route_file.csv illustrates the correct format.")
		os.Exit(1)
	}

	routesFile := args[1]

	fmt.Println("Route Testing!")

	iterateRoutes(routesFile)

	frisby.Global.PrintReport()
}

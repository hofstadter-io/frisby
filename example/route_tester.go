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

func getMethod(testname string, url string, header1 string, header2 string, header3 string, header4 string, status string, expected string) {
	F := frisby.Create(testname).
		Get(url)
	var res []string
	if header1 != "" {
		res = strings.Split(header1, ":")
	}
	if len(res) == 2 {
		F.SetHeader(res[0], res[1])
	}
	if header2 != "" {
		res = strings.Split(header2, ":")
	}
	if len(res) == 2 {
		F.SetHeader(res[0], res[1])
	}
	if header3 != "" {
		res = strings.Split(header3, ":")
	}
	if len(res) == 2 {
		F.SetHeader(res[0], res[1])
	}
	if header4 != "" {
		res = strings.Split(header4, ":")
	}
	if len(res) == 2 {
		F.SetHeader(res[0], res[1])
	}
	i, err := strconv.Atoi(status)
	if err != nil {
		log.Fatal(err)
	}
	if len(res) == 2 {
		F.SetHeader(res[0], res[1])
	}
	F.Send().ExpectStatus(i)
	if expected != "" {
		F.ExpectContent(expected)
	}
	// F.PrintBody()
	// F.PrintGoTestReport()
}

// func postMethod1(url string, contentType string, data string) {
// 	postData := strings.NewReader(data)
// 	response, err := http.Post(url, contentType, postData)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer response.Body.Close()
// 	body, err := ioutil.ReadAll(response.Body)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("%s", body)
// }

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
	header4 := record[6]
	status := record[7]
	expected := record[8]

	if method == "GET" {
		getMethod(testname, url, header1, header2, header3, header4, status, expected)
	}
	if method == "POST" {
		//postMethod1(url, record[2], record[3])
	}
}

func iterateroutes(routesfile string) {

	// Open file describing expected routes
	csvfile, err := os.Open(routesfile)
	if err != nil {
		log.Fatalln("Couldn't open the csv file "+routesfile, err)
	}

	// Parse the file
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
	len := len(args)
	if len != 2 {
		fmt.Println("Usage: " + args[0] + " <routes file>.csv")
		fmt.Println("\nThe routes file is a delimited text file that uses | to separate values")
		fmt.Println("\nThe format is:")
		fmt.Println("testName|method(e.g. GET)|url|header1|header2|header3|header4|responseStatus|expectedOutput")
		fmt.Println("\nThe file route_file.csv illustrates the correct format.")
		os.Exit(1)
	}

	routesFile := args[1]

	fmt.Println("Route Testing!")

	iterateroutes(routesFile)

	frisby.Global.PrintReport()
}

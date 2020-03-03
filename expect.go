package frisby

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
)

// ExpectFunc function type used as argument to Expect()
type ExpectFunc func(F *Frisby) (bool, string)

// Expect checks according to the given function, which allows you to describe any kind of assertion.
func (F *Frisby) Expect(foo ExpectFunc) *Frisby {
	Global.NumAsserts++
	if ok, errStr := foo(F); !ok {
		F.AddError(errStr)
	}
	return F
}

// ExpectStatus checks the response status code
func (F *Frisby) ExpectStatus(code int) *Frisby {
	Global.NumAsserts++
	status := F.Resp.StatusCode
	if status != code {
		errStr := fmt.Sprintf("Expected Status %d, but got %d: %q", code, status, F.Resp.Status)
		F.AddError(errStr)
	}
	return F
}

// ExpectHeader checks for header and if values match
func (F *Frisby) ExpectHeader(key, value string) *Frisby {
	Global.NumAsserts++
	chkVal := F.Resp.Header.Get(key)
	if chkVal == "" {
		errStr := fmt.Sprintf("Expected Header %q, but it was missing", key)
		F.AddError(errStr)
	} else if chkVal != value {
		errStr := fmt.Sprintf("Expected Header %q to be %q, but got %q", key, value, chkVal)
		F.AddError(errStr)
	}
	return F
}

// ExpectContent checks the response body for the given string
func (F *Frisby) ExpectContent(content string) *Frisby {
	Global.NumAsserts++
	text, err := F.Resp.Text()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	contains := strings.Contains(text, content)
	if !contains {
		errStr := fmt.Sprintf("Expected Body to contain '%q' in '%s'", content, text)
		F.AddError(errStr)
	}
	return F
}

// ExpectJSON uses the reflect.DeepEqual to compare the response
// JSON and the supplied JSON for structural and value equality
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJSON(path string, value interface{}) *Frisby {
	Global.NumAsserts++
	simpJSON, err := F.Resp.JSON()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		// Loop over each path item and progress down the json path.
		pathItems := strings.Split(path, Global.PathSeparator)
		for _, segment := range pathItems {
			processed := false
			// If the path segment is an integer, and we're at an array, access the index.
			index, err := strconv.Atoi(segment)
			if err == nil {
				if _, err := simpJSON.Array(); err == nil {
					simpJSON = simpJSON.GetIndex(index)
					processed = true
				}
			}

			if !processed {
				simpJSON = simpJSON.Get(segment)
			}
		}
	}
	json := simpJSON.Interface()

	equal := false
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		val, err := simpJSON.Int()
		if err != nil {
			F.AddError(err.Error())
			return F
		}
		equal = (val == value.(int))
	case reflect.Float64:
		val, err := simpJSON.Float64()
		if err != nil {
			F.AddError(err.Error())
			return F
		}
		equal = (val == value.(float64))
	default:
		equal = reflect.DeepEqual(value, json)
	}

	if !equal {
		errStr := fmt.Sprintf("ExpectJSON equality test failed for %q, got value: %v", path, json)
		F.AddError(errStr)
	}

	return F
}

// ExpectJSONType checks if the types of the response
// JSON and the supplied JSON match
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJSONType(path string, valType reflect.Kind) *Frisby {
	Global.NumAsserts++
	json, err := F.Resp.JSON()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		pathItems := strings.Split(path, Global.PathSeparator)
		json = json.GetPath(pathItems...)
	}

	jsonJSON := json.Interface()

	jsonVal := reflect.ValueOf(jsonJSON)
	if valType != jsonVal.Kind() {
		errStr := fmt.Sprintf("Expect Json %q type to be %q, but got %T", path, valType, jsonJSON)
		F.AddError(errStr)
	}

	return F
}

// ExpectJSONLength checks if the JSON at path
// is an array and has the correct length
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJSONLength(path string, length int) *Frisby {
	Global.NumAsserts++
	json, err := F.Resp.JSON()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		pathItems := strings.Split(path, Global.PathSeparator)
		json = json.GetPath(pathItems...)
	}

	ary, err := json.Array()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	L := len(ary)

	if L != length {
		errStr := fmt.Sprintf("Expect length to be %d, but got %d", length, L)
		F.AddError(errStr)
	}

	return F
}

// AfterContentFunc function type used as argument to AfterContent()
type AfterContentFunc func(F *Frisby, content []byte, err error)

// AfterContent allows you to write your own functions for inspecting the body of the response.
// You are also provided with the Frisby object.
//
// The function signiture is AfterContentFunc
//  type AfterContentFunc func(F *Frisby, content []byte, err error)
//
func (F *Frisby) AfterContent(foo AfterContentFunc) *Frisby {
	content, err := F.Resp.Content()
	foo(F, content, err)
	return F
}

// AfterTextFunc function type used as argument to AfterText()
type AfterTextFunc func(F *Frisby, text string, err error)

// AfterText allows you to write your own functions for inspecting the body of the response.
// You are also provided with the Frisby object.
//
// The function signiture is AfterTextFunc
//  type AfterTextFunc func(F *Frisby, text string, err error)
//
func (F *Frisby) AfterText(foo AfterTextFunc) *Frisby {
	text, err := F.Resp.Text()
	foo(F, text, err)
	return F
}

// AfterJSONFunc function type used as argument to AfterJSON()
type AfterJSONFunc func(F *Frisby, json *simplejson.Json, err error)

// AfterJSON allows you to write your own functions for inspecting the body of the response.
// You are also provided with the Frisby object.
//
// The function signiture is AfterJSONFunc
//  type AfterJSONFunc func(F *Frisby, json *simplejson.Json, err error)
//
// simplejson docs: https://github.com/bitly/go-simplejson
func (F *Frisby) AfterJSON(foo AfterJSONFunc) *Frisby {
	json, err := F.Resp.JSON()
	foo(F, json, err)
	return F
}

// PrintBody prints the body of the response
func (F *Frisby) PrintBody() *Frisby {
	str, err := F.Resp.Text()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	fmt.Println(str)
	return F
}

// PrintReport prints a report for the Frisby Object
//
// If there are any errors, they will all be printed as well
func (F *Frisby) PrintReport() *Frisby {
	if len(F.Errs) == 0 {
		fmt.Printf("Pass  [%s]\n", F.Name)
	} else {
		fmt.Printf("FAIL  [%s]\n", F.Name)
		for _, e := range F.Errs {
			fmt.Println("        - ", e)
		}
	}

	return F
}

// PrintGoTestReport prints a report for the Frisby Object in go_test format
//
// If there are any errors, they will all be printed as well
func (F *Frisby) PrintGoTestReport() *Frisby {
	if len(F.Errs) == 0 {
		fmt.Printf("=== RUN   %s\n--- PASS: %s (%.2fs)\n", F.Name, F.Name, F.ExecutionTime)
	} else {
		fmt.Printf("=== RUN   %s\n--- FAIL: %s (%.2fs)\n", F.Name, F.Name, F.ExecutionTime)
		for _, e := range F.Errs {
			fmt.Println("	", e)
		}
	}
	return F
}

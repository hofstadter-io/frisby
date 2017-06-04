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

// Expect Checks according to the given function, which allows you to describe any kind of assertion.
func (F *Frisby) Expect(foo ExpectFunc) *Frisby {
	Global.NumAsserts++
	if ok, err_str := foo(F); !ok {
		F.AddError(err_str)
	}
	return F
}

// Checks the response status code
func (F *Frisby) ExpectStatus(code int) *Frisby {
	Global.NumAsserts++
	status := F.Resp.StatusCode
	if status != code {
		err_str := fmt.Sprintf("Expected Status %d, but got %d: %q", code, status, F.Resp.Status)
		F.AddError(err_str)
	}
	return F
}

// Checks for header and if values match
func (F *Frisby) ExpectHeader(key, value string) *Frisby {
	Global.NumAsserts++
	chk_val := F.Resp.Header.Get(key)
	if chk_val == "" {
		err_str := fmt.Sprintf("Expected Header %q, but it was missing", key)
		F.AddError(err_str)
	} else if chk_val != value {
		err_str := fmt.Sprintf("Expected Header %q to be %q, but got %q", key, value, chk_val)
		F.AddError(err_str)
	}
	return F
}

// Checks the response body for the given string
func (F *Frisby) ExpectContent(content string) *Frisby {
	Global.NumAsserts++
	text, err := F.Resp.Text()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	contains := strings.Contains(text, content)
	if !contains {
		err_str := fmt.Sprintf("Expected Body to contain %q, but it was missing", content)
		F.AddError(err_str)
	}
	return F
}

// ExpectJson uses the reflect.DeepEqual to compare the response
// JSON and the supplied JSON for structural and value equality
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJson(path string, value interface{}) *Frisby {
	Global.NumAsserts++
	simp_json, err := F.Resp.Json()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		// Loop over each path item and progress down the json path.
		path_items := strings.Split(path, Global.PathSeparator)
		for _, segment := range path_items {
			processed := false
			// If the path segment is an integer, and we're at an array, access the index.
			index, err := strconv.Atoi(segment)
			if err == nil {
				if _, err := simp_json.Array(); err == nil {
					simp_json = simp_json.GetIndex(index)
					processed = true
				}
			}

			if !processed {
				simp_json = simp_json.Get(segment)
			}
		}
	}
	json := simp_json.Interface()

	equal := false
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		val, err := simp_json.Int()
		if err != nil {
			F.AddError(err.Error())
			return F
		} else {
			equal = (val == value.(int))
		}
	case reflect.Float64:
		val, err := simp_json.Float64()
		if err != nil {
			F.AddError(err.Error())
			return F
		} else {
			equal = (val == value.(float64))
		}
	default:
		equal = reflect.DeepEqual(value, json)
	}

	if !equal {
		err_str := fmt.Sprintf("ExpectJson equality test failed for %q, got value: %v", path, json)
		F.AddError(err_str)
	}

	return F
}

// ExpectJsonType checks if the types of the response
// JSON and the supplied JSON match
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJsonType(path string, val_type reflect.Kind) *Frisby {
	Global.NumAsserts++
	json, err := F.Resp.Json()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		path_items := strings.Split(path, Global.PathSeparator)
		json = json.GetPath(path_items...)
	}

	json_json := json.Interface()

	json_val := reflect.ValueOf(json_json)
	if val_type != json_val.Kind() {
		err_str := fmt.Sprintf("Expect Json %q type to be %q, but got %T", path, val_type, json_json)
		F.AddError(err_str)
	}

	return F
}

// ExpectJsonLength checks if the JSON at path
// is an array and has the correct length
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJsonLength(path string, length int) *Frisby {
	Global.NumAsserts++
	json, err := F.Resp.Json()
	if err != nil {
		F.AddError(err.Error())
		return F
	}

	if path != "" {
		path_items := strings.Split(path, Global.PathSeparator)
		json = json.GetPath(path_items...)
	}

	ary, err := json.Array()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	L := len(ary)

	if L != length {
		err_str := fmt.Sprintf("Expect length to be %d, but got %d", length, L)
		F.AddError(err_str)
	}

	return F
}

// function type used as argument to AfterContent()
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

// function type used as argument to AfterText()
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

// function type used as argument to AfterJson()
type AfterJsonFunc func(F *Frisby, json *simplejson.Json, err error)

// AfterJson allows you to write your own functions for inspecting the body of the response.
// You are also provided with the Frisby object.
//
// The function signiture is AfterJsonFunc
//  type AfterJsonFunc func(F *Frisby, json *simplejson.Json, err error)
//
// simplejson docs: https://github.com/bitly/go-simplejson
func (F *Frisby) AfterJson(foo AfterJsonFunc) *Frisby {
	json, err := F.Resp.Json()
	foo(F, json, err)
	return F
}

// Prints the body of the response
func (F *Frisby) PrintBody() *Frisby {
	str, err := F.Resp.Text()
	if err != nil {
		F.AddError(err.Error())
		return F
	}
	fmt.Println(str)
	return F
}

// Prints a report for the Frisby Object
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

// Prints a report for the Frisby Object in go_test format
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

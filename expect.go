package frisby

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/bitly/go-simplejson"
)

// Checks the response status code
func (F *Frisby) ExpectStatus(code int) *Frisby {
	status := F.resp.StatusCode
	if status != code {
		err_str := fmt.Sprintf("Expected Status %d, but got %d: %q", code, status, F.resp.Status)
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
	}
	return F
}

// Checks for header and if values match
func (F *Frisby) ExpectHeader(key, value string) *Frisby {
	chk_val := F.resp.Header.Get(key)
	if chk_val == "" {
		err_str := fmt.Sprintf("Expected Header %q, but it was missing", key)
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
	} else if chk_val != value {
		err_str := fmt.Sprintf("Expected Header %q to be %q, but got %q", key, value, chk_val)
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
	}
	return F
}

// Checks the response body for the given string
func (F *Frisby) ExpectContent(content string) *Frisby {
	text, err := F.resp.Text()
	if err != nil {
		F.errs = append(F.errs, err)
		return F
	}
	contains := strings.Contains(text, content)
	if !contains {
		err_str := fmt.Sprintf("Expected Body to contain %q, but it was missing", content)
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
	}
	return F
}

// ExpectJson uses the reflect.DeepEqual to compare the response
// JSON and the supplied JSON for structural and value equality
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJson(path string, value interface{}) *Frisby {
	simp_json, err := F.resp.Json()
	if err != nil {
		F.errs = append(F.errs, err)
		return F
	}

	if path != "" {
		path_items := strings.Split(path, ".")
		simp_json = simp_json.GetPath(path_items...)
	}
	json := simp_json.Interface()

	equal := false
	switch reflect.ValueOf(value).Kind() {
	case reflect.Int:
		val, err := simp_json.Int()
		if err != nil {
			F.errs = append(F.errs, err)
		} else {
			equal = (val == value.(int))
		}
	case reflect.Float64:
		val, err := simp_json.Float64()
		if err != nil {
			F.errs = append(F.errs, err)
		} else {
			equal = (val == value.(float64))
		}
	default:
		equal = reflect.DeepEqual(value, json)
	}

	if !equal {
		err_str := fmt.Sprintf("ExpectJson equality test failed for %q, got type: %v", path, reflect.TypeOf(json))
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
	}

	return F
}

// ExpectJsonType checks if the types of the response
// JSON and the supplied JSON match
//
// path can be a dot joined field names.
// ex:  'path.to.subobject.field'
func (F *Frisby) ExpectJsonType(path string, val_type reflect.Kind) *Frisby {
	json, err := F.resp.Json()
	if err != nil {
		F.errs = append(F.errs, err)
	}

	if path != "" {
		path_items := strings.Split(path, ".")
		json = json.GetPath(path_items...)
	}

	json_json := json.Interface()

	json_val := reflect.ValueOf(json_json)
	if val_type != json_val.Kind() {
		err_str := fmt.Sprintf("Expect Json %q type to be %q, but got %q", path, val_type, json_val.Type())
		err := errors.New(err_str)
		F.errs = append(F.errs, err)
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
	content, err := F.resp.Content()
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
	text, err := F.resp.Text()
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
	json, err := F.resp.Json()
	foo(F, json, err)
	return F
}

// Prints the body of the response
func (F *Frisby) PrintBody() *Frisby {
	str, err := F.resp.Text()
	if err != nil {
		F.errs = append(F.errs, err)
		fmt.Println("Error: ", err)
	}
	fmt.Println(str)
	return F
}

// Prints a report for the Frisby Object
//
// If there are any errors, they will all be printed as well
func (F *Frisby) PrintReport() *Frisby {
	if len(F.errs) == 0 {
		fmt.Printf("Pass  [%s]\n", F.Name)
	} else {
		fmt.Printf("FAIL  [%s]\n", F.Name)
		for _, e := range F.errs {
			fmt.Println("        - ", e)
		}
	}

	return F
}

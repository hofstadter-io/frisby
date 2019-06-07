package aidi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// ExpectFunc function type used as argument to Expect()
type ExpectFunc func(a *Aidi) (bool, string)

// Expect Checks according to the given function, which allows you to describe any kind of assertion.
func (F *Aidi) Expect(foo ExpectFunc) *Aidi {
	//Global.NumAsserts++
	if ok, err_str := foo(F); !ok {
		F.AddError(err_str)
	}
	return F
}

// Checks for header and if values match
func (a *Aidi) ExpectHeader(key, value string) *Aidi {
	Global.NumAsserts++
	chk_val := a.Resp.Header.Get(key)
	if chk_val == "" {
		err_str := fmt.Sprintf("Expected Header %q, but it was missing", key)
		a.AddError(err_str)
	} else if chk_val != value {
		err_str := fmt.Sprintf("Expected Header %q to be %q, but got %q", key, value, chk_val)
		a.AddError(err_str)
	}
	return a
}

// ExpectBodyJson checks if the body of the response
// equal the bpdyJson
//
// bodyJson the except response body json string
func (a *Aidi) ExpectBodyJson(bodyJson string) *Aidi {
	Global.NumAsserts++

	buf := new(bytes.Buffer)
	buf.ReadFrom(a.Resp.Body)
	s := buf.String()
	equal, err := areEqualJSON(bodyJson, s)
	if err != nil {
		a.AddError(err.Error())
		return a
	}
	if !equal {
		a.AddError(fmt.Sprintf("ExpectBody equality test failed for %s, got value: %s", bodyJson, s))
	}
	return a
}

func areEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string %s :: %s", s1, err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string %s :: %s", s2, err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

// Prints a report for the Aidi Object
//
// If there are any errors, they will all be printed as well
func (a *Aidi) PrintReport() *Aidi {
	if len(a.Errs) == 0 {
		fmt.Printf("Pass  [%s]\n", a.Name)
	} else {
		fmt.Printf("FAIL  [%s]\n", a.Name)
		for _, e := range a.Errs {
			fmt.Println("        - ", e)
		}
	}

	return a
}

// Prints a report for the Aidi Object in go_test format
//
// If there are any errors, they will all be printed as well
func (a *Aidi) PrintGoTestReport() *Aidi {
	if len(a.Errs) == 0 {
		fmt.Printf("=== RUN   %s\n--- PASS: %s (%.2fs)\n", a.Name, a.Name, a.ExecutionTime)
	} else {
		fmt.Printf("=== RUN   %s\n--- FAIL: %s (%.2fs)\n", a.Name, a.Name, a.ExecutionTime)
		for _, e := range a.Errs {
			fmt.Println("	", e)
		}
	}
	return a
}

// Checks the response status code
func (a *Aidi) ExpectStatus(code int) *Aidi {
	Global.NumAsserts++
	status := a.Resp.StatusCode
	if status != code {
		err_str := fmt.Sprintf("Expected Status %d, but got %d: %q", code, status, a.Resp.Status)
		a.AddError(err_str)
	}
	return a
}

package aidi

import (
	"errors"
	"fmt"
)

type global_data struct {
	Errs map[string][]error

	NumRequest int
	NumAsserts int
	NumErrored int

	PrintProgressName bool
	PrintProgressDot  bool

	PathSeparator string
}

const DefaultPathSeparator = "."

func init() {
	Global.Errs = make(map[string][]error, 0)
	Global.PrintProgressDot = true
	Global.PathSeparator = DefaultPathSeparator
}

// Manually add an error, if you need to
func (G *global_data) AddError(name, err_str string) *global_data {
	G.NumErrored++
	err := errors.New(err_str)
	G.Errs[name] = append(G.Errs[name], err)
	return G
}

// Get all errors for the global_data object
//
// This function should be called last
func (G *global_data) Errors() map[string][]error {
	return G.Errs
}

// Prints a report for the FrisbyGlobal Object
//
// If there are any errors, they will all be printed as well
func (G *global_data) PrintReport() *global_data {
	fmt.Printf("\nFor %d requests made\n", G.NumRequest)
	if len(G.Errs) == 0 {
		fmt.Printf("  All tests passed\n")
	} else {
		fmt.Printf("  FAILED  [%d/%d]\n", G.NumErrored, G.NumAsserts)
		for key, val := range G.Errs {
			fmt.Printf("      [%s]\n", key)
			for _, e := range val {
				fmt.Println("        - ", e)
			}
		}
	}

	return G
}

/*
rmparse_test.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package rmparse

import (
	"fmt"
	"os"

	// "fmt"
	"testing"
)

// TestRMParseV6RMFile tests for a remarkable version 3 file
func TestRMParseV6RMFile(t *testing.T) {

	filer, err := os.Open("../testfiles/version6.rm")
	if err != nil {
		t.Errorf("Could not open version6 rm file %v", err)
	}
	defer filer.Close()

	x, err := RMParse(filer)
	if err == nil {
		t.Errorf("expected error for v6 rm file")
	}
	fmt.Printf("%#v\n", x)
}

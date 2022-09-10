/*
pdf_test.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"os"

	colornames "golang.org/x/image/colornames"

	// "fmt"
	"io/ioutil"
	"testing"
)

// Test converting a PDF and associated files
// note that .pdf at the end of the uuid is optional
func TestConvertWithPDF(t *testing.T) {

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	defer os.Remove(tname)

	RM2PDF("../testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf", tname, "", false, []LocalColour{})
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

// Test converting an rm file bundle without a PDF
// the test file is a UUID to indicate the rm bundle of interest
// A template A4 is provided in lieu of a background PDF
func TestConvertWithoutPDF(t *testing.T) {

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	defer os.Remove(tname)

	colours := []LocalColour{
		LocalColour{
			Name:   "darkseagreen",
			Colour: colornames.Darkseagreen,
		},
	}

	RM2PDF("../testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7", tname, "../templates/A4.pdf", false, colours)
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

// Test converting an rm file bundle with an inserted page
func TestConvertWithInsertedPage(t *testing.T) {

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	defer os.Remove(tname)

	colours := []LocalColour{
		LocalColour{
			Name:   "darkseagreen",
			Colour: colornames.Darkseagreen,
		},
	}

	RM2PDF("../testfiles/fbe9f971-03ba-4c21-a0e8-78dd921f9c4c", tname, "../templates/A4.pdf", false, colours)
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

// TestConvertWithLandscape tests converting an rm file bundle in horizontal format
func TestConvertWithLandscape(t *testing.T) {

	testUUID := "e724bba2-266f-434d-aaf2-935d2b405aee"
	template := ""

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	tname = tname + ".pdf"
	defer os.Remove(tname)

	colours := []LocalColour{
		LocalColour{
			Name:   "blue",
			Colour: colornames.Blueviolet,
		},
	}

	RM2PDF("../testfiles/"+testUUID, tname, template, false, colours)
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

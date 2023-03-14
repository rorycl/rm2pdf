/*
pdf_test.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"fmt"
	"os"

	colornames "golang.org/x/image/colornames"

	"io/ioutil"
	"testing"

	"github.com/rorycl/rm2pdf/pdfutil"
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

	RM2PDF("../testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf", tname, "", "", false, []LocalColour{})
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
		{
			Name:   "darkseagreen",
			Colour: colornames.Darkseagreen,
		},
	}

	RM2PDF("../testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7", tname, "../templates/A4.pdf", "", false, colours)
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
		{
			Name:   "darkseagreen",
			Colour: colornames.Darkseagreen,
		},
	}

	RM2PDF("../testfiles/fbe9f971-03ba-4c21-a0e8-78dd921f9c4c", tname, "../templates/A4.pdf", "", false, colours)
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
		{
			Name:   "blue",
			Colour: colornames.Blueviolet,
		},
	}

	RM2PDF("../testfiles/"+testUUID, tname, template, "", false, colours)
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

// TestWithCustomSettings tests setting custom pens
func TestWithCustomSettings(t *testing.T) {

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

	// write custom configuration to temporary file
	configFile, err := ioutil.TempFile("", "config")
	if err != nil {
		t.Error(err)
	}
	cName := configFile.Name()
	os.Remove(cName)
	cName = cName + ".yaml"
	defer os.Remove(cName)

	fo, err := os.Create(cName)
	if err != nil {
		t.Fatalf("could not open file %s for writing", cName)
	}
	_, _ = fo.Write([]byte(`
---
all:
  - pen:     pen
    weight:  standard
    width:   3.0
    color:   red
    opacity: 0.7
`))
	fo.Sync()

	colours := []LocalColour{}
	RM2PDF("../testfiles/"+testUUID, tname, template, cName, false, colours)
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}
}

// TestConvertZip tests converting an rm file bundle from a zip file
func TestConvertZip(t *testing.T) {

	file := "../testfiles/horizontal_rmapi.zip"
	template := ""

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	tname = tname + ".pdf"
	defer os.Remove(tname)

	RM2PDF(file, tname, template, "", false, []LocalColour{})
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}

	thisPDF, err := pdfutil.NewPDFFile(tname)
	if err != nil {
		t.Errorf("could not get pdf info %s", err)
	}
	if fmt.Sprint(thisPDF.Orientation) != "landscape" {
		t.Errorf("pdf orientation not horizontal, got %s", fmt.Sprint(thisPDF.Orientation))
	}
	if thisPDF.Pages != 2 {
		t.Errorf("pdf pages should be 2, got %d", thisPDF.Pages)
	}
}

// TestConvertZipNoMetadata tests converting an rm zip file bundle from
// before 2021 which holds no metadata. See
// https://github.com/rorycl/rm2pdf/issues/9. Thanks to
// https://github.com/qwert2003 for the bug report.
func TestConvertZipNoMetadata(t *testing.T) {

	file := "../testfiles/no-metadata.zip"
	template := ""

	// make temporary file
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
	}
	tname := tmpfile.Name()
	tname = tname + ".pdf"
	// defer os.Remove(tname)

	fmt.Printf(tname)
	RM2PDF(file, tname, template, "", true, []LocalColour{})
	if err != nil {
		t.Errorf("An rm2pdf error occurred: %v", err)
	}

	thisPDF, err := pdfutil.NewPDFFile(tname)
	if err != nil {
		t.Fatalf("could not get pdf info %s", err)
	}
	if fmt.Sprint(thisPDF.Orientation) != "portrait" {
		t.Errorf("pdf orientation not portrait , got %s", fmt.Sprint(thisPDF.Orientation))
	}
	if thisPDF.Pages != 1 {
		t.Errorf("pdf pages should be 1, got %d", thisPDF.Pages)
	}
}

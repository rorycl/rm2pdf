package pdfutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {

	p, err := NewPDFFile("../testfiles/e724bba2-266f-434d-aaf2-935d2b405aee.pdf")
	if err != nil {
		t.Error(err)
	}

	if p.Pages != 2 {
		t.Errorf("pages should equal 2, got %d", p.Pages)
	}
	if fmt.Sprint(p.Orientation) != "landscape" {
		t.Errorf("orientation should be landscape, got %s", fmt.Sprint(p.Orientation))
	}
}

// copyFile copies files
func copyFile(inPath, outPath string) error {
	input, err := ioutil.ReadFile(inPath)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outPath, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func TestRotate(t *testing.T) {

	testPDF := "../testfiles/e724bba2-266f-434d-aaf2-935d2b405aee.pdf"

	tmpfile, err := ioutil.TempFile("", "tmppdfcopy")
	if err != nil {
		t.Error(err)
	}
	tName := tmpfile.Name() + ".pdf"
	os.Remove(tName) // remove file

	p, err := NewPDFFile(testPDF)
	if err != nil {
		t.Error(err)
	}

	err = copyFile(testPDF, tName)
	if err != nil {
		t.Fatal(err)
	}

	pCopy, err := NewPDFFile(tName)

	pCopy.Rotate(-90)

	if p.Orientation == pCopy.Orientation {
		t.Errorf("orig %s should not equal rotated %s\n", fmt.Sprint(p.Orientation), fmt.Sprint(pCopy.Orientation))
	}
	os.Remove(tName) // remove file

}

func TestRotateCopy(t *testing.T) {

	p, err := NewPDFFile("../testfiles/e724bba2-266f-434d-aaf2-935d2b405aee.pdf")
	if err != nil {
		t.Error(err)
	}

	tmpfile, err := ioutil.TempFile("", "tmppdf")
	if err != nil {
		t.Error(err)
	}
	tName := tmpfile.Name() + ".pdf"
	os.Remove(tName) // remove file

	p.RotateCopy(-90, tName)

	tPDF, err := NewPDFFile(tName) // recreates file
	if err != nil {
		t.Error(err)
	}

	if p.Pages != tPDF.Pages {
		t.Errorf("orig %d pages should %d\n", p.Pages, tPDF.Pages)
	}

	if p.Orientation == tPDF.Orientation {
		t.Errorf("orig %s should not equal rotated %s\n", fmt.Sprint(p.Orientation), fmt.Sprint(tPDF.Orientation))
	}
	os.Remove(tName) // remove file

}

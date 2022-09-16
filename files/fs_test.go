package files

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func lenFilePaths(g, e int, t *testing.T) {
	if e != g {
		t.Errorf("Expected %d filepaths, got %d", e, g)
	}
}

func testReadSeek(irs io.ReadSeeker) error {

	if irs == nil {
		return errors.New("read seeker is nil")
	}

	// read once
	buf := make([]byte, 10)
	n, err := irs.Read(buf)
	if err != nil {
		return fmt.Errorf("error when reading from readseeker %v", err)
	}
	if n != 10 {
		return fmt.Errorf("read bytes expected 10 got %d", n)
	}

	// seek
	_, err = irs.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("could not seek back %s", err)
	}

	// read a second time
	buf2 := make([]byte, 10)
	n2, err := irs.Read(buf2)
	if err != nil {
		return fmt.Errorf("error when reading for a second time from readseeker %v", err)
	}
	if n2 != 10 {
		return fmt.Errorf("read bytes (second time) expected 10 got %d", n)
	}
	if string(buf) != string(buf2) {
		return fmt.Errorf("buf %s != buf2 %s", string(buf), string(buf2))
	}
	return nil
}

func TestRmFSZipFile(t *testing.T) {

	zipPath := "../testfiles/horizontal_rmapi.zip"
	template := ""

	rmFS, err := NewZipRmFS(zipPath, template)
	if err != nil {
		t.Fatalf("could not mount zip fs: %s", err)
	}
	err = rmFS.Scan()
	if err != nil {
		t.Fatalf("could not scan zip fs: %s", err)
	}

	expectedID := "e724bba2-266f-434d-aaf2-935d2b405aee.pdf"
	if rmFS.IdentifyPDF(false) != expectedID {
		t.Errorf("identity %s != %s", rmFS.IdentifyPDF(false), expectedID)
	}

	lenFilePaths(len(rmFS.rmFiles), 1, t)

	err = testReadSeek(rmFS.pdfReader)
	if err != nil {
		t.Error(err)
	}
	err = testReadSeek(rmFS.templateReader)
	if err != nil {
		t.Error(err)
	}

}

func TestRmFSDirectory(t *testing.T) {

	dirPath := "../testfiles/"
	basePath := "d34df12d-e72b-4939-a791-5b34b3a810e7"
	template := "../templates/A4.pdf"

	rmFS, err := NewDirRmFS(dirPath, basePath, template)
	if err != nil {
		t.Fatalf("could not mount directory fs: %s", err)
	}
	err = rmFS.Scan()
	if err != nil {
		t.Fatalf("could not scan directory fs: %s", err)
	}

	expectedID := "../templates/A4.pdf"
	if rmFS.IdentifyPDF(true) != expectedID {
		t.Errorf("identity %s != %s", rmFS.IdentifyPDF(true), expectedID)
	}

	lenFilePaths(len(rmFS.rmFiles), 1, t)

	err = testReadSeek(rmFS.pdfReader)
	if err == nil {
		t.Error("should error, as there is no pdf")
	}
	err = testReadSeek(rmFS.templateReader)
	if err != nil {
		t.Error(err)
	}
}

func TestRmFSDirectory2(t *testing.T) {

	dirPath := "../testfiles/"
	basePath := "cc8313bb-5fab-4ab5-af39-46e6d4160df3"
	template := ""

	rmFS, err := NewDirRmFS(dirPath, basePath, template)
	if err != nil {
		t.Fatalf("could not mount directory fs: %s", err)
	}
	err = rmFS.Scan()
	if err != nil {
		t.Fatalf("could not scan directory fs: %s", err)
	}
	expectedID := "cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf"
	if rmFS.IdentifyPDF(false) != expectedID {
		t.Errorf("identity %s != %s", rmFS.IdentifyPDF(false), expectedID)
	}
	lenFilePaths(len(rmFS.rmFiles), 2, t)

	err = testReadSeek(rmFS.pdfReader)
	if err != nil {
		t.Error(err)
	}
	err = testReadSeek(rmFS.templateReader)
	if err != nil {
		t.Error(err)
	}

}

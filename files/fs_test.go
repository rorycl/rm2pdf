package files

import (
	"fmt"
	"testing"
)

func TestRmFSZipFile(t *testing.T) {

	zipPath := "../testfiles/horizontal_rmapi.zip"

	rmFS, err := NewZipRmFS(zipPath)
	if err != nil {
		t.Fatalf("could not mount zip fs: %s", err)
	}
	err = rmFS.Scan()
	if err != nil {
		t.Fatalf("could not scan zip fs: %s", err)
	}
	fmt.Printf("%+v\n", rmFS)
}

func TestRmFSDirectory(t *testing.T) {

	dirPath := "../testfiles/"
	basePath := "d34df12d-e72b-4939-a791-5b34b3a810e7"

	rmFS, err := NewDirRmFS(dirPath, basePath)
	if err != nil {
		t.Fatalf("could not mount directory fs: %s", err)
	}
	err = rmFS.Scan()
	if err != nil {
		t.Fatalf("could not scan directory fs: %s", err)
	}
	fmt.Printf("%+v\n", rmFS)
}

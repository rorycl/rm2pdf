/*
Test programme to use fs.FS filesystem interface for both os and zipfile
file reading operations
*/
package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
)

func main() {

	fsys := os.DirFS(".")
	// example file system list
	listing, err := fs.Glob(fsys, "*.pdf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("file listing:\n%v\n", listing)

	// example file system walk (os)
	fmt.Printf("\nos walk listing:\n")
	fw(fsys)

	fsys2, err := zip.OpenReader("rmapi/test.zip")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// example file system walk (zip)
	fmt.Printf("\nzipfile walk listing:\n")
	fw(fsys2)

}

func fw(fst fs.FS) {

	// printer provides a simple implementation of the fs.Walkdir
	// specification, which is
	// type WalkDirFunc func(path string, d DirEntry, err error) error
	printer := func(path string, d fs.DirEntry, err error) error {
		dirPrinter := func(d fs.DirEntry) string {
			if d.IsDir() {
				return " [d]"
			}
			return ""
		}
		fmt.Printf("%s%s\n", path, dirPrinter(d))
		return nil
	}

	fs.WalkDir(fst, ".", printer)
}

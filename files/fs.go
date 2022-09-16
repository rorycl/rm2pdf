package files

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// https://godocs.io/testing/fstest
// https://github.com/bitfield/tpg-tools/blob/main/7.5/findgo_test.go
// https://bitfieldconsulting.com/golang/filesystems

// rmFileExtensions describe the file extensions to be drawn out of the
// RmFS file structure
var rmFileExtensions = []string{
	".rm", // lines file
	".content",
	".metadata",
	// ".pagedata", // only needed for rmapi zip files?
	".pdf",
	"-metadata.json",
}

type rmFSType int

const (
	zipfile rmFSType = iota
	directory
)

// String returns a string for an rmFSType
func (t rmFSType) String() string {
	return [...]string{"zipfile", "directory"}[t]
}

// RmFS represents a remarkable file system of files
type RmFS struct {
	fs        fs.FS
	fsPath    string
	rmfsType  rmFSType
	filter    string
	filePaths []string
}

// NewZipRmFS mounts an RmFS from a zip file
func NewZipRmFS(zipPath string) (*RmFS, error) {
	rm := RmFS{
		fsPath:   zipPath,
		rmfsType: zipfile,
	}
	z, err := zip.OpenReader(zipPath)
	if err != nil {
		return &rm, err
	}
	type f fs.FS
	rm.fs = f(z)
	return &rm, nil
}

// NewDirRmFS mounts an RmFS from a Directory from which the needed
// files are selected using a filter
func NewDirRmFS(path, basename string) (*RmFS, error) {
	rm := RmFS{
		fsPath:   path,
		filter:   basename,
		rmfsType: directory,
	}
	rm.fs = os.DirFS(path)
	return &rm, nil
}

// String provides a string reprentation of an RmFS filesystem
func (rf *RmFS) String() string {
	t := `RmFS
    Type       : %s
    Mounted at : %s
    Filtered   : %t
    Files :
`
	filtered := func() bool {
		if rf.filter == "" {
			return false
		}
		return true
	}()
	s := fmt.Sprintf(t, rf.rmfsType, rf.fsPath, filtered)
	for _, f := range rf.filePaths {
		s = s + fmt.Sprintf("       %s\n", f)
	}
	return s
}

// Scan scans for the files of interest and stores these in rf.filePaths
func (rf *RmFS) Scan() error {

	err := fs.WalkDir(rf.fs, ".", func(path string, d fs.DirEntry, err error) error {

		// simple filter function
		filtered := func(s string) bool {
			if rf.filter == "" {
				return true
			}
			return strings.Contains(s, rf.filter)
		}

		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !filtered(path) {
			return nil
		}
		for _, e := range rmFileExtensions {
			if strings.HasSuffix(path, e) {
				rf.filePaths = append(rf.filePaths, path)
				break
			}
		}
		return nil
	})
	return err
}

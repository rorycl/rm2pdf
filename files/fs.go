package files

import (
	"archive/zip"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// https://godocs.io/testing/fstest
// https://github.com/bitfield/tpg-tools/blob/main/7.5/findgo_test.go
// https://bitfieldconsulting.com/golang/filesystems

// rmExtensions describe the extensions for remarkable lines and
// metadata files
var rmExtensions = []string{
	".rm", // lines file
	"-metadata.json",
}

// rmAllExtensions describe the file extensions to be drawn out of
// the RmFS file structure
var rmAllExtensions = append(
	rmExtensions,
	[]string{
		".content",
		".metadata",
		// ".pagedata", // only needed for rmapi zip files?
		".pdf",
	}...,
)

// rmFSType describes a type of RmFS type, such as "zipfile" or "directory"
type rmFSType int

const (
	zipfile rmFSType = iota
	directory
)

// String returns a string for an rmFSType
func (t rmFSType) String() string {
	return [...]string{"zipfile", "directory"}[t]
}

// rmFileDesc describes an rm lines file and associated metadata
type rmFileDesc struct {
	rmPath       string
	rm           fs.File
	metadataPath string
	metadata     fs.File
}

// String returns a terse string version of rmFileDesc
func (r rmFileDesc) String() string {
	return fmt.Sprintf(
		"\n         path       %s"+"\n         metadata   %s",
		r.rmPath,
		r.metadataPath,
	)
}

// RmFS represents a remarkable file system of files
type RmFS struct {
	fs         fs.FS
	fsPath     string
	rmfsType   rmFSType
	filter     string
	identifier string

	// metadata
	contentPath  string
	content      fs.File // content
	metadataPath string
	metadata     fs.File // metadata

	// pdf, if any
	pdfPath   string
	pdf       fs.File // backing pdf
	pdfBytes  []byte  // needed to readseek a zip file
	pdfReader io.ReadSeeker

	// template
	templatePath   string
	template       fs.File // either the user provided or embedded template
	templateBytes  []byte  // needed to readseek a zip file
	templateReader io.ReadSeeker

	// per-page rm file metadata and stroke .rm file
	rmFiles map[string]rmFileDesc // base path to rmFileDesc mapping
}

//go:embed "A4.pdf"
var embeddedA4File embed.FS

// NewZipRmFS mounts an RmFS from a zip file
func NewZipRmFS(zipPath, tplPath string) (*RmFS, error) {

	files := make(map[string]rmFileDesc)
	rm := RmFS{
		fsPath:   zipPath,
		rmfsType: zipfile,
		rmFiles:  files,
	}
	var err error
	rm.fs, err = zip.OpenReader(zipPath)
	if err != nil {
		return &rm, err
	}

	if tplPath != "" {
		rm.templatePath = tplPath
		rm.template, err = os.Open(tplPath)
		if err != nil {
			return &rm, err
		}
		err = rm.templateReadSeeker()
		if err != nil {
			return &rm, fmt.Errorf("could not make template readseeker: %w", err)
		}
	} else {
		rm.templatePath = "embedded A4.pdf"
		rm.template, err = embeddedA4File.Open("A4.pdf")
		if err != nil {
			return &rm, fmt.Errorf("could not open embedded template file %s", err)
		}
		err = rm.templateReadSeeker()
		if err != nil {
			return &rm, fmt.Errorf("could not make template readseeker: %w", err)
		}
	}
	return &rm, nil
}

// NewDirRmFS mounts an RmFS from a Directory from which the needed
// files are selected using a filter
func NewDirRmFS(path, basename, tplPath string) (*RmFS, error) {

	files := make(map[string]rmFileDesc)
	rm := RmFS{
		fsPath:   path,
		filter:   basename,
		rmfsType: directory,
		rmFiles:  files,
	}
	rm.fs = os.DirFS(path)

	var err error

	if tplPath != "" {
		rm.templatePath = tplPath
		rm.template, err = os.Open(tplPath)
		if err != nil {
			return &rm, err
		}
		err = rm.templateReadSeeker()
		if err != nil {
			return &rm, fmt.Errorf("could not make template readseeker: %w", err)
		}
	} else {
		rm.templatePath = "embedded A4.pdf"
		rm.template, err = embeddedA4File.Open("A4.pdf")
		if err != nil {
			return &rm, fmt.Errorf("could not open embedded template file %s", err)
		}
		err = rm.templateReadSeeker()
		if err != nil {
			return &rm, fmt.Errorf("could not make template readseeker: %w", err)
		}
	}

	return &rm, nil
}

// IdentifyPDF shows the pdf path in use
func (rf *RmFS) IdentifyPDF(isTpl bool) string {
	if isTpl {
		return rf.templatePath
	}
	return rf.pdfPath
}

// pdfReadSeeker makes the PDF an io.ReadSeeker. Underlying the fs.File,
// an os.File supports seeking, but files from a zip do not, so detect
// that and return a bytes.NewReader if necessary
func (rf *RmFS) pdfReadSeeker() error {
	var (
		err error
		ok  bool
	)
	if rf.pdf == nil {
		errors.New("pdf file has no content, cannot make readseeker")
	}
	if rf.pdfReader, ok = rf.pdf.(io.ReadSeeker); ok {
		return nil
	}
	if rf.pdfReader != nil {
		return nil
	}

	rf.pdfBytes, err = io.ReadAll(rf.pdf)
	if err != nil {
		fmt.Errorf("error reading pdf file, cannot make readseeker bytes: %w", err)
	}
	rf.pdfReader = bytes.NewReader(rf.pdfBytes)
	return nil
}

// templateReadSeeker makes the Template an io.ReadSeeker. Underlying
// the fs.File, an os.File supports seeking, but files from a zip do
// not, so detect that and return a bytes.NewReader if necessary
func (rf *RmFS) templateReadSeeker() error {
	var (
		err error
		ok  bool
	)
	if rf.template == nil {
		errors.New("template file has no content, cannot make readseeker")
	}
	if rf.templateReader, ok = rf.template.(io.ReadSeeker); ok {
		return nil
	}
	if rf.templateReader != nil {
		return nil
	}

	rf.templateBytes, err = io.ReadAll(rf.template)
	if err != nil {
		fmt.Errorf("error reading template file, cannot make readseeker bytes: %w", err)
	}
	rf.templateReader = bytes.NewReader(rf.templateBytes)
	return nil
}

// func identify extracts a uuid identifier from a filesystem
func (rf *RmFS) identify(s string) error {
	if _, err := uuid.Parse(s); err != nil {
		return fmt.Errorf("uuid '%s' is invalid", s)
	}
	rf.identifier = s
	return nil
}

// String provides a string reprentation of an RmFS filesystem
func (rf *RmFS) String() string {
	t := `RmFS
    Identifier    : %s
    Type          : %s
    Mounted at    : %s
    Filtered      : %t
    --
    Content file  : %s
    Metadata file : %s
    PDF file      : %s loaded %t
    RM Files      :
`
	filtered := func() bool {
		if rf.filter == "" {
			return false
		}
		return true
	}()
	hasPDF := func() bool {
		if rf.pdf == nil {
			return false
		}
		return true
	}()
	s := fmt.Sprintf(
		t,
		rf.identifier,
		rf.rmfsType,
		rf.fsPath,
		filtered,
		rf.contentPath,
		rf.metadataPath,
		rf.pdfPath,
		hasPDF,
	)
	for k, f := range rf.rmFiles {
		s = s + fmt.Sprintf("       %s : %s\n", k, f)
	}
	return s
}

// Scan scans for the files of interest and stores the metadata and rm
// lines information in the RmFS struct
func (rf *RmFS) Scan() error {

	var ErrNoSuffix error

	addRMFileData := func(data string) error {

		var suffix string
		for _, thisSuffix := range rmExtensions {
			if strings.Contains(data, thisSuffix) {
				suffix = thisSuffix
				break
			}
		}
		if suffix == "" {
			return ErrNoSuffix
		}

		base := strings.ReplaceAll(data, suffix, "")
		rfd := rf.rmFiles[base]

		var err error

		switch suffix {
		case ".rm":
			rfd.rmPath = data
			rfd.rm, err = rf.fs.Open(data)
			if err != nil {
				return fmt.Errorf("could not open rm file: %w", err)
			}
		case "-metadata.json":
			rfd.metadataPath = data
			rfd.metadata, err = rf.fs.Open(data)
			if err != nil {
				return fmt.Errorf("could not open rm metadata file: %w", err)
			}
		}
		rf.rmFiles[base] = rfd
		return nil
	}

	// find files
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

		// add rm file info
		err = addRMFileData(path)
		switch err {
		case ErrNoSuffix: // not really an error
			break
		case nil: // made an rm record, continue
			return nil
		default:
			return err // an error has occurred
		}

		for _, e := range rmAllExtensions {
			// rmfiles already excluded in addRMFileData
			if strings.HasSuffix(path, e) {
				switch e {
				case ".content":
					rf.contentPath = path
					rf.content, err = rf.fs.Open(path)
					if err != nil {
						return fmt.Errorf("could not open content file: %w", err)
					}
					// extract identifier
					err = rf.identify(strings.TrimSuffix(filepath.Base(path), e))
					if err != nil {
						return err
					}
				case ".metadata":
					rf.metadataPath = path
					rf.metadata, err = rf.fs.Open(path)
					if err != nil {
						return fmt.Errorf("could not open metadata file: %w", err)
					}
				case ".pdf":
					rf.pdfPath = path
					rf.pdf, err = rf.fs.Open(path)
					if err != nil {
						return fmt.Errorf("could not open pdf file: %w", err)
					}
					err = rf.pdfReadSeeker()
					if err != nil {
						return fmt.Errorf("could not make pdf readseeker: %w", err)
					}
				}
				break
			}
		}
		return nil
	})

	return err
}

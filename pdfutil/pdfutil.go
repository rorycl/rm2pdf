// Package pdfutil provides info on and rotates pdf files. The page
// sizing information is simplified to only report on the first page.
package pdfutil

import (
	"errors"
	"fmt"
	"io"
	"os"

	pdfapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

// Orientation determines the orientation of a pdf file
type Orientation int

const (
	landscape Orientation = iota
	portrait
)

// String represents the string of the enum
func (o Orientation) String() string {
	switch o {
	case landscape:
		return "landscape"
	case portrait:
		return "portrait"
	}
	return ""
}

// PDFFile represents a pdf file
type PDFFile struct {
	FilePath    string
	Pages       int
	Width       float64
	Height      float64
	Orientation Orientation
}

// NewPDFFile returns a new PDFFile
func NewPDFFile(path string) (*PDFFile, error) {

	var err error
	p := PDFFile{FilePath: path}

	f, err := os.Open(path)
	if err != nil {
		return &p, err
	}
	defer f.Close()

	p.Pages, err = pdfapi.PageCount(f, nil)
	if err != nil {
		return &p, fmt.Errorf("pagecount error: %s", err)
	}

	// get page sizes
	err = p.dimensions(f)
	if err != nil {
		return &p, err
	}

	return &p, nil
}

func (p *PDFFile) dimensions(f io.ReadSeeker) error {

	var err error
	if f == nil {
		g, err := os.Open(p.FilePath)
		if err != nil {
			return err
		}
		defer g.Close()
		f = io.ReadSeeker(g)
	}

	pageSizes, err := pdfapi.PageDims(f, nil)
	if err != nil {
		return fmt.Errorf("pagesize error: %s", err)
	}
	// use first page only as simplification
	if len(pageSizes) < 1 {
		return errors.New("could not retrieve first page size")
	}
	p.Width = pageSizes[0].Width
	p.Height = pageSizes[0].Height

	if p.Width > p.Height {
		p.Orientation = landscape
	} else {
		p.Orientation = portrait
	}
	return nil
}

// String returns a string representation of a PDFFile
func (p *PDFFile) String() string {
	tpl := `
Filepath    : %s
Pages       : %d
Dimensions  : %0.7f (w) %0.7f (h)
Orientation : %s
`
	return fmt.Sprintf(tpl, p.FilePath, p.Pages, p.Width, p.Height, p.Orientation)
}

func (p *PDFFile) rotate(rotation int, copyFile string) error {
	if rotation == 0 {
		return nil
	}

	err := pdfapi.RotateFile(p.FilePath, copyFile, rotation, nil, nil)
	if err != nil {
		return err
	}

	// update dimensions if in-place copy
	if copyFile == "" {
		err = p.dimensions(nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// Rotate rotates a pdf file by the specified degrees; copyFile can be
// an empty string to overwrite the original file
func (p *PDFFile) Rotate(rotation int) error {
	return p.rotate(rotation, "")
}

// RotateCopy rotates a pdf file by the specified degrees; copyFile can
// be an empty string to overwrite the original file
func (p *PDFFile) RotateCopy(rotation int, copyFile string) error {
	return p.rotate(rotation, copyFile)
}

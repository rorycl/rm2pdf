/*
Collection of information relating to an .rm file bundle.

MIT licensed, please see LICENCE
RCL December 2019
*/

package files

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RMFileInfo is a struct defining the collected metadata about a PDF
// from the reMarkable file collection
type RMFileInfo struct {
	RelPDFPath         string // full relative path to PDF
	RelPDFTemplatePath string // full relative path to PDF template
	Identifier         string // the uuid used to identify the PDF file
	Version            int    // version from metadata
	VisibleName        string // visibleName from metadata (used in reMarkable interface)
	LastModified       time.Time
	OriginalPageCount  int
	PageCount          int
	Pages              []RMPage
	RedirectionPageMap []int // page insertion info
	Debugging          bool
	// show inserted pages
	insertedPages
	// page number used for processing
	thisPageNo int
}

// Debug prints a message if the debugging switch is on
func (r *RMFileInfo) Debug(d string) {
	if r.Debugging {
		fmt.Println(d)
	}
}

// deal with inserted pages
type insertedPages []int

// pageNos shows human-readable inserted page numbers
func (ip insertedPages) insertedPageNos() []int {
	var o []int
	for _, v := range ip {
		o = append(o, v+1)
	}
	return o
}

// format which human readable pages are inserted
func (ip insertedPages) insertedPageNumbers() string {
	var s []string
	for _, n := range ip.insertedPageNos() {
		s = append(s, strconv.Itoa(n))
	}
	o := strings.Join(s, " and ")
	n := strings.Count(o, " and ")
	if n > 1 {
		o = strings.Replace(o, " and ", ", ", n-1)
	}
	return o
}

// register inserted pages
func (r *RMFileInfo) registerInsertedPages() {
	for i, v := range r.RedirectionPageMap {
		if v == -1 {
			r.insertedPages = append(r.insertedPages, i)
		}
	}
	return
}

// PageIterate iterates over pages using the rmfile iterator which
// provides a page number and the pdf to use (either the annotated
// pdf or the template). For annotated pdfs with inserted pages one
// might receive the following output from the iterator:
//
// pageno | pdfPage | inserted | template      |
// -------+---------+----------+---------------+
// 0      | 0       | no       | annotated.pdf |
// 1      | 0       | yes      | template.pdf  |
// 2      | 1       | no       | annotated.pdf |
//
// This function returns 0-indexed pdf pages
func (r *RMFileInfo) PageIterate() (pageNo, pdfPageNo int, inserted, isTemplate bool) {
	pageNo = r.thisPageNo
	r.thisPageNo++

	// if there is only a template, always return the first page
	if r.RelPDFPath == "" {
		pdfPageNo = 0
		isTemplate = true
		return
	}

	// return the template if this is an inserted page
	if r.RedirectionPageMap[pageNo] == -1 {
		pdfPageNo = 0
		inserted = true
		isTemplate = true
		return
	}

	// if the annotated pdf has inserted pages, calculate the offset of
	// the original pdf to use
	if r.PageCount != r.OriginalPageCount {
		pdfPageNo = pageNo
		for i := 0; i <= pageNo; i++ {
			if r.RedirectionPageMap[i] == -1 {
				pdfPageNo--
			}
		}
		return
	}

	// fall through: the annotated pdf has no inserted pages
	pdfPageNo = pageNo
	return

}

// RMPage is a struct defining metadata about each .rm file associated
// with the PDF described in an RMFileInfo. Note that while the .content
// file records page UUIDs for each page of the original PDF, .rm and
// related file are only made for those pages which have marks
type RMPage struct {
	PageNo     int
	Identifier string   // the uuid used to identify the RM file
	RelRMPath  string   // full relative path to the .rm file
	Exists     bool     // file exists on disk
	LayerNames []string // layer names by implicit index
}

// Per-rm file json .metadata file decoding (layers.name)
type rmMetadataLayer struct {
	Layer string `json:"name"`
}

// Per-rm file json .metadata file decoding (layers)
type rmMetadata struct {
	Layers []rmMetadataLayer `json:"layers"`
}

// Per-pdf file .content json file decoding
type content struct {
	FileType           string   `json:fileType`
	Orientation        string   `json:orientation`
	Pages              []string `json:pages`
	RedirectionPageMap []int    `json:redirectionPageMap`
	OriginalPageCount  int      `json:originalPageCount`
	PageCount          int      `json:pageCount`
}

// Per-pdf file .metadata json file decoding: epoch time property
type epochTime time.Time

// Per-pdf file .metadata json file decoding: general metadata
type pdfMetadata struct {
	LastModified epochTime `json:lastmodified`
	Type         string    `json:type`
	Version      int       `json:version`
	VisibleName  string    `json:tpl`
}

// Custom json decoder for unix epochs, with reference to
// https://gist.github.com/alexmcroberts/219127816e7a16c7bd70
func (t *epochTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	eT := time.Unix(q/1000, 0)
	// fmt.Printf("eT, %+v | %s\n", eT, string(eT.Format(time.RFC822)))
	*(*time.Time)(t) = eT
	return
}

// Custom json decoder for unix epochs: string representation
func (t epochTime) String() string {
	return time.Time(t).String()
}

// Custom json decoder for unix epochs: formatter
func (t epochTime) Format(str string) string {
	return time.Time(t).Format(str)
}

// Check if a file exists
func checkFileExists(f string) error {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return err
	}
	return nil
}

// RMFiler collects information from the reMarkable files associated
// with the uuid of interest. Either a pdf at <path/uuid.pdf> is
// expected, or a single A4 page template is to be provided instead. The
// uuid (identified by its filepath plus <uuid>), is used to collect
// information from the .metadata and .content files. It then collects
// layer information for each associated .rm file in a directory named
// by the uuid of the pdf.
func RMFiler(inputpath string, template string) (RMFileInfo, error) {

	rm := RMFileInfo{}

	// if the inputpath has '.pdf' at the end, chop it off
	inputpath = strings.TrimSuffix(inputpath, ".pdf")

	// split path and uuid
	dir, hUUID := filepath.Split(inputpath)

	// verify uuid
	if _, err := uuid.Parse(hUUID); err != nil {
		return rm, fmt.Errorf("uuid '%s' is invalid", hUUID)
	}
	rm.Identifier = hUUID

	// construct paths to .content and .metadata and check the paths exist
	fbase := filepath.Join(dir, hUUID)
	fmetadata := fbase + ".metadata"
	fcontent := fbase + ".content"

	// metadata only exists on xochitl files
	if err := checkFileExists(fmetadata); err == nil {

		body, err := ioutil.ReadFile(fmetadata)
		if err != nil {
			return rm, err
		}
		var p pdfMetadata
		err = json.Unmarshal(body, &p)
		if err != nil {
			return rm, err
		}

		rm.Version = p.Version
		rm.VisibleName = p.VisibleName
		rm.LastModified = time.Time(p.LastModified)
	}

	if err := checkFileExists(fcontent); err != nil {
		return rm, fmt.Errorf("PDF content file %s not found", fcontent)
	}

	// read content
	body, err := ioutil.ReadFile(fcontent)
	if err != nil {
		return rm, err
	}
	var c content
	err = json.Unmarshal(body, &c)
	if err != nil {
		return rm, err
	}

	// load content into rm struct and calculate the inserted pages
	rm.PageCount = c.PageCount
	rm.OriginalPageCount = c.OriginalPageCount
	rm.RedirectionPageMap = c.RedirectionPageMap
	rm.registerInsertedPages()

	if len(c.Pages) != rm.PageCount {
		return rm, fmt.Errorf(
			"number of rm pages %d != json pageCount %d", len(c.Pages), rm.PageCount)
	}

	// check base pdf exists and/or template pdf file
	if c.FileType == "pdf" {
		rm.RelPDFPath = fbase + ".pdf"
		if err := checkFileExists(rm.RelPDFPath); err != nil {
			return rm, fmt.Errorf("PDF file %s not found", rm.RelPDFPath)
		}
	}
	if template != "" {
		err := checkFileExists(template)
		if err != nil {
			return rm, fmt.Errorf("template %s not found", template)
		}
		rm.RelPDFTemplatePath = template
	}

	if rm.RelPDFPath == "" && rm.RelPDFTemplatePath == "" {
		return rm, errors.New("neither a base pdf or template pdf were found")
	}

	// extract each rm json page and construct the path to the .rm file
	// itself
	for i, rmj := range c.Pages {

		if err := checkFileExists(filepath.Join(fbase, rmj+"-metadata.json")); err != nil {
			// swap explicit page uuid for rmapi index-based system
			rmj = strconv.Itoa(i)
		}
		rmJSONPath := filepath.Join(fbase, rmj+"-metadata.json")
		rmPath := filepath.Join(fbase, rmj+".rm")
		rmid := strings.Replace(filepath.Base(rmJSONPath), filepath.Ext(rmJSONPath), "", 1)

		rmP := RMPage{
			PageNo:     i,
			Identifier: rmid,
			RelRMPath:  rmPath,
			Exists:     true,
		}

		err := checkFileExists(rmJSONPath)
		if err != nil {
			rmP.Exists = false
			rm.Pages = append(rm.Pages, rmP)
			continue
		}

		err = checkFileExists(rmPath)
		if err != nil {
			rmP.Exists = false
			rm.Pages = append(rm.Pages, rmP)
			continue
		}

		body, err := ioutil.ReadFile(rmJSONPath)
		if err != nil {
			panic(err)
		}

		// read json from rm .json file
		var m rmMetadata
		err = json.Unmarshal(body, &m)
		if err != nil {
			panic(err)
		}

		for _, v := range m.Layers {
			rmP.LayerNames = append(rmP.LayerNames, v.Layer)
		}

		// append page
		rm.Pages = append(rm.Pages, rmP)

	}

	return rm, nil
}

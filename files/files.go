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
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// RMFileInfo is a struct defining the collected metadata about a PDF
// from the reMarkable file collection
type RMFileInfo struct {
	*RmFS                     // embedded rm filesystem
	Version            int    // version from metadata
	VisibleName        string // visibleName from metadata (used in reMarkable interface)
	LastModified       time.Time
	OriginalPageCount  int
	PageCount          int
	Pages              []RMPage
	Orientation        string
	RedirectionPageMap []int // page insertion info
	// show inserted pages
	insertedPages
	// page number used for processing
	thisPageNo int
	Debugging  bool
}

// Debug prints a message if the debugging switch is on
func (r *RMFileInfo) Debug(d string) {
	if r.Debugging {
		fmt.Println(d)
	}
}

// InsertedPages is a public export of the embedded insertedPages human
// readable page numbers func
func (r *RMFileInfo) InsertedPages() string {
	return r.insertedPages.insertedPageNumbers()
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
//
// Returning an io.ReadSeeker from an fs.File is described by Ian Lance
// Taylor at https://github.com/golang/go/issues/44175#issuecomment-775545730
func (r *RMFileInfo) PageIterate() (pageNo, pdfPageNo int, inserted, isTemplate bool, reader *io.ReadSeeker) {
	pageNo = r.thisPageNo
	r.thisPageNo++

	// if there is only a template, always return the first page
	if r.pdfPath == "" {
		pdfPageNo = 0
		isTemplate = true
		reader = &r.templateReader
		return
	}

	// older remarkable bundles don't report inserted pages; ignore
	hasRedir := func() bool { return len(r.RedirectionPageMap) > 0 }()

	// return the template if this is an inserted page
	if hasRedir && r.RedirectionPageMap[pageNo] == -1 {
		pdfPageNo = 0
		inserted = true
		isTemplate = true
		reader = &r.templateReader
		return
	}

	// remaining target is the annotated file
	reader = &r.pdfReader

	// if the annotated pdf has inserted pages, calculate the offset of
	// the original pdf to use
	if hasRedir && r.PageCount != r.OriginalPageCount {
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
	*rmFileDesc // the rm file descriptor
	PageNo      int
	Identifier  string   // the uuid used to identify the RM file
	Exists      bool     // file exists on disk
	LayerNames  []string // layer names by implicit index
}

// RMFile returns the fs.File pointing to the .rm file
func (r *RMPage) RMFile() fs.File {
	return r.rm
}

// RMFilePath returns the .rm file path
func (r *RMPage) RMFilePath() string {
	return r.rmPath
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
// with the uuid of interest.
//
// Either a pdf at <path/uuid.pdf> is expected, or a single A4 page
// template. If a template is not explictly provided an embedded A4
// template is used. This is managed by fs.go
//
// The uuid (identified by its filepath plus <uuid>), is used to collect
// information from the .metadata and .content files. It then collects
// layer information for each associated .rm file in a directory named
// by the uuid of the pdf.
func RMFiler(inputpath string, template string) (RMFileInfo, error) {

	rm := RMFileInfo{}
	var err error

	// make a remarkable file system of files and scan the file system
	if filepath.Ext(strings.ToLower(inputpath)) == ".zip" {
		rm.RmFS, err = NewZipRmFS(inputpath, template)
		if err != nil {
			return rm, fmt.Errorf("could not init new zip fs: %s", err)
		}
	} else {
		// trim suffix, so if a suffix is provided by mistake it is ignored
		inputpath = strings.TrimSuffix(inputpath, filepath.Ext(inputpath))
		path, base := filepath.Split(inputpath)
		if path == "" {
			path = "."
		}
		rm.RmFS, err = NewDirRmFS(path, base, template)
		if err != nil {
			return rm, fmt.Errorf("could not init new directory fs: %s", err)
		}
	}

	// scan the filesystem for bundle files
	err = rm.Scan()
	if err != nil {
		return rm, err
	}

	// metadata
	// load and read the metadata if the metadata file exists (which it
	// may not do in older, per 2021 rm bundles; see
	// https://github.com/rorycl/rm2pdf/issues/9)
	var body []byte
	if rm.metadata != nil {
		body, err = io.ReadAll(rm.metadata)
		if err != nil {
			return rm, fmt.Errorf("could not read metadata file %s : %s", rm.metadataPath, err)
		}

		var p pdfMetadata
		err = json.Unmarshal(body, &p)
		if err != nil {
			return rm, fmt.Errorf("could not unmarshal %s: %s", rm.metadataPath, err)
		}
		rm.Version = p.Version
		rm.VisibleName = p.VisibleName
		rm.LastModified = time.Time(p.LastModified)
	}

	// content
	// load content into rm struct and calculate the inserted pages
	// assume that if OriginalPageCount is 0 this is from an historic
	// .rm file (which did not have this field) and set it to be the
	// same as PageCount
	if rm.content == nil {
		return rm, errors.New("content file does not exist")
	}
	body, err = io.ReadAll(rm.content)
	if err != nil {
		return rm, err
	}
	var c content
	err = json.Unmarshal(body, &c)
	if err != nil {
		return rm, fmt.Errorf("could not unmarshal content file %s : %s", rm.contentPath, err)
	}
	rm.Orientation = c.Orientation
	rm.PageCount = c.PageCount
	rm.OriginalPageCount = c.OriginalPageCount
	if rm.OriginalPageCount == 0 {
		rm.OriginalPageCount = rm.PageCount
	}
	rm.RedirectionPageMap = c.RedirectionPageMap
	rm.registerInsertedPages()
	if len(c.Pages) != rm.PageCount {
		return rm, fmt.Errorf(
			"number of rm pages %d != json pageCount %d", len(c.Pages), rm.PageCount)
	}

	// note that template switching is done in fs.go

	// extract each rm json page and construct the path to the .rm file
	// itself
	for i, rmj := range c.Pages {

		rmP := RMPage{
			PageNo:     i,
			Identifier: rmj,
			Exists:     true,
		}

		// some rm files described in the content json file don't
		// necessarily get written to disk. If there is no file, set the
		// page.Exists flag to false and continue processing.
		//
		// rmfs.rmFiles map needs a path/uuid to extract the rmFileDesc
		// note, however, that some older pre-2021 rmapi zip files use
		// 0-page indexing instead of using uuids for the rm files.
		ok := false
		var rmfd rmFileDesc
		for _, p := range []string{rmj, strconv.Itoa(i)} {
			lkPath := filepath.Join(rm.identifier, p)
			rmfd, ok = rm.rmFiles[lkPath]
			if ok {
				break
			}
		}
		if !ok {
			rmP.Exists = false
			continue
		}
		rmP.rmFileDesc = &rmfd

		// open and read json from rm .json file
		body, err := io.ReadAll(rmfd.metadata)
		if err != nil {
			return rm, fmt.Errorf("could not read metadata file %s: %s", rmfd.metadataPath, err)
		}
		var m rmMetadata
		err = json.Unmarshal(body, &m)
		if err != nil {
			return rm, fmt.Errorf("could not unmarshal metadata file %s: %s", rmfd.metadataPath, err)
		}
		for _, v := range m.Layers {
			rmP.LayerNames = append(rmP.LayerNames, v.Layer)
		}

		// append page
		rm.Pages = append(rm.Pages, rmP)

	}
	return rm, nil
}

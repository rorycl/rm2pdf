/*
Tests for files.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package files

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func ptime(ti string) time.Time {
	tp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ti)
	if err != nil {
		panic(err)
	}
	return tp
}

// TestFilesXochitlWithPDF tests the xochitl file format for a test with
// a backing pdf
func TestFilesXochitlWithPDF(t *testing.T) {

	template := ""
	rmf, err := RMFiler("../testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf", template)
	if err != nil {
		t.Fatalf("Could not open file %v", err)
	}

	expected := RMFileInfo{
		RmFS: &RmFS{
			pdfPath:    "cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf",
			identifier: "cc8313bb-5fab-4ab5-af39-46e6d4160df3",
		},
		Version:      17,
		VisibleName:  "tpl",
		LastModified: ptime("2019-12-28 23:17:19 +0000 GMT"),
		PageCount:    2,
		Pages: []RMPage{
			{
				PageNo:     0,
				Identifier: "da7f9a41-c2b2-4cbc-9c1b-5a20b5d54224",
				rmFileDesc: &rmFileDesc{rmPath: "cc8313bb-5fab-4ab5-af39-46e6d4160df3/da7f9a41-c2b2-4cbc-9c1b-5a20b5d54224.rm"},
				LayerNames: []string{"Layer 1", "Layer 2 is empty"},
			},
			{
				PageNo:     1,
				Identifier: "7794dbce-2506-4fb0-99fd-9ec031426d57",
				rmFileDesc: &rmFileDesc{rmPath: "cc8313bb-5fab-4ab5-af39-46e6d4160df3/7794dbce-2506-4fb0-99fd-9ec031426d57.rm"},
				LayerNames: []string{"Layer 1", "Layer 2"},
			},
		},
	}

	if rmf.pdfPath != expected.pdfPath {
		t.Errorf("pdfPath got %v wanted %v", rmf.pdfPath, expected.pdfPath)
	}
	if rmf.identifier != expected.identifier {
		t.Errorf("identifier got %v wanted %v", rmf.identifier, expected.identifier)
	}
	if rmf.Version != expected.Version {
		t.Errorf("Version got %v wanted %v", rmf.Version, expected.Version)
	}
	if rmf.VisibleName != expected.VisibleName {
		t.Errorf("VisibleName got %v wanted %v", rmf.VisibleName, expected.VisibleName)
	}
	if rmf.PageCount != expected.PageCount {
		t.Errorf("PageCount got %v wanted %v", rmf.PageCount, expected.PageCount)
	}

	if rmf.Pages[1].PageNo != expected.Pages[1].PageNo {
		t.Errorf("Page two PageNo got %v wanted %v", rmf.Pages[1].PageNo, expected.Pages[1].PageNo)
	}
	if rmf.Pages[1].Identifier != expected.Pages[1].Identifier {
		t.Errorf("Page two Identifier got %v wanted %v", rmf.Pages[1].Identifier, expected.Pages[1].Identifier)
	}
	if rmf.Pages[1].rmPath != expected.Pages[1].rmPath {
		t.Errorf("Page two rmPath got %v wanted %v", rmf.Pages[1].rmPath, expected.Pages[1].rmPath)
	}
	if rmf.Pages[1].LayerNames[1] != expected.Pages[1].LayerNames[1] {
		t.Error("Page two second layer names not the same")
	}

	// https://stackoverflow.com/a/29339052
	redirStdOut := func(log string) string {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w
		// debug!
		rmf.Debug(log)
		w.Close()
		s, _ := ioutil.ReadAll(r)
		r.Close()
		os.Stdout = oldStdout
		return string(s)
	}

	rmf.Debugging = false
	s := redirStdOut("hi")
	if s != "" {
		t.Error("debug should be nil")
	}
	rmf.Debugging = true
	s = redirStdOut("hi")
	if s != "hi\n" {
		t.Errorf("debug got %s not %s", s, "hi")
	}

}

// TestFilesXochitlWithoutPDF tests xochitl format files without a pdf
func TestFilesXochitlWithoutPDF(t *testing.T) {

	template := "../templates/A4.pdf"
	rmf, err := RMFiler("../testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7", template)
	if err != nil {
		t.Fatalf("Could not open file %v", err)
	}

	expected := RMFileInfo{
		RmFS: &RmFS{
			pdfPath:    "", // no pdf
			identifier: "d34df12d-e72b-4939-a791-5b34b3a810e7",
		},
		Version:      0,
		VisibleName:  "toolbox",
		LastModified: ptime("2020-01-05 13:03:52 +0000 GMT"),
		PageCount:    1,
		Pages: []RMPage{
			{
				PageNo:     0,
				Identifier: "2c277cdb-79a5-4f69-b583-4901d944e77e",
				rmFileDesc: &rmFileDesc{rmPath: "d34df12d-e72b-4939-a791-5b34b3a810e7/2c277cdb-79a5-4f69-b583-4901d944e77e.rm"},
				LayerNames: []string{"Layer 1"},
			},
		},
	}

	if rmf.pdfPath != expected.pdfPath {
		t.Errorf("pdfPath got %v wanted %v", rmf.pdfPath, expected.pdfPath)
	}
	if rmf.identifier != expected.identifier {
		t.Errorf("identifier got %v wanted %v", rmf.identifier, expected.identifier)
	}
	if rmf.Version != expected.Version {
		t.Errorf("Version got %v wanted %v", rmf.Version, expected.Version)
	}
	if rmf.VisibleName != expected.VisibleName {
		t.Errorf("VisibleName got %v wanted %v", rmf.VisibleName, expected.VisibleName)
	}
	if rmf.PageCount != expected.PageCount {
		t.Errorf("PageCount got %v wanted %v", rmf.PageCount, expected.PageCount)
	}

	if rmf.Pages[0].PageNo != expected.Pages[0].PageNo {
		t.Errorf("Page one PageNo got %v wanted %v", rmf.Pages[0].PageNo, expected.Pages[0].PageNo)
	}
	if rmf.Pages[0].Identifier != expected.Pages[0].Identifier {
		t.Errorf("Page one Identifier got %v wanted %v", rmf.Pages[0].Identifier, expected.Pages[0].Identifier)
	}
	if rmf.Pages[0].rmPath != expected.Pages[0].rmPath {
		t.Errorf("Page one rmPath got %v wanted %v", rmf.Pages[0].rmPath, expected.Pages[0].rmPath)
	}
	if rmf.Pages[0].LayerNames[0] != expected.Pages[0].LayerNames[0] {
		t.Error("Page one second layer names not the same")
	}
}

// TestInsertedPage checks if an inserted page is detected correctly
func TestInsertedPage(t *testing.T) {

	testUUID := "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c"
	template := "../templates/A4.pdf"

	rmf, err := RMFiler("../testfiles/"+testUUID, template)
	if err != nil {
		t.Fatalf("Could not open file %v", err)
	}

	expected := RMFileInfo{
		RmFS: &RmFS{
			pdfPath:    "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c.pdf",
			identifier: "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c",
		},
		Version:           0,
		VisibleName:       "insert-pages",
		LastModified:      ptime("2022-09-09 14:13:39 +0100 BST"),
		Orientation:       "portrait",
		OriginalPageCount: 2,
		PageCount:         3,
		Pages: []RMPage{
			{
				PageNo:     0,
				Identifier: "fa678373-8530-465d-a988-a0b158d957e4",
				rmFileDesc: &rmFileDesc{rmPath: "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c/fa678373-8530-465d-a988-a0b158d957e4.rm"},
				LayerNames: []string{"Layer 1"},
			},
			{
				PageNo:     1,
				Identifier: "0b8b6e65-926c-4269-9109-36fca8718c94",
				rmFileDesc: &rmFileDesc{rmPath: "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c/0b8b6e65-926c-4269-9109-36fca8718c94.rm"},
				LayerNames: []string{"Layer 1"},
			},
			{
				PageNo:     2,
				Identifier: "e2a69ab6-5c11-42d1-8d2d-9ce6569d9fdf",
				rmFileDesc: &rmFileDesc{rmPath: "fbe9f971-03ba-4c21-a0e8-78dd921f9c4c/e2a69ab6-5c11-42d1-8d2d-9ce6569d9fdf.rm"},
				LayerNames: []string{"Layer 1"},
			},
		},
		RedirectionPageMap: []int{0, -1, 1},
		Debugging:          false,
	}

	opt := cmp.Comparer(func(x, y RMFileInfo) bool {
		if x.pdfPath != y.pdfPath {
			t.Errorf("path %s != %s", x.pdfPath, y.pdfPath)
			return false
		}
		if x.identifier != y.identifier {
			t.Errorf("identifier %s != %s", x.pdfPath, y.pdfPath)
			return false
		}
		if x.Version != y.Version ||
			x.VisibleName != y.VisibleName ||
			x.Orientation != y.Orientation ||
			x.OriginalPageCount != y.OriginalPageCount ||
			x.PageCount != y.PageCount {
			t.Error("version, visiblename, orientation, originalpagecount or pagecount differ")
			return false
		}
		if len(x.RedirectionPageMap) != len(y.RedirectionPageMap) {
			t.Errorf("redirection length %d != %d", len(x.RedirectionPageMap), len(y.RedirectionPageMap))
			return false
		}
		for i, rpm := range x.RedirectionPageMap {
			if rpm != y.RedirectionPageMap[i] {
				t.Errorf("redirection page map %d %d != %d", i, rpm, y.RedirectionPageMap[i])
				return false
			}
		}
		if len(x.Pages) != len(y.Pages) {
			t.Errorf("page lengths different %d != %d", len(x.Pages), len(y.Pages))
			return false
		}
		for i, xPage := range x.Pages {
			yPage := y.Pages[i]
			if xPage.PageNo != yPage.PageNo {
				t.Errorf("page %d != %d", xPage.PageNo, yPage.PageNo)
				return false
			}
			if xPage.Identifier != yPage.Identifier {
				t.Errorf("identifier %s != %s", xPage.Identifier, yPage.Identifier)
				return false
			}
			if xPage.rmPath != yPage.rmPath {
				t.Errorf("rmpath %x != %s", xPage.rmPath, yPage.rmPath)
				return false
			}
			if len(xPage.LayerNames) != len(yPage.LayerNames) {
				t.Errorf("layer len %d != %d", len(xPage.LayerNames), len(yPage.LayerNames))
				return false
			}
		}
		return true
	})

	// if !cmp.Equal(rmf, expected, cmpopts.IgnoreUnexported(rmf), cmpopts.IgnoreInterfaces(struct{ io.Reader }{})) {
	if !cmp.Equal(rmf, expected, opt) {
		t.Errorf("rmf != expected for insert page test")
	}
	if len(expected.Pages) != rmf.PageCount {
		t.Errorf("expected pages %d != rmf pages %d", len(expected.Pages), rmf.PageCount)
	}

	if len(rmf.insertedPages) != 1 || rmf.insertedPages[0] != 1 {
		t.Errorf(
			"inserted pages %v should equal [1]",
			rmf.insertedPages,
		)
	}
	if !cmp.Equal(rmf.insertedPages.insertedPageNos(), []int{2}) {
		t.Errorf(
			"human inserted pages %v should equal {2}",
			rmf.insertedPages.insertedPageNos(),
		)
	}
	if rmf.insertedPages.insertedPageNumbers() != "2" {
		t.Errorf(
			"human inserted pages as string %v should equal '2'",
			rmf.insertedPages.insertedPageNumbers(),
		)
	}

	type iterExpected struct {
		pageNo     int
		pdfPageNo  int
		inserted   bool
		isTemplate bool
	}
	iExpectArray := []iterExpected{
		{0, 0, false, false},
		{1, 0, true, true},
		{2, 1, false, false},
	}

	for i := 0; i < rmf.PageCount; i++ {
		// ignore filehandle in last assignment
		pageNo, pdfPageNo, inserted, isTemplate, _ := rmf.PageIterate()
		j := iterExpected{pageNo, pdfPageNo, inserted, isTemplate}
		e := iExpectArray[i]
		if j.pageNo != e.pageNo ||
			j.pdfPageNo != e.pdfPageNo ||
			j.inserted != e.inserted ||
			j.isTemplate != e.isTemplate {
			t.Errorf("iter i %d expected %+v got %+v", i, e, j)
		}
	}
}

// TestHorizontal checks if a horizontal PDF is detected correctly
func TestHorizontal(t *testing.T) {

	testUUID := "e724bba2-266f-434d-aaf2-935d2b405aee"
	template := ""

	rmf, err := RMFiler("../testfiles/"+testUUID, template)
	if err != nil {
		t.Errorf("Could not open file %v", err)
	}

	if rmf.Orientation != "landscape" {
		t.Errorf("Expected landscape orientation, got %s", rmf.Orientation)
	}
}

// TestExtensionIgnored checks that when providing an input with an extension
// the extension is ignored
func TestExtensionIgnored(t *testing.T) {

	testUUID := "e724bba2-266f-434d-aaf2-935d2b405aee.arbitrary"
	template := ""

	_, err := RMFiler("../testfiles/"+testUUID, template)
	if err != nil {
		t.Errorf("Could not open file %v", err)
	}
}

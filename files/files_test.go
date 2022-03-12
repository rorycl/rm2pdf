/*
Tests for files.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package files

import (
	"testing"
	"time"
	// "fmt"
)

func ptime(ti string) time.Time {
	tp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ti)
	if err != nil {
		panic(err)
	}
	return tp
}

func TestFilesWithPDF(t *testing.T) {

	template := ""
	rmf, err := RMFiler("../testfiles/xochitl/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf", template)
	if err != nil {
		t.Errorf("Could not open file %v", err)
	}

	// fmt.Printf("%+v", rmf)

	expected := RMFileInfo{
		RelPDFPath:   "../testfiles/xochitl/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf",
		Identifier:   "cc8313bb-5fab-4ab5-af39-46e6d4160df3",
		Version:      17,
		VisibleName:  "tpl",
		LastModified: ptime("2019-12-28 23:17:19 +0000 GMT"),
		PageCount:    2,
		Pages: []RMPage{
			{
				PageNo:     0,
				Identifier: "da7f9a41-c2b2-4cbc-9c1b-5a20b5d54224-metadata",
				RelRMPath:  "../testfiles/xochitl/cc8313bb-5fab-4ab5-af39-46e6d4160df3/da7f9a41-c2b2-4cbc-9c1b-5a20b5d54224.rm",
				LayerNames: []string{"Layer 1", "Layer 2 is empty"},
			},
			{
				PageNo:     1,
				Identifier: "7794dbce-2506-4fb0-99fd-9ec031426d57-metadata",
				RelRMPath:  "../testfiles/xochitl/cc8313bb-5fab-4ab5-af39-46e6d4160df3/7794dbce-2506-4fb0-99fd-9ec031426d57.rm",
				LayerNames: []string{"Layer 1", "Layer 2"},
			},
		},
	}

	if rmf.RelPDFPath != expected.RelPDFPath {
		t.Errorf("RelPDFPath wanted %v got %v", rmf.RelPDFPath, expected.RelPDFPath)
	}
	if rmf.Identifier != expected.Identifier {
		t.Errorf("Identifier wanted %v got %v", rmf.Identifier, expected.Identifier)
	}
	if rmf.Version != expected.Version {
		t.Errorf("Version wanted %v got %v", rmf.Version, expected.Version)
	}
	if rmf.VisibleName != expected.VisibleName {
		t.Errorf("VisibleName wanted %v got %v", rmf.VisibleName, expected.VisibleName)
	}
	if rmf.PageCount != expected.PageCount {
		t.Errorf("PageCount wanted %v got %v", rmf.PageCount, expected.PageCount)
	}

	if rmf.Pages[1].PageNo != expected.Pages[1].PageNo {
		t.Errorf("Page two PageNo wanted %v got %v", rmf.Pages[1].PageNo, expected.Pages[1].PageNo)
	}
	if rmf.Pages[1].Identifier != expected.Pages[1].Identifier {
		t.Errorf("Page two Identifier wanted %v got %v", rmf.Pages[1].Identifier, expected.Pages[1].Identifier)
	}
	if rmf.Pages[1].RelRMPath != expected.Pages[1].RelRMPath {
		t.Errorf("Page two RelRMPath wanted %v got %v", rmf.Pages[1].RelRMPath, expected.Pages[1].RelRMPath)
	}
	if rmf.Pages[1].LayerNames[1] != expected.Pages[1].LayerNames[1] {
		t.Error("Page two second layer names not the same")
	}

}

func TestFilesWithoutPDF(t *testing.T) {

	template := "../templates/A4.pdf"
	rmf, err := RMFiler("../testfiles/xochitl/d34df12d-e72b-4939-a791-5b34b3a810e7", template)
	if err != nil {
		t.Errorf("Could not open file %v", err)
	}

	// fmt.Printf("%+v", rmf)

	expected := RMFileInfo{
		RelPDFPath:   "../templates/A4.pdf",
		Identifier:   "d34df12d-e72b-4939-a791-5b34b3a810e7",
		Version:      0,
		VisibleName:  "toolbox",
		LastModified: ptime("2020-01-05 13:03:52 +0000 GMT"),
		PageCount:    1,
		Pages: []RMPage{
			{
				PageNo:     0,
				Identifier: "2c277cdb-79a5-4f69-b583-4901d944e77e-metadata",
				RelRMPath:  "../testfiles/xochitl/d34df12d-e72b-4939-a791-5b34b3a810e7/2c277cdb-79a5-4f69-b583-4901d944e77e.rm",
				LayerNames: []string{"Layer 1"},
			},
		},
	}

	if rmf.RelPDFPath != expected.RelPDFPath {
		t.Errorf("RelPDFPath wanted %v got %v", rmf.RelPDFPath, expected.RelPDFPath)
	}
	if rmf.Identifier != expected.Identifier {
		t.Errorf("Identifier wanted %v got %v", rmf.Identifier, expected.Identifier)
	}
	if rmf.Version != expected.Version {
		t.Errorf("Version wanted %v got %v", rmf.Version, expected.Version)
	}
	if rmf.VisibleName != expected.VisibleName {
		t.Errorf("VisibleName wanted %v got %v", rmf.VisibleName, expected.VisibleName)
	}
	if rmf.PageCount != expected.PageCount {
		t.Errorf("PageCount wanted %v got %v", rmf.PageCount, expected.PageCount)
	}

	if rmf.Pages[0].PageNo != expected.Pages[0].PageNo {
		t.Errorf("Page one PageNo wanted %v got %v", rmf.Pages[0].PageNo, expected.Pages[0].PageNo)
	}
	if rmf.Pages[0].Identifier != expected.Pages[0].Identifier {
		t.Errorf("Page one Identifier wanted %v got %v", rmf.Pages[0].Identifier, expected.Pages[0].Identifier)
	}
	if rmf.Pages[0].RelRMPath != expected.Pages[0].RelRMPath {
		t.Errorf("Page one RelRMPath wanted %v got %v", rmf.Pages[0].RelRMPath, expected.Pages[0].RelRMPath)
	}
	if rmf.Pages[0].LayerNames[0] != expected.Pages[0].LayerNames[0] {
		t.Error("Page one second layer names not the same")
	}
}

/*
rmparse_test.go
MIT licenced, please see LICENCE
RCL January 2020
*/

package rmparse

import (
	"os"
	// "fmt"
	"testing"
)

func Last(i []Segment) Segment {
	return i[len(i)-1]
}

func TestRMParse(t *testing.T) {

	filer, err := os.Open("../testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3/da7f9a41-c2b2-4cbc-9c1b-5a20b5d54224.rm")
	if err != nil {
		t.Errorf("Could not open file %v", err)
	}
	defer filer.Close()

	rm, err := RMParse(filer)
	if err != nil {
		panic(err)
	}

	lastPath := RMPath{}
	for rm.Parse() {
		lastPath = rm.Path
	}

	// fmt.Printf("%+v", lastPath)
	// {Layer:1 Path:{Pen:17 Colour:0 _:0 Width:2 _:0 NumSegments:226} Segments:[... {X:1033.4183 Y:1429.1265 Pressure:0.33935595 Tilt:0.35699552 _:0 _:0}]}
	thisPath := RMPath{
		Layer: 1,
		Path: Path{
			Pen:         17,
			Colour:      0,
			Width:       2,
			NumSegments: 226,
		},
		Segments: []Segment{
			// only the last segment
			Segment{
				X:        1033.4183,
				Y:        1429.1265,
				Pressure: 0.33935595,
				Tilt:     0.35699552,
			},
		},
	}

	if lastPath.Layer != thisPath.Layer {
		t.Errorf("Layer not %v", thisPath.Layer)
	}
	if lastPath.Path.Pen != thisPath.Path.Pen {
		t.Errorf("Path.Pen not %v", thisPath.Path.Pen)
	}
	if lastPath.Path.Colour != thisPath.Path.Colour {
		t.Errorf("Path.Colour not %v", thisPath.Path.Colour)
	}
	if lastPath.Path.Width != thisPath.Path.Width {
		t.Errorf("Path.Width not %v", thisPath.Path.Width)
	}
	if lastPath.Path.NumSegments != thisPath.Path.NumSegments {
		t.Errorf("Path.NumSegments not %v", thisPath.Path.NumSegments)
	}

	if Last(lastPath.Segments).X != Last(thisPath.Segments).X {
		t.Errorf("Last(Segments).X not %v", Last(thisPath.Segments).X)
	}
	if Last(lastPath.Segments).Y != Last(thisPath.Segments).Y {
		t.Errorf("Last(Segments).Y not %v", Last(thisPath.Segments).Y)
	}
	if Last(lastPath.Segments).Pressure != Last(thisPath.Segments).Pressure {
		t.Errorf("Last(Segments).Pressure not %v", Last(thisPath.Segments).Pressure)
	}
	if Last(lastPath.Segments).Tilt != Last(thisPath.Segments).Tilt {
		t.Errorf("Last(Segments).Tilt not %v", Last(thisPath.Segments).Tilt)
	}

	// fmt.Printf("%+v", rm)
	// &{File:0xc00008c560 Header:[..] LayerNo:2 ThisLayer:2 PathNo:0 ThisPath:1 Path:{Layer:0 Path:{Pen:0 Colour:0 _:0 Width:0 _:0 NumSegments:0} Segments:[]} MaxCoordinates:{X:1404.1321 Y:1873.1632} Verbose:false}--- FAIL: TestMain (0.02s)
	// MaxCoordinates:{X:1404.1321 Y:1873.1632} Verbose:false}
	thisRM := RMFile{
		LayerNo:   2,
		ThisLayer: 2,
		PathNo:    0,
		ThisPath:  1,
		Path: RMPath{
			Layer: 0,
			Path: Path{
				Pen:         0,
				Colour:      0,
				Width:       0,
				NumSegments: 0,
			},
			Segments: []Segment{},
		},
		MaxCoordinates: MaxCoordinates{
			X: 1404.1321,
			Y: 1873.1632,
		},
	}
	if rm.LayerNo != thisRM.LayerNo {
		t.Errorf("Layer number not %v", thisRM.LayerNo)
	}
	if rm.PathNo != thisRM.PathNo {
		t.Errorf("Path number not %v", thisRM.PathNo)
	}
	if rm.ThisPath != thisRM.ThisPath {
		t.Errorf("ThisPath not %v", thisRM.ThisPath)
	}
	if rm.MaxCoordinates != thisRM.MaxCoordinates {
		t.Errorf("MaxCoordinates not %+v", thisRM.MaxCoordinates)
	}

}

func TestRMParseCorrupt(t *testing.T) {

	filer, err := os.Open("../testfiles/7cbc50c9-8d68-48cf-8f77-e70f2e87b732.rm")
	if err != nil {
		t.Errorf("Could not open corrupt rm file %v", err)
	}
	defer filer.Close()

	rm, err := RMParse(filer)
	if err != nil {
		t.Error("The corrupted rm file could not be setup for parsing")
	}
	ok := true
	for rm.Parse() {
		ok = false // should not get here
	}
	if ok != true {
		t.Error("The corrupted rm file could not be parsed")
	}

}

/*
Go port of reMarkable tablet "lines" or ".rm" file parser, with binary
decoding hints drawn from rm2svg
https://github.com/reHackable/maxio/blob/master/tools/rM2svg which in
turn refers to https://github.com/lschwetlick/maxio/tree/master/tools.

Python struct format codes referred to below, such as "<{}sI" are from
rm2svg.

This package provides a python-like iterator based on bufio.Scan, which
iterates over the referenced reMarkable .rm file returning a data
structure consisting of each path with its associated layer and path
segments.

MIT licensed, please see LICENCE
RCL January 2020
*/

package rmparse

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
)

// Header is an rm file header
var Header = "reMarkable .lines file, version=5         "

// RMFile is the reMarkable .rm file File parser metadata base structure
type RMFile struct {
	File           fs.File
	Header         [43]byte
	LayerNo        uint32
	ThisLayer      uint32
	PathNo         uint32
	ThisPath       uint32
	Path           RMPath
	MaxCoordinates MaxCoordinates
	Verbose        bool
}

// RMPath is the reMarkable parsed data structure, returned by Parse()
type RMPath struct {
	Layer    uint32
	Path     Path
	Segments []Segment
}

// HeaderLayers is the header and number of layers
// format <{}sI
type HeaderLayers struct {
	Header [43]byte
	Layers uint32
}

// Paths is the number of paths in this layer
// format <I
type Paths struct {
	Number uint32
}

// Path describes a path
// format <IIIfII
type Path struct {
	Pen         uint32
	Colour      uint32
	_           uint32 // unknown
	Width       float32
	_           uint32 // unknown
	NumSegments uint32
}

// Segment describes a path segment
// format <ffffff
type Segment struct {
	X        float32
	Y        float32
	Pressure float32
	Tilt     float32
	_        float32 // unknown
	_        float32 // unknown
}

// MaxCoordinates stores the maximum and minimum path segments
type MaxCoordinates struct {
	X float32
	Y float32
}

var ms = MaxCoordinates{}

// RMParse instantiates a parser by registering a file to parse and
// initialising the header, layer count and related counters. Continue
// parsing using the "Parse()" iterator-type function.
func RMParse(f fs.File) (*RMFile, error) {

	rm := &RMFile{}
	rm.File = f

	headerLayers, err := HeaderParse(f)
	if err != nil {
		return nil, err
	}

	rm.Header = headerLayers.Header
	rm.LayerNo = headerLayers.Layers

	// last byte is 0-terminated, chop it off
	if string(rm.Header[:len(rm.Header)-1]) != Header {
		return nil, fmt.Errorf("Header does not match %s", Header)
	}

	if rm.LayerNo < 1 {
		return nil, errors.New("Number of layers less than 1")
	}

	// init counters
	rm.ThisLayer = 1
	rm.ThisPath = 1

	return rm, nil
}

// Parse an .rm file, returning an RMPath data structure until depleted.
// The Parse() function collects all the segments in a path and collects
// it in an RMFile.Path struct, stored in RMFile.Path. The Parse
// function is based loosely on bufio.Scan() so it may be called using
// "for" as follows:
//
//    rm := rmparse.RMParse(filename)
//    for rm.Parse() {
//        path = rm.Path
//    }
//
// Catching io.ErrUnexpectedEOF errors is necessary because it is
// possible to have corrupt .rm files which, for example, report layers
// but do not have any content.
func (rm *RMFile) Parse() bool {

	if rm.PathNo == 0 {
		// find number of paths in this layer
		paths, err := ParseLayers(rm.File)
		if err != nil {
			panic("Cannot determine number of paths")
		}

		rm.PathNo = paths.Number
	}

	// init return structure
	rm.Path = RMPath{}

	// complete processing
	if rm.ThisPath > rm.PathNo && rm.ThisLayer == rm.LayerNo {
		rm.MaxCoordinates = ms // record max segment in rm struct
		return false
	}

	// some layers are empty
	if rm.PathNo > 0 {

		// get next path
		path, err := ParsePath(rm.File)
		if err != nil {
			panic("ParsePath failed")
		}

		rm.Path.Layer = rm.ThisLayer
		rm.Path.Path = path

		// retrieve segments
		for s := 1; s <= int(rm.Path.Path.NumSegments); s++ {
			segment, err := ParseSegment(rm.File)
			if err != nil {
				panic("ParseSegment failed")
			}
			rm.Path.Segments = append(rm.Path.Segments, segment)
		}

	}

	// increment path counter
	rm.ThisPath++

	// advance to next layer if necessary
	if rm.ThisPath > rm.PathNo && rm.ThisLayer < rm.LayerNo {
		rm.ThisLayer++
		rm.ThisPath = 1
		rm.PathNo = 0
	}

	return true
}

// HeaderParse starts parsing an .rm file, returning the header and number of layers
func HeaderParse(f fs.File) (HeaderLayers, error) {

	hl := HeaderLayers{}

	err := binary.Read(f, binary.LittleEndian, &hl)
	if err == io.ErrUnexpectedEOF {
		return hl, err
	} else if err != nil {
		panic(err)
	}

	return hl, nil
}

// ParseLayers returns the number of paths for each layer in the .rm
// file
func ParseLayers(f fs.File) (Paths, error) {

	pths := Paths{}

	err := binary.Read(f, binary.LittleEndian, &pths)
	if err == io.ErrUnexpectedEOF {
		return pths, err
	} else if err != nil {
		panic(err)
	}

	return pths, nil
}

// ParsePath returns the path for each path in the layer.paths
func ParsePath(f fs.File) (Path, error) {

	path := Path{}

	err := binary.Read(f, binary.LittleEndian, &path)
	if err == io.ErrUnexpectedEOF {
		return path, err
	} else if err != nil {
		fmt.Printf("%T %+v", err, err)
		panic(err)
	}

	return path, nil
}

// ParseSegment returns the segment for each segment in a path
func ParseSegment(f fs.File) (Segment, error) {

	sg := Segment{}

	err := binary.Read(f, binary.LittleEndian, &sg)
	if err == io.EOF {
		return sg, err
	} else if err == io.ErrUnexpectedEOF {
		return sg, err
	} else if err != nil {
		panic(err)
	}

	// record maximum segment coordinates
	if sg.X > ms.X {
		ms.X = sg.X
	}
	if sg.Y > ms.Y {
		ms.Y = sg.Y
	}

	return sg, nil
}

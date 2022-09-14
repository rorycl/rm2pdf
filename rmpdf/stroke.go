/*
Define settings for describing .rm file paths as PDF strokes using the
go fpdf package.

MIT licenced, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"image/color"
	"math"

	colornames "golang.org/x/image/colornames"
)

// StrokeSetting describes strokes from an .rm file in a pdf document.
// Although Colours are defined as RGBA values, they all have solid
// (255) Alpha values. The width of each stroke is a value representing
// the medium-sized pen width of each pen type (the middle of three
// values), although StdWidth is an eyeballed/very approximate value
// which is further adjusted through StrokeSetting.Width(). In future
// it may be better to set the widths explictly in this struct.
// The Alpha value is set separately using the Opacity value. The
// ColourOverride property determines if the colour of the stroke may be
// manually overridden by command-line options.
type StrokeSetting struct {
	Colour         color.RGBA
	StdWidth       float32
	Opacity        float64
	ColourOverride bool
}

// StrokeMap is a Map of pen numbers in a reMarkable binary .rm file
var StrokeMap = map[int]string{
	2: "pen",
	4: "fineliner",
	3: "marker",
	5: "highligher",
	6: "eraser",
	7: "sharp pencil",
	8: "erase area",
	// v5 pen new types?
	12: "paint",
	13: "mechanical pencil",
	14: "pencil",
	15: "ballpoint",
	16: "marker",
	17: "pen",
	18: "highlighter",
}

// StrokeSettings sets the pen default settings
var StrokeSettings = map[string]StrokeSetting{
	"pen": {
		Colour:         colornames.Black,
		StdWidth:       2.0,
		Opacity:        1,
		ColourOverride: true,
	},
	"highlighter": {
		Colour:         colornames.Blue,
		StdWidth:       15.0,
		Opacity:        0.4,
		ColourOverride: true,
	},
	"fineliner": {
		Colour:         colornames.Blue,
		StdWidth:       1.0,
		Opacity:        1,
		ColourOverride: true,
	},
	"marker": {
		Colour:         colornames.Black,
		StdWidth:       3.8,
		Opacity:        1,
		ColourOverride: true,
	},
	"ballpoint": {
		// Colour  : color.RGBA{68, 68, 68, 225}, // greyish
		Colour:   colornames.Slategray,
		StdWidth: 1.75,
		Opacity:  0.8,
	},
	"pencil": {
		Colour:   colornames.Black,
		StdWidth: 1.9,
		Opacity:  1,
	},
	"mechanical pencil": {
		Colour:   colornames.Black,
		StdWidth: 1.2,
		Opacity:  0.7,
	},
	"paint": {
		Colour:   color.RGBA{55, 55, 55, 220}, // dark grey
		StdWidth: 4.8,
		Opacity:  0.8,
	},
	"eraser": {
		Colour:   colornames.White,
		StdWidth: 9.0,
		Opacity:  0,
	},
	"erase area": {
		Colour:   colornames.White,
		StdWidth: 9.0,
		Opacity:  0,
	},
}

// Width sets pen widths. Each rm pen comes in three widths, 1.875,
// 2.000, 2.125, so provide a fractional width calculation done by
// eyeballing what seems about right. It probably makes sense to move
// the widths to the map of pens in future.
func (s *StrokeSetting) Width(penwidth float32) float64 {
	r := 0.0
	t := float64(s.StdWidth)
	p := float64(penwidth)

	switch math.Round(p*1000) / 1000 {
	case 1.875:
		r = 0.60 * t
	case 2.125:
		r = 1.20 * t
	default:
		r = 0.85 * t
	}
	return float64(r)
}

// NaturalWidth reports pen widths as "narrow", "standard" or "broad"
func (s *StrokeSetting) NaturalWidth(penwidth float32) string {

	p := float64(penwidth)

	switch math.Round(p*1000) / 1000 {
	case 1.875:
		return "narrow"
	case 2.125:
		return "broad"
	}
	return "standard"
}

// Return the rbg components of the stroke's colour
func (s *StrokeSetting) toRGB() (int, int, int) {
	r := int(s.Colour.R)
	g := int(s.Colour.G)
	b := int(s.Colour.B)
	return int(r), int(g), int(b)
}

// Given a colour, determine if the stroke is overrideable (using the
// ColourOverride attribute); if so return the RGB of the
// given colour, else return the RGB of the native colour
func (s *StrokeSetting) selectColour(lc *LocalColour, force bool) (int, int, int) {
	c := color.RGBA{}
	if lc.Name == "" || lc.Name == "empty" {
		c = s.Colour
	} else if !s.ColourOverride && !force {
		c = s.Colour
	} else {
		c = lc.Colour
	}
	r := int(c.R)
	g := int(c.G)
	b := int(c.B)
	return int(r), int(g), int(b)
}

// Return the cmyk components of the stroke's colour
// func (s *StrokeSetting) toCMYK() color.CMYK {
// 	return color.CMYKModel.Convert(s.Colour)
// }

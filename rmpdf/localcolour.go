/*
Local Colour struct, used for custom layer colour specification

MIT licenced, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"image/color"
	"strings"

	colornames "golang.org/x/image/colornames"
)

// LocalColour describes a color by name and RGBA value
type LocalColour struct {
	Name   string
	Colour color.RGBA
}

// UnmarshalFlag generates the colour value for a colour string
func (l *LocalColour) UnmarshalFlag(value string) error {
	// empty or mismatched values are defaulted to an empty color.RGBA
	// instance
	if value == "" {
		value = "empty"
	}
	c := LocalColour{
		Name:   value,
		Colour: colornames.Map[strings.ToLower(value)],
	}
	*l = c
	return nil
}

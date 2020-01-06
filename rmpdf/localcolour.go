/*
Local Colour struct, used for custom layer colour specification

MIT licenced, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"image/color"
	colornames "golang.org/x/image/colornames"
	"strings"
)

type LocalColour struct {
	Name   string
	Colour color.RGBA
}

func (c *LocalColour) Usage() string {
	return "some help here"
}

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

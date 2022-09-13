/*
LocalPenColour set pen colours from an rgba, rgb or hex colour string or from
an image/color colourname as set out in image/colornames.
*/

package penconfig

import (
	"fmt"
	"image/color"
	"strings"

	playcolors "github.com/go-playground/colors"
	colornames "golang.org/x/image/colornames"
)

// LocalPenColour describes a color by name and RGBA value
type LocalPenColour struct {
	Name   string
	Colour color.RGBA
}

// Unmarshal generates the colour value for a colour string
func (l *LocalPenColour) Unmarshal(value string) error {

	var c color.RGBA
	if len(value) > 4 && value[0:4] == "rgba" {
		pColour, err := playcolors.ParseRGBA(value)
		if err != nil {
			return fmt.Errorf("rgba value %s invalid: %s", value, err)
		}
		tmp := pColour.ToRGBA()
		// ignore alpha channel
		c = color.RGBA{R: tmp.R, G: tmp.G, B: tmp.B}

	} else if len(value) > 3 && value[0:3] == "rgb" {
		pColour, err := playcolors.ParseRGB(value)
		if err != nil {
			return fmt.Errorf("rgb value %s invalid: %s", value, err)
		}
		tmp := pColour.ToRGBA()
		// ignore alpha channel
		c = color.RGBA{R: tmp.R, G: tmp.G, B: tmp.B}

	} else if len(value) > 1 && value[0] == '#' {
		pColour, err := playcolors.ParseHEX(value)
		if err != nil {
			return fmt.Errorf("hex value %s invalid: %s", value, err)
		}
		tmp := pColour.ToRGBA()
		// ignore alpha channel
		c = color.RGBA{R: tmp.R, G: tmp.G, B: tmp.B}

	} else {
		var ok bool
		c, ok = colornames.Map[strings.ToLower(value)]
		if !ok {
			return fmt.Errorf("color name %s not in image/color names", value)
		}
	}

	co := LocalPenColour{
		Name:   value,
		Colour: c,
	}
	*l = co
	return nil
}

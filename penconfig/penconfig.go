// Package penconfig allows custom specification of pen strokes from a
// yaml format file on a per-layer basis, for example
//
//   all:
//     - pen:     fineliner
//       weight:  narrow
//       width:   0.95
//       color:   black
//       opacity: 0.9
//
//   "1":
//     - pen:     fineliner
//       weight:  narrow
//       width:   0.8
//       color:   blue
//       opacity: 0.8
//
package penconfig

import (
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// penTypes are the currently understood pen types
var penTypes = []string{
	"pen",
	"fineliner",
	"marker",
	"highligher",
	"eraser",
	"sharp pencil",
	"erase area",
	"paint",
	"mechanical pencil",
	"pencil",
	"ballpoint",
	"marker",
	"pen",
	"highlighter",
}

// penWeights are the currently understood pen weights
var penWeights = []string{"narrow", "standard", "broad"}

// PenConfig allows the configuration of a s
type PenConfig struct {
	Pen     string         `yaml:"pen"`
	Weight  string         `yaml:"weight"`
	Width   float32        `yaml:"width"`
	Colour  LocalPenColour `yaml:"color"`
	Opacity float64        `yaml:"opacity"`
}

// LayerPenConfigs defines StrokeSettings by layer
type LayerPenConfigs map[string][]PenConfig

// UnmarshalYAML is a custom unmarshaller
func (pc *PenConfig) UnmarshalYAML(value *yaml.Node) (err error) {

	// auxilliary unmarshal struct
	type AuxPenConfig struct {
		Pen     string  `yaml:"pen"`
		Weight  string  `yaml:"weight"`
		Width   float32 `yaml:"width"`
		Colour  string  `yaml:"color"`
		Opacity float64 `yaml:"opacity"`
	}

	var apc AuxPenConfig
	err = value.Decode(&apc)
	if err != nil {
		return fmt.Errorf("Yaml parsing error: %v", err)
	}

	lpc := LocalPenColour{}
	err = lpc.colourConvert(apc.Colour)
	if err != nil {
		return fmt.Errorf("colour convert error: %w", err)
	}

	*pc = PenConfig{
		Pen:     apc.Pen,
		Weight:  apc.Weight,
		Width:   apc.Width,
		Colour:  lpc,
		Opacity: apc.Opacity,
	}

	return nil
}

// LoadYaml reads bytes into a PenConfig structure
func LoadYaml(yamlByte []byte) (LayerPenConfigs, error) {

	var lpc LayerPenConfigs
	err := yaml.Unmarshal(yamlByte, &lpc)
	if err != nil {
		return lpc, err
	}
	err = lpc.check()
	return lpc, err
}

// check checks the validity of the configuration file
func (lpc LayerPenConfigs) check() error {

	// layers to list of pens
	for layer, penList := range lpc {
		if layer != "all" {
			_, err := strconv.Atoi(layer)
			if err != nil {
				return fmt.Errorf("penconfig layer %s needs to be 'all' or a layer number", layer)
			}
		}

		// list of pen configs in layer
		for i, pen := range penList {

			// check pen type
			penOK := false
			for _, penType := range penTypes {
				if pen.Pen == penType {
					penOK = true
					break
				}
			}
			if !penOK {
				return fmt.Errorf(
					"layer %s, item %d pen type %s not in\n%s",
					layer, i, pen.Pen, strings.Join(penTypes, " "),
				)
			}

			// check pen weight (should be checked by pen type too)
			weightOK := false
			for _, weight := range penWeights {
				if pen.Weight == weight {
					weightOK = true
					break
				}
			}
			if !weightOK {
				fmt.Printf("error pen %+v\n", pen)
				return fmt.Errorf(
					"layer %s, item %d weight type %s not in \n%s",
					layer, i, pen.Weight, strings.Join(penWeights, " "),
				)
			}

			// check pen opacity
			if pen.Opacity < 0.0 || pen.Opacity > 1.0 {
				return fmt.Errorf("layer %s, item %d opacity %f invalid", layer, i, pen.Opacity)
			}

			// check pen width
			if pen.Width < 0.0 || pen.Width > 30.0 {
				return fmt.Errorf("layer %s, item %d width %f invalid", layer, i, pen.Width)
			}

		}
	}
	return nil
}

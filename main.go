/*
rm2pdf entry point

MIT licensed, please see LICENCE
RCL 28 March 2021
*/

package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	rmpdf "github.com/rorycl/rm2pdf/rmpdf"
)

const version string = "0.1.6"

const usage = `InputPath OutputFile

rm2pdf version %s

Render layered PDF files from reMarkable tablet file bundles with
customisable pen widths and colours.

rm2pdf does not currently support remarkable software version 3 .rm
'lines' files.

Note that PDF files from sources such as Microsoft Word do not always
work well. It can help to rewrite them using the pdftk tool, e.g. by
doing

	pdftk word.pdf cat output word.pdf.pdftk \
	&& mv word.pdf word.pdf.bkp \
	&& mv word.pdf.pdftk word.pdf

For notebooks without a backing pdf file a template can be specified, of
which only the first page is used. If no template is provided the
embedded A4 template is used.

rm2pdf [-v] [-s pens.yaml] [-t A4red.pdf] [-c red]  `

// Options are flag options
type Options struct {
	Verbose  bool                `short:"v" long:"verbose"  description:"show verbose output\nthis presently does not do much"`
	Settings string              `short:"s" long:"settings" description:"path to customised pen settings file\nsee config_example.yaml for an example"`
	Template string              `short:"t" long:"template" description:"path to a single page A4 template to use when no UUID.pdf exists\nuseful for processing sketches without a backing PDF"`
	Colours  []rmpdf.LocalColour `short:"c" long:"colours"  description:"colour by layer\nuse several -c flags in series to select different colours\ne.g. -c red -c blue -c green for layers 1, 2 and 3.\nSee golang.org/x/image/colornames for the colours that can be used"`
	Args     struct {
		InputPath  string `description:"input path and uuid, optionally ending in '.pdf'"`
		OutputFile string `description:"output pdf file to write to"`
	} `positional-args:"yes" required:"yes"`
}

func main() {

	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = fmt.Sprintf(usage, version)

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	err := rmpdf.RM2PDF(options.Args.InputPath, options.Args.OutputFile, options.Template, options.Settings, options.Verbose, options.Colours)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

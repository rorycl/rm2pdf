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

const usage = `InputPath OutputFile

This programme attempts to create annotated PDF files from reMarkable
tablet file groups (RM bundles), including .rm files recording marks.

Normally these files will be in a local directory, such as an xochitl
directory synchronised to a reMarkable tablet over sshfs.

The programme takes as input either:
* The path to the PDF file which has had annotations made to it
* The path to the RM bundle with uuid, such as <path>/<uuid> with no
  filename extension, together with a PDF template to use for the
  background (a blank A4 template is provided in templates/A4.pdf).

The resulting PDF is layered with the background and .rm file layers
each in a separated PDF layer. The .rm file marks are stroked using the
fpdf PDF library, although .rm tilt and pressure characteristics are not
represented in the PDF output.

PDF files from sources such as Microsoft Word do not always work
well. It can help to rewrite them using the pdftk tool, e.g. by doing

	pdftk word.pdf cat output word.pdf.pdftk \
	&& mv word.pdf word.pdf.bkp \
	&& mv word.pdf.pdftk word.pdf

Custom colours for some pens can be specified using the -c or --colours
switch, which overrides the default pen selection. A second -c switch
sets the colours on the second layer, and so on.

Note that this programme has only been tested on a reMarkable v1 tablet.

rm2pdf -t templates/A4.pdf \
       testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 /tmp/output.pdf

rm2pdf [-v] [-c red] [-c green] [-c ...] `

// Options are flag options
type Options struct {
	Verbose  bool                `short:"v" long:"verbose"  description:"show verbose output\nthis presently does not do much"`
	Template string              `short:"t" long:"template" description:"path to a single page A4 template to use when no UUID.pdf exists\nuseful for processing sketches without a backing PDF"`
	Colours  []rmpdf.LocalColour `short:"c" long:"colours"  description:"colour by layer\nuse several -c flags in series to select different colours\ne.g. -c red -c blue -c green for layers 1, 2 and 3.\nSee golang.org/x/image/colornames for the colours that can be used"`
	Args     struct {
		InputPath  string `description:"input path and uuid, optionally ending in '.pdf'"`
		OutputFile string `description:"output pdf file to write to"`
	} `positional-args:"yes" required:"yes"`
}

// See pdf.rm2pdf for further details
func main() {

	var options Options
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = usage

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	err := rmpdf.RM2PDF(options.Args.InputPath, options.Args.OutputFile, options.Template, options.Verbose, options.Colours)
	if err != nil {
		fmt.Printf("An error occurred: %v", err)
	}
}

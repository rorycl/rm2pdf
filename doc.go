/*
rm2pdf

MIT Licensed
RCL January 2020


Introduction

This programme attempts to create annotated A4 PDF files from reMarkable
tablet file groups (RM bundles), including .rm files recording marks.

Normally these files will be in a local directory, such as an xochitl
directory synchronised to a tablet over sshfs.

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

Example of processing an rm bundle without a pdf:
	rm2pdf -t templates/A4.pdf \
	testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 \
	/tmp/output.pdf

Example of processing an rm bundle with a pdf, and per-layer colours:
	rm2pdf -c chartreuse -c firebrick \
	testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf \
	/tmp/output.pdf

General options:

	rm2pdf [-v] [-c red] [-c green] [-c ...] [-t template] InputPath OutputFile

Warning: the OutputFile will be overwritten if it exists.


ReMarkable .rm file parser

The parser is a go port of reMarkable tablet "lines" or ".rm" file
parser, with binary decoding hints drawn from rm2svg
https://github.com/reHackable/maxio/blob/master/tools/rM2svg which in
turn refers to https://github.com/lschwetlick/maxio/tree/master/tools.

Python struct format codes referred to in the parser, such as "<{}sI"
are from rm2svg.

RMParser provides a python-like iterator based on bufio.Scan, which
iterates over the referenced reMarkable .rm file returning a data
structure consisting of each path with its associated layer and path
segments.

Usage example:

	rm, err := rmparse.RMParse("filename.rm")
	// start parsing; note that pdflayers are dealt with sequentially
	for rm.Parse() {
		path := rm.Path.Path
		penName := StrokeMap[int(path.Pen)]
		for s := 1; s <= int(path.NumSegments); s++ {
			segment := rm.Path.Segments[s-1]
			// do something with path and/or segment
		}
	}


PDF paths, strokes and colours

Pen selections are hard-coded in stroke.go with widths, opacities and
colours. The StrokeSetting interface "Width" is used to scale strokes
based on nothing more than what seems to be about right.

Resolving the page sizes and reMarkable output resolution was based on
the reMarkable png templates and viewing the reMarkable app's output x
and y widths. These dimensions are noted in pdf.go in PDF_WIDTH_IN_MM
and PDF_HEIGHT_IN_MM. Conversion from mm to points (MM_TO_RMPOINTS) and
from points to the resolution of the reMarkable tablet (PTS_2_RMPTS) is
also set in pdf.go. The theoretical conversion factor is slightly
altered based on the output from various tests, including those in the
testfiles directory.

To view the testfiles after processing use or alter the paths used in
the tests.
*/
package main

/*
Generate a pdf with each page from an original pdf overlaid with marks
from a reMarkable tablet .rm file.

MIT licensed, please see LICENCE
RCL January 2020
*/

package rmpdf

import (
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/gofpdi"
	"github.com/rorycl/rm2pdf/files"
	"github.com/rorycl/rm2pdf/penconfig"
	"github.com/rorycl/rm2pdf/rmparse"
)

// reMarkable png templates (in /usr/share/remarkable/templates) are
// 1404x1873px or 495.4x660.9mm at 2.834x2.834 pixels/mm reMarkable
// output PDF sizes with fixed y and variable x use reMarkable png
// templates as a model

// PDFWidthInMM is the width of a PDF in millimetres
const PDFWidthInMM = 222.6264

// PDFHeightInMM is the width of a PDF in millimetres
const PDFHeightInMM = 297.0000

// MMtoRMPoints is the conversion between millimeters and standard
// postscript points
const MMtoRMPoints = 2.83465

// Pts2RMPoints is the conversion from rm pixels to points, theoretically
// 2.2253
const Pts2RMPoints = 2.222 // eyeballed conversion

// LayerRegister is a Layer names
var LayerRegister = map[string]int{}

// UnknownPens is an unknown pen register
var UnknownPens = make(map[int]int)

// List of colours to use, by layer, for pens with ColourOverride ==
// true
var layerColours = []LocalColour{}

// penConfigs are the set of custom pen settings by layer
var penConfigs = make(penconfig.LayerPenConfigs)

// Extract layer id from register by name, first initialising the PDF
// layerid for that name if necessary
func layerIDFromRegister(name string, pdf *gofpdf.Fpdf) int {
	if _, ok := LayerRegister[name]; !ok {
		LayerRegister[name] = pdf.AddLayer(name, true) // true == visible
	}
	return LayerRegister[name]
}

// Construct a pdf page with layers from rm files described by rmf
// RMFileInfo at 0-indexed page number, to be added to pdf. The existing
// pdf (annotated pdf, template pdf, or embedded template), described by
// sourceFH, is put in a "Background" layer, and the other .rm file
// layers are put into subsequent layers with a default PDF visibility
// of "true".
//
// Note that eraser types are presently skipped.
func constructPageWithLayers(rmf files.RMFileInfo, rmPageNo, pdfPageNo int, useTemplate bool, sourceFH *io.ReadSeeker, pdf *gofpdf.Fpdf) error {

	// add a new page
	pdf.AddPage()

	// add the base PDF within a PDF layer named "Background". Only A4
	// import files are presently supported, but it should be possible
	// to use base PDFs of other sizes although I can't find a
	// convenient way of determining the size of an imported page in go.
	layerID := layerIDFromRegister("Background", pdf)
	pdf.BeginLayer(layerID)

	rmf.Debug(fmt.Sprintf("%s rm page %d pdf page %d", rmf.IdentifyPDF(useTemplate), rmPageNo+1, pdfPageNo+1))

	// if an annotated pdf is provided, use the next page from that
	// if using the A4 template, recycle page use, based on output from
	// rmf.PageIterate from caller, whose pagenumbers are 0-indexed
	pdfImportPage := pdfPageNo + 1

	bgpdf := gofpdi.ImportPageFromStream(pdf, sourceFH, pdfImportPage, "/MediaBox")
	rmf.Debug(fmt.Sprintf("orientation %s", rmf.Orientation))
	if rmf.Orientation == "portrait" {
		gofpdi.UseImportedTemplate(pdf, bgpdf, 0, 0, 210*MMtoRMPoints, 297*MMtoRMPoints)
	} else {
		gofpdi.UseImportedTemplate(pdf, bgpdf, 0, 0, 297*MMtoRMPoints, 210*MMtoRMPoints)
	}
	pdf.EndLayer()

	// Initialise the .rm file parser if the .rm file exists, else return
	rmPage := rmf.Pages[rmPageNo]
	if !rmPage.Exists {
		rmf.Debug(fmt.Sprintf("no rm file for page %d ...skipping", rmPageNo+1))
		return nil
	}
	rmf.Debug(fmt.Sprintf("rmfile %s", rmPage.RMFilePath()))
	rm, err := rmparse.RMParse(rmPage.RMFile())
	if err != nil {
		return err
	}

	// set custom colours for layers, if provided
	pageLayerColours := map[int]LocalColour{}
	for c := 0; c < len(rmPage.LayerNames); c++ {
		if c <= len(layerColours)-1 {
			pageLayerColours[c] = layerColours[c]
		}
	}

	// layer setup
	// note that layers recorded in RMParse are 1-indexed, while the
	// LayerNames are 0 indexed
	layerNo := 1
	layerName := rmPage.LayerNames[layerNo-1]
	layerID = layerIDFromRegister(layerName, pdf)
	pdf.BeginLayer(layerID)
	rmf.Debug(fmt.Sprintf("Beginning layer %d", layerNo))

	// start parsing; note that pdflayers are dealt with sequentially
	// rm.Parse works on a per-path basis, implicity therefore on a
	// per-pen basis
	pathNum := 0
	for rm.Parse() {

		// start a new PDF layer if necessary
		if rm.Path.Layer != uint32(layerNo) {
			pdf.EndLayer()
			layerNo++
			rmf.Debug(fmt.Sprintf("Beginning layer %d", layerNo))
			layerName := rmPage.LayerNames[layerNo-1]
			layerID = layerIDFromRegister(layerName, pdf)
			pdf.BeginLayer(layerID)
		}

		path := rm.Path.Path

		// Skip eraser types
		penName := StrokeMap[int(path.Pen)]
		if penName == "eraser" || penName == "erase area" {
			continue
		}

		// set stroke colour, transparent fill color and line width
		// if opacity is not 1.0, set the alpha blending channel to the
		// required fraction of 1.0
		// Also record if an pen type is not found.
		penName, ok := StrokeMap[int(path.Pen)]
		if !ok {
			UnknownPens[int(path.Pen)]++
			penName = "fineliner"
		}
		ss := StrokeSettings[penName]

		width := ss.Width(path.Width)
		opacity := ss.Opacity // inclusive range [0,1]

		// load custom pen settings if any exist
		penWidthName := ss.NaturalWidth(path.Width)
		customPen, ok := penConfigs.GetPen(layerNo-1, penName, penWidthName)
		if ok {
			rmf.Debug(fmt.Sprintf("  path %4d : using custom pen %+v", pathNum, customPen))
			width = customPen.Width
			opacity = customPen.Opacity
		}

		// set colours, first checking to see if there is a custom pen
		// defined in the configuration file, then setting a colour
		// override if set
		//
		// pdf.SetFillSpotColor("White", 100) // 0% tint
		ok = false
		var layerCustomColour LocalColour

		layerCustomColour, ok = pageLayerColours[layerNo-1]
		if ok {
			rmf.Debug(fmt.Sprintf("  path %4d : using general layer colour %s", pathNum, layerCustomColour.Name))
			pdf.SetDrawColor(ss.selectColour(&layerCustomColour, false))
		} else if customPen != nil {
			// force
			pdf.SetDrawColor(ss.selectColour(
				&LocalColour{customPen.Colour.Name, customPen.Colour.Colour},
				true,
			))
		}

		// set width
		pdf.SetLineWidth(width)

		// set opacity
		if opacity != 1.0 {
			pdf.SetAlpha(opacity, "Normal")
		}

		// rmf.Debug(fmt.Sprintf("Pen : %s Width : %f, calcwidth %f, opacity %f", penName, path.Width, ss.Width(path.Width), ss.Opacity))

		for s := 1; s <= int(path.NumSegments); s++ {
			segment := rm.Path.Segments[s-1]

			// write rm segment to pdf path
			if rmf.Orientation == "portrait" {
				// portrait
				if s == 1 {
					pdf.MoveTo(float64(segment.X/Pts2RMPoints),
						float64(segment.Y/Pts2RMPoints))
				} else {
					pdf.LineTo(float64(segment.X/Pts2RMPoints),
						float64(segment.Y/Pts2RMPoints))
				}
			} else {
				// landscape format files need to be flipped
				yBasis := (297 * MMtoRMPoints)
				if s == 1 {
					pdf.MoveTo(yBasis-float64(segment.Y/Pts2RMPoints),
						float64(segment.X/Pts2RMPoints))
				} else {
					pdf.LineTo(yBasis-float64(segment.Y/Pts2RMPoints),
						float64(segment.X/Pts2RMPoints))
				}
			}
		}
		pdf.DrawPath("D") // outlined only; use FD for filled and outlined

		// reset opacity
		if opacity != 1.0 {
			pdf.SetAlpha(1.0, "Normal")
		}

		pathNum++
	}

	// close the layer
	pdf.EndLayer()

	rmf.Debug(fmt.Sprintf("Maximum coordinates : %+v\n", rm.MaxCoordinates))

	return nil
}

// RM2PDF is the main entry point for the programme. It takes a single
// string pointing to a valid PDF file (or the replacement A4 template)
// with an associated set of reMarkable metadata and .rm files. It then
// makes a PDF page for each page in the original PDF (although the
// template file's first page is recycled) and then adds each layer of
// the associated page's .rm file on top of that, finally writing the
// resulting pdf to outfile. Custom colours may be specified for each
// layer. Settings may also be supplied from a settings configuration
// file.
func RM2PDF(inputpath, outfile, template, settings string, verbose bool, colours []LocalColour) error {

	// initialise struct containing information about the files
	rmfile, err := files.RMFiler(inputpath, template)
	if err != nil {
		return err
	}

	if verbose {
		rmfile.Debugging = true
	}

	if (rmfile.OriginalPageCount != rmfile.OriginalPageCount) && template == "" {
		return fmt.Errorf(
			"bundle has inserted page/s %s and no template was provided",
			rmfile.InsertedPages(),
		)
	}

	// See fpdf PageSize example
	var pdf *gofpdf.Fpdf
	if rmfile.Orientation == "portrait" {
		pdf = gofpdf.NewCustom(&gofpdf.InitType{
			UnitStr: "pt",
			Size: gofpdf.SizeType{
				Ht: PDFHeightInMM * MMtoRMPoints,
				Wd: PDFWidthInMM * MMtoRMPoints,
			},
		})
	} else {
		pdf = gofpdf.NewCustom(&gofpdf.InitType{
			UnitStr: "pt",
			Size: gofpdf.SizeType{
				Wd: PDFHeightInMM * MMtoRMPoints,
				Ht: PDFWidthInMM * MMtoRMPoints,
			},
		})
	}

	// set custom layer colours if provided
	if len(colours) > 0 {
		layerColours = colours
	}

	// pen configuration file
	if settings != "" {
		var err error
		penConfigs, err = penconfig.NewPenConfigFromFile(settings)
		if err != nil {
			return fmt.Errorf("settings file load error: %w", err)
		}
	}

	// Make colour (White in CMYK notation) for transparent fill
	// pdf.AddSpotColor("White", 0, 0, 0, 0)

	// Add general line styles
	pdf.SetLineCapStyle("round")
	pdf.SetLineJoinStyle("round")

	// Iterate over pages using the rmfile iterator which provides a
	// page number and the pdf to use (either the annotated pdf or the
	// template). For annotated pdfs with inserted pages one might
	// receive the following output from the iterator:
	// pageno | inserted | template      | templatepageno
	// -------+----------+---------------+---------------
	// 0      | no       | annotated.pdf | 0
	// 1      | yes      | template.pdf  | 0
	// 2      | no       | annotated.pdf | 1

	// Iterate over each page in the pdf
	for i := 0; i < len(rmfile.Pages); i++ {
		pageNo, pdfPageNo, inserted, isTemplate, pdfFH := rmfile.PageIterate()
		rmfile.Debug(fmt.Sprintf(
			"processing page %d %d inserted %t template %t",
			pageNo, pdfPageNo, inserted, isTemplate,
		))
		constructPageWithLayers(rmfile, pageNo, pdfPageNo, isTemplate, pdfFH, pdf)
	}

	err = pdf.OutputFileAndClose(outfile)
	if err != nil {
		return err
	}

	if len(UnknownPens) > 0 {
		fmt.Println("Some pen types were not found, and were forced to the fineliner style")
		for k, v := range UnknownPens {
			fmt.Printf("pen: %02d occurrences: %d\n", k, v)
		}
	}

	return nil
}

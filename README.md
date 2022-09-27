# rm2pdf

version 0.1.5 : 27 September 2022

Convert reMarkable tablet file 'bundles' to layered PDFs, with optional
per-layer colours for selected pens.

## Update

Support older rmapi zip files and file bundles

Fix for [issue 9](https://github.com/rorycl/rm2pdf/issues/9), to skip
missing bundle metadata files for older reMarkable bundles, and support
0-indexed (rather than uuid-indexed) .rm files for older rmapi zip
files. Thanks to [qwert2003](https://github.com/qwert2003) for the
report and invaluable help.

`rm2pdf path_to.zip output.pdf`

Note that go 1.16+ is needed for `rm2pdf` due to the use of embedded
files, added in v0.1.3.

Recent releases:
* 0.1.4 : add support for rmapi zip files
* 0.1.3 : add embedded A4 template
* 0.1.2 : support custom pen configuration, see `config_example.yaml`.
* 0.1.1 : allow input paths with suffixes (such as `.content`).
* 0.1.0 : added support for landscape mode files.
* 0.0.3 : support for pages inserted while annotating a PDF.

## Examples

```
./rm2pdf -h

Usage:
  rm2pdf InputPath OutputFile

rm2pdf version 0.1.4

...

rm2pdf [-v] [-c red] [-c green] [-c ...]  InputPath OutputFile

Application Options:
  -v, --verbose     show verbose output
                    this presently does not do much
  -s, --settings=   path to customised pen settings file
  -t, --template=   path to a single page A4 template to use when no UUID.pdf exists
                    useful for processing sketches without a backing PDF
  -c, --colours=    colour by layer
                    use several -c flags in series to select different colours
                    e.g. -c red -c blue -c green for layers 1, 2 and 3.
                    See golang.org/x/image/colornames for the colours that can be used

Help Options:
  -h, --help        Show this help message

Arguments:
  InputPath:        input path and uuid, optionally ending in '.pdf'
  OutputFile:       output pdf file to write to

```

Invocation examples for annotated PDFs using the test files in `testfiles`:

```
rm2pdf testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf output.pdf

rm2pdf -c orange -c olivegreen \
       testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf output2.pdf
```

Invocation examples for reMarkable notebooks using the test files in `testfiles`
and the A4 template in `templates`.

```
rm2pdf -c blue -c red -t templates/A4.pdf \
       testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 output4.pdf

```

The embedded template is used where one is not provided, so the above command is
the same as

```
rm2pdf -c blue -c red \
       testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 output4.pdf
```

## Details

rm2pdf requires "bundles" of files created on the reMarkable tablet, including
the associated `.metadata` and `.content` files, together with the `.rm` binary
and `-metadata.json` files associated with each page of reMarkable marks.

rm2pdf aims to create PDFs from both PDFs that are annotated on the reMarkable
and reMarkable notebooks using these files. The latter uses an empty template
PDF as the background. PDF templates can be made from the reMarkable png
templates (in /usr/share/remarkable/templates) but should fit the standard
222.6264mm x 297.0000mm reMarkable output PDF size, or be A4.

Output PDFs are layered with the background PDF forming a "Background" layer and
subsequent layers using the layer names created on the tablet. The layers can be
turned on and off using tools provided by PDF readers such as Evince.

The pen widths and opacities provided by default are estimates. Colours, base
width and opacity are set for each pen are set in rmpdf/stroke.go. Those pens
with ColourOverride true will have their colour overridden by the command-line
options or pen configuration yaml file. Note that at present the .rm file
pressure and tilt information are not presently used. 

Some PDF files, notably those created by Microsoft Word, cannot be imported
reliably, causing the programme to panic. Reprocessing problem PDFs with the
`pdftk` tool seems to fix the problem.

Note that rm2pdf has only been tested on a reMarkable v1 tablet.

## Background

The project includes rmparse/rmparse.go, a remarkable tablet Go port of
reMarkable tablet "lines" or ".rm" file parser, with binary decoding hints drawn
from rm2svg https://github.com/reHackable/maxio/blob/master/tools/rM2svg which
in turn refers to https://github.com/lschwetlick/maxio/tree/master/tools.

The project makes extensive use of the go PDF `fpdf` library and the contrib
module `gofpdi`. The latter is used for including pages from existing PDF
documents.

If your pdf causes fpdf to fail, resave the pdf using the `pdftk`
programme.

### Build and test

Developed with go 1.18 on 64bit Linux.

Test:  `go test -v ./...`

Build : `go build`; this should produce an executable called `rm2pdf`.

## License

This project is licensed under the [MIT Licence](LICENCE).

# rm2pdf

version 0.0.1-alpha

Convert reMarkable tablet file 'bundles' to layered PDFs, with optional
per-layer colours for selected pens.

```
rm2pdf -h
```

Invocation examples for annotated PDFs using the test files in `testfiles`:

```
rm2pdf testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf /tmp/output.pdf
rm2pdf -c orange -c olivegreen testfiles/cc8313bb-5fab-4ab5-af39-46e6d4160df3.pdf /tmp/output2.pdf
```

Invocation examples for reMarkable notebooks using the test files in `testfiles`
and the A4 template in `templates`.

```
rm2pdf -t templates/A4.pdf testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 /tmp/output3.pdf
rm2pdf -c blue -c red -t templates/A4.pdf testfiles/d34df12d-e72b-4939-a791-5b34b3a810e7 /tmp/output4.pdf
```

# Details

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

The pen widths and opacities are estimates and could be improved. The .rm file
pressure and tilt information is not presently used. 

Colours, base width and opacity are set for each pen are set in rmpdf/stroke.go.
Those pens with ColourOverride true will have their colour overridden by the
command-line options.

Some PDF files, notably those created by Microsoft Word, cannot be imported
reliably, causing the programme to panic. Reprocessing problem PDFs with the
`pdftk` tool seems to fix the problem.

## Background

This is a holiday project to learn go. Error handling in particular is pretty
poor.

The project includes rmparse/rmparse.go, a remarkable tablet Go port of
reMarkable tablet "lines" or ".rm" file parser, with binary decoding hints drawn
from rm2svg https://github.com/reHackable/maxio/blob/master/tools/rM2svg which
in turn refers to https://github.com/lschwetlick/maxio/tree/master/tools.

The project makes extensive use of the go PDF `fpdf` library and the contrib
module `gofpdi`. The latter is used for including pages from existing PDF
documents.

### Build and test

Developed with go 1.13 on 64bit Linux.

Test:  `go test -v ./...`

Build : `go build`; this should produce an executable called `rm2svg`.

## License

This project is licensed under the [MIT Licence](LICENCE).

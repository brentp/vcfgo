package vcfgo

import (
	"fmt"
	"io"
	"strings"
)

// Writer allows writing VCF files.
type Writer struct {
	io.Writer
	Header *Header
}

// NewWriter returns a writer after writing the header.
func NewWriter(w io.Writer, h *Header) (*Writer, error) {
	fmt.Fprintf(w, "##fileformat=VCFv%s\n", h.FileFormat)

	for k, v := range h.Extras {
		fmt.Fprintf(w, "## %s=%s\n", k, v)
	}

	for _, imap := range h.Contigs {
		fmt.Fprintf(w, "##contig=<ID=%s", imap["ID"])

		for k, v := range imap {
			if k == "ID" {
				continue
			}

			fmt.Fprintf(w, ",%s=%s", k, v)
		}
		fmt.Fprintln(w, ">")
	}

	for k, v := range h.Filters {
		fmt.Fprintf(w, "##FILTER=<ID=,Description=\"%s\">\n", k, v)
	}

	for _, i := range h.Infos {
		fmt.Fprintf(w, "%s\n", i)
	}

	for _, i := range h.SampleFormats {
		fmt.Fprintf(w, "%s\n", i)
	}

	fmt.Fprint(w, "#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT")
	var s string
	if len(h.SampleNames) > 0 {
		s = "\t" + strings.Join(h.SampleNames, "\t")
	}

	fmt.Fprint(w, s+"\n")
	return &Writer{w, h}, nil
}

// WriteVariant writes a single variant
func (w *Writer) WriteVariant(v *Variant) {
	fmt.Println(w, v)
}

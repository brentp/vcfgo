package vcfgo

import (
	"fmt"
	"io"
	"sort"
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

	// Samples
	keys := make([]string, 0, len(h.Samples))
	for sampleId := range h.Samples {
		keys = append(keys, sampleId)
	}
	sort.Strings(keys)
	for _, sampleId := range keys {
		fmt.Fprintln(w, h.Samples[sampleId])
	}

	for i := range h.Pedigrees {
		fmt.Fprintln(w, h.Pedigrees[i])
	}

	// Filters
	keys = keys[:0]
	for k := range h.Filters {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "##FILTER=<ID=%s,Description=\"%s\">\n", k, h.Filters[k])
	}

	// Infos
	keys = keys[:0]
	for k := range h.Infos {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s\n", h.Infos[k])
	}

	// SampleFormats
	keys = keys[:0]
	for k := range h.SampleFormats {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(w, "%s\n", h.SampleFormats[k])
	}
	for _, line := range h.Extras {
		fmt.Fprintf(w, "%s\n", line)
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
	fmt.Fprintln(w, v)
}

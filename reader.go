// Package vcfgo implements a Reader and Writer for variant call format.
// It eases reading, filtering modifying VCF's even if they are not to spec.
// Example:
//  f, _ := os.Open("examples/test.auto_dom.no_parents.vcf")
//  rdr, err := vcfgo.NewReader(f)
//  if err != nil {
//  	panic(err)
//  }
//  for {
//  	variant := rdr.Read()
//  	if variant == nil {
//  		break
//  	}
//  	fmt.Printf("%s\t%d\t%s\t%s\n", variant.Chromosome, variant.Pos, variant.Ref, variant.Alt)
//  	fmt.Printf("%s", variant.Info["DP"].(int) > 10)
//  	sample := variant.Samples[0]
//  	// we can get the PL field as a list (-1 is default in case of missing value)
//  	fmt.Println("%s", variant.GetGenotypeField(sample, "PL", -1))
//  	_ = sample.DP
//  }
//  fmt.Fprintln(os.Stderr, rdr.Error())
package vcfgo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	"github.com/brentp/irelate/interfaces"
)

// used for the quality score which is 0 to 255, but allows "."
const MISSING_VAL = 256

// Reader holds information about the current line number (for errors) and
// The VCF header that indicates the structure of records.
type Reader struct {
	buf         *bufio.Reader
	Header      *Header
	verr        *VCFError
	LineNumber  int64
	lazySamples bool
	r           io.Reader
}

func NewWithHeader(r io.Reader, h *Header, lazySamples bool) (*Reader, error) {
	buf := bufio.NewReaderSize(r, 32768/8)
	var verr = NewVCFError()
	return &Reader{buf, h, verr, 1, lazySamples, r}, nil
}

// NewReader returns a Reader.
// If lazySamples is true, then the user will have to call Reader.ParseSamples()
// in order to access simple info.
func NewReader(r io.Reader, lazySamples bool) (*Reader, error) {
	buffered := bufio.NewReaderSize(r, 32768/8)

	var verr = NewVCFError()

	var LineNumber int64
	h := NewHeader()

	for {

		LineNumber++
		line, err := buffered.ReadString('\n')
		if err != nil && err != io.EOF {
			verr.Add(err, LineNumber)
		}
		if len(line) > 1 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		if LineNumber == 1 {
			v, err := parseHeaderFileVersion(line)
			verr.Add(err, LineNumber)
			h.FileFormat = v

		} else if strings.HasPrefix(line, "##FORMAT") {
			format, err := parseHeaderFormat(line)
			verr.Add(err, LineNumber)
			if format != nil {
				h.SampleFormats[format.Id] = format
			}

		} else if strings.HasPrefix(line, "##INFO") {
			info, err := parseHeaderInfo(line)
			verr.Add(err, LineNumber)
			if info != nil {
				h.Infos[info.Id] = info
			}

		} else if strings.HasPrefix(line, "##FILTER") {
			filter, err := parseHeaderFilter(line)
			verr.Add(err, LineNumber)
			if filter != nil && len(filter) == 2 {
				h.Filters[filter[0]] = filter[1]
			}

		} else if strings.HasPrefix(line, "##contig") {
			contig, err := parseHeaderContig(line)
			verr.Add(err, LineNumber)
			if contig != nil {
				if _, ok := contig["ID"]; ok {
					h.Contigs = append(h.Contigs, contig)
				} else {
					verr.Add(fmt.Errorf("bad contig: %v", line), LineNumber)
				}
			}
		} else if strings.HasPrefix(line, "##SAMPLE") {
			sample, err := parseHeaderSample(line)
			verr.Add(err, LineNumber)
			if sample != "" {
				h.Samples[sample] = line
			} else {
				verr.Add(fmt.Errorf("bad sample: %v", line), LineNumber)
			}
		} else if strings.HasPrefix(line, "##PEDIGREE") {
			h.Pedigrees = append(h.Pedigrees, line)
		} else if strings.HasPrefix(line, "##") {
			kv, err := parseHeaderExtraKV(line)
			verr.Add(err, LineNumber)

			if kv != nil && len(kv) == 2 {
				h.Extras = append(h.Extras, line)
			}

		} else if strings.HasPrefix(line, "#CHROM") {
			var err error
			h.SampleNames, err = parseSampleLine(line)
			verr.Add(err, LineNumber)
			//h.Validate(verr)
			break

		} else {
			e := fmt.Errorf("unexpected header line: %s", line)
			return nil, e
		}
	}
	reader := &Reader{buffered, h, verr, LineNumber, lazySamples, r}
	return reader, reader.Error()
}

// Read returns a pointer to a Variant. Upon reading the caller is assumed
// to check Reader.Err()
func (vr *Reader) Read() interfaces.IVariant {

	line, err := vr.buf.ReadBytes('\n')
	if err != nil {
		if len(line) == 0 && err == io.EOF {
			return nil
		} else if err != io.EOF {
			vr.verr.Add(err, vr.LineNumber)
		}
	}

	vr.LineNumber++
	if line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	fields := bytes.SplitN(line, []byte{'\t'}, 10)

	return vr.Parse(fields)
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (vr *Reader) Parse(fields [][]byte) interfaces.IVariant {

	pos, err := strconv.ParseUint(unsafeString(fields[1]), 10, 64)
	vr.verr.Add(err, vr.LineNumber)

	var qual float64
	if len(fields[5]) == 1 && fields[5][0] == '.' {
		qual = MISSING_VAL
	} else {
		qual, err = strconv.ParseFloat(unsafeString(fields[5]), 32)
	}

	vr.verr.Add(err, vr.LineNumber)

	v := &Variant{Chromosome: string(fields[0]), Pos: pos, Id_: string(fields[2]), Reference: string(fields[3]), Alternate: strings.Split(string(fields[4]), ","), Quality: float32(qual),
		Filter: string(fields[6]), Header: vr.Header}

	if len(fields) > 8 {
		v.Format = strings.Split(string(fields[8]), ":")
		if len(fields) > 9 {
			v.sampleString = string(fields[9])
		}
		if !vr.lazySamples {
			vr.Header.ParseSamples(v)
		}
	}
	v.LineNumber = vr.LineNumber

	v.Info_ = NewInfoByte(fields[7], vr.Header)
	vr.verr.Add(err, vr.LineNumber)
	return v
}

// Force parsing of the sample fields.
func (h *Header) ParseSamples(v *Variant) error {
	if v.Format == nil || v.sampleString == "" || v.Samples != nil {
		return nil
	}
	var errors []error
	v.Samples = make([]*SampleGenotype, len(h.SampleNames))

	for i, sample := range strings.Split(v.sampleString, "\t") {
		var geno *SampleGenotype
		geno, errors = h.parseSample(v.Format, sample)

		v.Samples[i] = geno
	}
	v.sampleString = ""
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// Add a INFO field to the header.
func (vr *Reader) AddInfoToHeader(id string, num string, stype string, desc string) {
	h := vr.Header
	h.Infos[id] = &Info{Id: id, Number: num, Type: stype, Description: desc}
}

// Error() aggregates the multiple errors that can occur into a single object.
func (vr *Reader) Error() error {
	if vr.verr.IsEmpty() {
		return nil
	}
	return vr.verr
}

// Clear empties the cache of errors.
func (vr *Reader) Clear() {
	vr.verr.Clear()
}

func (vr *Reader) Close() error {
	if rc, ok := vr.r.(io.ReadCloser); ok {
		return rc.Close()
	}
	return nil
}

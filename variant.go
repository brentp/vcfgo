package vcfgo

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/brentp/irelate/interfaces"
)

// Variant holds the information about a single site. It is analagous to a row in a VCF file.
type Variant struct {
	Chromosome string
	Pos        uint64
	Id_        string
	Reference  string
	Alternate  []string
	Quality    float32
	Filter     string
	Info_      interfaces.Info
	Format     []string
	Samples    []*SampleGenotype
	// if lazy parsing, then just save the sample strings here.
	sampleString string
	Header       *Header
	LineNumber   int64
}

func (v *Variant) Info() interfaces.Info {
	return v.Info_
}

func (v *Variant) Id() string {
	return v.Id_
}

func (v *Variant) Ref() string {
	return v.Reference
}

func (v *Variant) Alt() []string {
	return v.Alternate
}

// Chrom returns the chromosome name.
func (v *Variant) Chrom() string {
	return v.Chromosome
}

// Start returns the 0-based start
func (v *Variant) Start() uint32 {
	return uint32(v.Pos - 1)
}

// CIPos reports the Left and Right end of an SV using the CIPOS tag. It is in
// bed format so the end is +1'ed. E.g. If there is not CIPOS, the return value
// is v.Start(), v.Start() + 1
func (v *Variant) CIPos() (uint32, uint32, bool) {
	if v.Header == nil {
		v.Header = NewHeader()
	}
	s := v.Start()
	v.Header.RLock()
	if _, ok := v.Header.Infos["CIPOS"]; !ok {
		v.Header.RUnlock()
		v.Header.Lock()
		v.Header.Infos["CIPOS"] = &Info{Id: "CIPOS", Number: "2", Description: "CIPOS", Type: "Integer"}
		v.Header.Unlock()
	} else {
		v.Header.RUnlock()
	}

	ipair, err := v.Info().Get("CIPOS")
	if ipair == nil && err != nil {
		return s, s + 1, false
	}
	pair, ok := ipair.([]int)
	if !ok {
		return s, s + 1, false
	}
	left := pair[0]
	right := pair[1]
	return uint32(int(s) + left), uint32(int(s) + right + 1), true
}

// CIEnd reports the Left and Right end of an SV using the CIEND tag. It is in
// bed format so the end is +1'ed. E.g. If there is no CIEND, the return value
// is v.End() - 1, v.End()
func (v *Variant) CIEnd() (uint32, uint32, bool) {
	if v.Header == nil {
		v.Header = NewHeader()
	}
	e := v.End()
	v.Header.RLock()
	if _, ok := v.Header.Infos["CIEND"]; !ok {
		v.Header.RUnlock()
		v.Header.Lock()
		v.Header.Infos["CIEND"] = &Info{Id: "CIEND", Number: "2", Description: "CIEND", Type: "Integer"}
		v.Header.Unlock()
	} else {
		v.Header.RUnlock()
	}
	ipair, err := v.Info().Get("CIEND")
	if ipair == nil && err != nil {
		return e - 1, e, false
	}
	pair, ok := ipair.([]int)
	if !ok {
		return e - 1, e, false
	}
	left := pair[0]
	right := pair[1]
	return uint32(int(e)+left) - 1, uint32(int(e) + right), true
}

// End returns the 0-based start + the length of the reference allele.
func (v *Variant) End() uint32 {
	/*
		if len(v.Alt()[0]) == 0 {
			log.Println(v.Chrom(), v.Start()+1, v.Ref(), v.Alt())
		}*/
	a := v.Alt()[0]
	if len(a) == 0 || a[0] != '<' {
		return uint32(v.Pos-1) + uint32(len(v.Ref()))
	}
	if strings.HasPrefix(v.Alt()[0], "<DEL") || strings.HasPrefix(v.Alt()[0], "<DUP") || strings.HasPrefix(v.Alt()[0], "<IN") || strings.HasPrefix(v.Alt()[0], "<CN") {
		if svlen, err := v.Info().Get("SVLEN"); err == nil || (strings.Contains(err.Error(), "not found in header") && svlen != nil) {
			var slen int
			err = nil
			switch svlen.(type) {
			case int:
				slen = svlen.(int)
			case string:
				var e error
				if svlen.(string) == "" {
					return uint32(v.Pos)
				}
				slen, e = strconv.Atoi(svlen.(string))
				if e != nil {
					log.Fatalf("bad value for svlen: %s\n", svlen)
				}
			case float64:
				slen = int(svlen.(float64))
			case float32:
				slen = int(svlen.(float32))
			case []interface{}:
				slen = svlen.([]interface{})[0].(int)
			case interface{}:
				slen = svlen.(int)
			default:
				log.Printf("non int type for SVLEN:%s at %s:%d\n", svlen, v.Chrom(), v.Pos)
				slen = 1
			}
			if slen < 0 {
				slen = -slen
			}
			return uint32(int(v.Pos) + slen)

		} else if end, err := v.Info().Get("END"); err == nil || end != nil {
			if end == nil || end == "" {
				if a != "<CN0>" {
					log.Printf("non int type for END and SVLEN:%s at %s:%d\n", svlen, v.Chrom(), v.Pos)
				}
				return uint32(v.Pos + 1)
			}

			if e, ok := end.(int); ok {
				return uint32(e)
			}
			if s, ok := end.(string); ok {
				if v, err := strconv.Atoi(s); err == nil {
					return uint32(v)
				} else {
					log.Printf("error parsing INFO/END: %s", err)
				}
			}
		}
		log.Printf("no svlen for variant %s:%d\n%s\nUsing %d", v.Chromosome, v.Pos, v, v.Pos+1)
	}
	// <INS and BND's get handled by this.
	return uint32(v.Pos-1) + uint32(len(v.Ref()))
}

func fmtFloat32(v float32) string {
	var val string
	if v > 0.02 || v < -0.02 {
		val = fmt.Sprintf("%.4f", v)
	} else {
		val = fmt.Sprintf("%.5g", v)
	}
	val = strings.TrimRight(strings.TrimRight(val, "0"), ".")
	if val == "" || val == "-" {
		val = "0"
	}
	return val
}

func fmtFloat64(v float64) string {
	var val string
	if v > 0.02 || v < -0.02 {
		val = fmt.Sprintf("%.4f", v)
	} else {
		val = fmt.Sprintf("%.5g", v)
	}
	val = strings.TrimRight(strings.TrimRight(val, "0"), ".")
	if val == "" || val == "-" {
		val = "0"
	}
	return val
}

// SampleGenotype holds the information about a sample. Several fields are pre-parsed, but
// all fields are kept in Fields as well.
type SampleGenotype struct {
	Phased bool
	GT     []int
	DP     int
	GL     []float64
	GQ     int
	MQ     int
	Fields map[string]string
}

// String returns the string representation of the sample field.
func (s *SampleGenotype) String() string {
	return fmt.Sprintf("GT: %v, DP: %d, GL: %v, GQ: %d, MQ: %d, Fields: %v", s.GT, s.DP, s.GL, s.GQ, s.MQ, s.Fields)
}

// RefDepth returns the depths of the alternates for this sample
func (s *SampleGenotype) RefDepth() (int, error) {
	if ad, ok := s.Fields["AD"]; ok {
		idx := strings.Index(ad, ",")
		return strconv.Atoi(ad[:idx])
	}
	if ro, ok := s.Fields["RO"]; ok {
		return strconv.Atoi(ro)
	}
	return 0, errors.New("only freebayes and gatk depths supported at this time")
}

// AltDepths returns the depths of the alternates for this sample
func (s *SampleGenotype) AltDepths() ([]int, error) {
	var svals []string
	if ad, ok := s.Fields["AD"]; ok {
		idx := strings.Index(ad, ",")
		svals = strings.Split(ad[idx+1:], ",")
	} else if ro, ok := s.Fields["AO"]; ok {
		svals = strings.Split(ro, ",")
	} else {
		return []int{}, errors.New("only freebayes and GATK supported for ref/alt depths")
	}
	vals := make([]int, len(svals))
	for i := range svals {
		v, err := strconv.Atoi(svals[i])
		if err != nil {
			return []int{}, err
		}
		vals[i] = v
	}
	return vals, nil
}

// ToString returns the string representation of the sample field.
func (sg *SampleGenotype) ToString(fields []string) string {
	if len(fields) == 0 {
		return "."
	}

	s := make([]string, len(fields))
	for i, f := range fields {
		s[i] = sg.Fields[f]
	}
	return strings.Join(s, ":")
}

// NewSampleGenotype allocates the internals and returns a *SampleGenotype
func NewSampleGenotype() *SampleGenotype {
	s := &SampleGenotype{}
	s.GT = make([]int, 0, 2)
	s.GL = make([]float64, 0, 3)
	s.Fields = make(map[string]string)
	return s
}

// String gives a string representation of a variant
func (v *Variant) String() string {
	var qual string
	
	if v.Quality < 0 {
		qual = "."
	} else {
		qual = fmt.Sprintf("%.1f", v.Quality)
	}

	s := fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s", v.Chromosome, v.Pos, v.Id_, v.Ref(), strings.Join(v.Alt(), ","), qual, v.Filter, v.Info())
	if len(v.Samples) > 0 {
		samps := make([]string, len(v.Samples))
		for i, s := range v.Samples {
			samps[i] = s.ToString(v.Format)
		}
		s += fmt.Sprintf("\t%s\t%s", strings.Join(v.Format, ":"), strings.Join(samps, "\t"))
	} else if v.sampleString != "" {
		s += fmt.Sprintf("\t%s\t%s", strings.Join(v.Format, ":"), v.sampleString)
	}
	return s
}

// GetGenotypeField uses the information from the header to parse the correct time from a genotype field.
// It returns an interface that can be asserted to the expected type.
func (v *Variant) GetGenotypeField(g *SampleGenotype, field string, missing interface{}) (interface{}, error) {
	if g == nil {
		return missing, fmt.Errorf("GetGenotypeField: empty genotype when requesting %s", field)
	}
	h := v.Header
	format, ok := h.SampleFormats[field]
	if !ok {
		return nil, fmt.Errorf("GetGenotypeField: field not found in formats: %s", field)
	}
	value, ok := g.Fields[field]
	if !ok {
		return nil, fmt.Errorf("GetGenotypeField: field not found in genotypes: %s", field)
	}
	switch format.Type {
	case "Integer":
		var mv int
		var ok bool
		if mv, ok = missing.(int); !ok {
			return nil, fmt.Errorf("GetGenotypeField: bad non-int missing value: %v", missing)
		}
		return handleNumberType(format.Number, value, len(v.Alt()), len(g.GT), true, mv)

	case "Float":
		var mv float32
		var ok bool
		if mv, ok = missing.(float32); !ok {
			return nil, fmt.Errorf("GetGenotypeField: bad non-float missing value: %v", missing)
		}
		return handleNumberType(format.Number, value, len(v.Alt()), len(g.GT), false, mv)

	case "String", "Character", "Unknown":
		return value, nil

	case "Flag":
		return field, nil

	}

	return nil, fmt.Errorf("unknown format: %s", format.Type)
}

func handleNumberType(number string, value string, nAlts int, nGTs int, isInt bool, mv interface{}) (interface{}, error) {
	if number == "1" || !strings.Contains(value, ",") || number == "." || number == "" {
		if isInt {
			if value == "" || value == "." {
				return (mv).(int), nil
			}
			return strconv.Atoi(value)
		}
		if value == "" || value == "." {
			return (mv).(float32), nil
		}
		return strconv.ParseFloat(value, 32)
	}
	if count, err := strconv.Atoi(number); err == nil || number == "G" || number == "A" || number == "R" {
		if err != nil {
			switch number {
			case "G":
				count = nGTs * (nGTs + 1) / 2
			case "A":
				count = nAlts
			case "R":
				count = nAlts + 1
			}
			err = nil
		}
		var ret interface{}
		split := strings.Split(value, ",")
		if isInt {
			ret = make([]int, len(split), len(split))
		} else {
			ret = make([]float32, len(split), len(split))
		}

		var countErr error

		// caller can ignore error if they want, we still fill what we can.
		if len(split) != count {
			countErr = fmt.Errorf("number of fields (%d) does not match expected (%d) in '%s'", len(split), count, value)
		}
		for i, s := range split {
			if isInt {
				ri, err := strconv.Atoi(s) //, 10, 32)
				if err != nil {
					// if it's an error, we allow empty
					if s == "" || s == "." {
						ret.([]int)[i] = mv.(int)
						err = nil
					} else {
						return nil, fmt.Errorf("non integer type: %s", s)
					}
				} else {
					ret.([]int)[i] = int(ri)
				}
			} else {
				rf, err := strconv.ParseFloat(s, 32)
				if err != nil {
					// if it's an error, we allow empty
					if s == "" || s == "." {
						ret.([]float32)[i] = mv.(float32)
						err = nil
					} else {
						return nil, fmt.Errorf("non float type: %s", s)
					}
				} else {
					ret.([]float32)[i] = float32(rf)
				}
			}
		}
		return ret, countErr
	} else if number == "." || number == "" {
		return value, nil
	} else {
		return nil, fmt.Errorf("unknown number field: %s", number)
	}
}

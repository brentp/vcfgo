package vcfgo

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// InfoMap holds the parsed Info field which can contain floats, ints and lists thereof.
type InfoMap map[string]interface{}

func (i InfoMap) Add(key string, o interface{}) {
	i[key] = o
	i["__order"] = append(i["__order"].([]string), key)
}

// Variant holds the information about a single site. It is analagous to a row in a VCF file.
type Variant struct {
	Chromosome string
	Pos        uint64
	Id         string
	Ref        string
	Alt        []string
	Quality    float32
	Filter     string
	Info       InfoMap
	Format     []string
	Samples    []*SampleGenotype
	// if lazy parsing, then just save the sample strings here.
	sampleStrings []string
	Header        *Header
	LineNumber    int64
}

// Is returns true if variants are the same by position and share at least 1 alternate allele.
func (v *Variant) Is(o *Variant) bool {
	if v.Pos != o.Pos || v.Chromosome != o.Chromosome || v.Ref != o.Ref {
		return false
	}
	for _, av := range v.Alt {
		for _, ov := range o.Alt {
			if av == ov {
				return true
			}
		}
	}
	return false
}

// Chrom returns the chromosome name.
func (v *Variant) Chrom() string {
	return v.Chromosome
}

// Start returns the 0-based start
func (v *Variant) Start() uint32 {
	return uint32(v.Pos - 1)
}

// End returns the 0-based start + the length of the reference allele.
func (v *Variant) End() uint32 {
	return uint32(v.Pos-1) + uint32(len(v.Ref))
}

func fmtFloat(v float32) string {
	var val string
	if v > 0.02 || v < -0.02 {
		val = fmt.Sprintf("%.2f", v)
	} else {
		val = fmt.Sprintf("%.5gf", v)
	}
	return val
}

func fmtFloat64(v float64) string {
	var val string
	if v > 0.02 || v < -0.02 {
		val = fmt.Sprintf("%.2f", v)
	} else {
		val = fmt.Sprintf("%.5gf", v)
	}
	return val
}

// String returns a string that matches the original info field.
func (m InfoMap) String() string {
	var order []string
	// use __order internally to keep order of keys.
	order, ok := m["__order"].([]string)
	if !ok {
		order = make([]string, 0)
		for k := range m {
			order = append(order, k)
		}
		sort.Strings(order)

	}
	s := ""
	for j, k := range order {
		v := m[k]
		if b, ok := v.(bool); ok && b {
			s += k
		} else {
			switch v.(type) {
			case float32:
				s += k + "=" + fmtFloat(v.(float32))
			case float64:
				s += k + "=" + fmtFloat64(v.(float64))
			case int:
				s += fmt.Sprintf("%s=%d", k, v.(int))
			case []interface{}:

				switch v.([]interface{})[0].(type) {
				case float64:
					for _, vv := range v.([]interface{}) {
						s += k + "=" + fmtFloat64(vv.(float64))
					}
				case int:
					for _, vv := range v.([]interface{}) {
						s += fmt.Sprintf("%s=%d", k, vv.(int))
					}
				}
			default:
				s += fmt.Sprintf("%s=%s", k, v.(string))
			}
		}
		if j < len(order)-1 {
			s += ";"
		}
	}
	return s
}

// SampleGenotype holds the information about a sample. Several fields are pre-parsed, but
// all fields are kept in Fields as well.
type SampleGenotype struct {
	Phased bool
	GT     []int
	DP     int
	GL     []float32
	GQ     int
	MQ     int
	// TODO: add methods for Ref, Alt depth.
	Fields map[string]string
}

// String returns the string representation of the sample field.
func (sg *SampleGenotype) String(fields []string) string {
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
	s.GL = make([]float32, 0, 3)
	s.Fields = make(map[string]string)
	return s
}

// String gives a string representation of a variant
func (v *Variant) String() string {
	//#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	1_dad	1_mom	1_kid	2_dad	2_mom	2_kid	3_dad	3_mom	3_kid
	s := fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%.1f\t%s\t%s\t", v.Chromosome, v.Pos, v.Id, v.Ref, strings.Join(v.Alt, ","), v.Quality, v.Filter, v.Info)
	if len(v.Samples) > 0 {
		samps := make([]string, len(v.Samples))
		for i, s := range v.Samples {
			samps[i] = s.String(v.Format)
		}
		s += fmt.Sprintf("%s\t%s", strings.Join(v.Format, ":"), strings.Join(samps, "\t"))
	} else if v.sampleStrings != nil && len(v.sampleStrings) != 0 {
		s += fmt.Sprintf("%s\t%s", strings.Join(v.Format, ":"), strings.Join(v.sampleStrings, "\t"))
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
		return handleNumberType(format.Number, value, len(v.Alt), len(g.GT), true, mv)

	case "Float":
		var mv float32
		var ok bool
		if mv, ok = missing.(float32); !ok {
			return nil, fmt.Errorf("GetGenotypeField: bad non-float missing value: %v", missing)
		}
		return handleNumberType(format.Number, value, len(v.Alt), len(g.GT), false, mv)

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

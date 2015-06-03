package vcfgo

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// InfoMap holds the parsed Info field which can contain floats, ints and lists thereof.
type InfoMap map[string]interface{}

func (i InfoMap) Add(key string, o interface{}) {
	i[key] = o
	var order []string
	if ord, ok := i["__order"]; !ok {
		order = make([]string, 0)
		for k := range i {
			order = append(order, k)
		}
	} else {
		order = ord.([]string)
	}
	for i := len(order) - 1; i > -1; i-- {
		if key == order[i] {
			return
		}
	}
	i["__order"] = append(order, key)
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
				s += k + "=" + fmtFloat32(v.(float32))
			case float64:
				s += (k + "=" + fmtFloat64(v.(float64)))
			case int:
				s += fmt.Sprintf("%s=%d", k, v.(int))
			case uint32:
				s += fmt.Sprintf("%s=%d", k, v.(uint32))
			case []interface{}:
				vals := v.([]interface{})
				svals := make([]string, len(vals))
				switch vals[0].(type) {
				case float64:
					for i, val := range vals {
						svals[i] = fmtFloat64(val.(float64))
					}
				case float32:
					for i, val := range vals {
						svals[i] = fmtFloat32(val.(float32))
					}
				case int:
					for i, val := range vals {
						svals[i] = strconv.Itoa(val.(int))
					}
				default:
					for i, val := range vals {
						svals[i] = fmt.Sprintf("%v", val)
					}
				}

				s += fmt.Sprintf("%s=%s", k, strings.Join(svals, ","))
			default:
				s += fmt.Sprintf("%s=%s", k, v.(string))
			}
		}
		if j < len(order)-1 && !strings.HasSuffix(s, ";") {
			s += ";"
		}
	}
	return s
}

func (i InfoMap) Keys() []string {
	if keys, ok := i["__order"]; ok {
		return keys.([]string)
	}
	keys := make([]string, 0)
	for k := range i {
		keys = append(keys, k)
	}
	return keys
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
	if v.Alt[0][0] != '<' {
		return uint32(v.Pos-1) + uint32(len(v.Ref))
	}
	if strings.HasPrefix(v.Alt[0], "<DEL") || strings.HasPrefix(v.Alt[0], "<DUP") || strings.HasPrefix(v.Alt[0], "<INV") || strings.HasPrefix(v.Alt[0], "<CN") {
		if svlen, ok := v.Info["SVLEN"]; ok {
			var slen int
			switch svlen.(type) {
			case int:
				slen = svlen.(int)
			case []interface{}:
				slen = svlen.([]interface{})[0].(int)
			case interface{}:
				slen = svlen.(int)
			default:
				log.Fatalf("non int type for SVLEN")
			}
			if slen < 0 {
				slen = -slen
			}
			return uint32(int(v.Pos) + slen)

		} else if end, ok := v.Info["END"]; ok {
			return uint32(end.(int))
		}
		log.Fatalf("no svlen for variant %s:%d", v.Chromosome, v.Pos)
	}
	// <INS and BND's get handled by this.
	return uint32(v.Pos)
}

func fmtFloat32(v float32) string {
	var val string
	if v > 0.02 || v < -0.02 {
		val = fmt.Sprintf("%.4f", v)
	} else {
		val = fmt.Sprintf("%.5g", v)
	}
	val = strings.TrimRight(strings.TrimRight(val, "0"), ".")
	if val == "" {
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
	if val == "" {
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

// RefDepths returns the depths of the alternates for this sample
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

// String returns the string representation of the sample field.
func (sg *SampleGenotype) String(fields []string) string {
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
	if v.Quality == MISSING_VAL {
		qual = "."
	} else {
		qual = fmt.Sprintf("%.1f", v.Quality)
	}
	s := fmt.Sprintf("%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s", v.Chromosome, v.Pos, v.Id, v.Ref, strings.Join(v.Alt, ","), qual, v.Filter, v.Info)
	if len(v.Samples) > 0 {
		samps := make([]string, len(v.Samples))
		for i, s := range v.Samples {
			samps[i] = s.String(v.Format)
		}
		s += fmt.Sprintf("\t%s\t%s", strings.Join(v.Format, ":"), strings.Join(samps, "\t"))
	} else if v.sampleStrings != nil && len(v.sampleStrings) != 0 {
		s += fmt.Sprintf("\t%s\t%s", strings.Join(v.Format, ":"), strings.Join(v.sampleStrings, "\t"))
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

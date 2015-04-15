package vcfgo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var typeRe = `String|Integer|Float|Flag|Character|Unknown`
var infoRegexp = regexp.MustCompile(fmt.Sprintf(`##INFO=<ID=(.+),Number=([\dAGR\.]?),Type=(%s),Description="(.+)">`, typeRe))
var formatRegexp = regexp.MustCompile(fmt.Sprintf(`##FORMAT=<ID=(.+),Number=([\dAGR\.]?),Type=(%s),Description="(.+)">`, typeRe))
var filterRegexp = regexp.MustCompile(`##FILTER=<ID=(.+),Description="(.+)">`)
var contigRegexp = regexp.MustCompile(`contig=<.*((\w+)=([^,>]+))`)
var fileVersionRegexp = regexp.MustCompile(`##fileformat=VCFv(.+)`)

// Info holds the Info and Format fields
type Info struct {
	Id          string
	Description string
	Number      string // A G R . ''
	Type        string // STRING INTEGER FLOAT FLAG CHARACTER UNKONWN
}

// SampleFormat holds the type info for Format fields.
type SampleFormat Info

// Header holds all the type and format information for the variants.
type Header struct {
	SampleNames   []string
	Infos         map[string]*Info
	SampleFormats map[string]*SampleFormat
	Filters       map[string]string
	Extras        map[string]string
	FileFormat    string
	// contid id maps to a map of length, URL, etc.
	Contigs map[string]map[string]string
	// TOOD: make Contig struct
}

// String returns a string representation.
func (i *Info) String() string {
	return fmt.Sprintf("##INFO=<ID=%s,Number=%s,Type=%s,Description=\"%s\">", i.Id, i.Number, i.Type, i.Description)
}

// String returns a string representation.
func (i *SampleFormat) String() string {
	return fmt.Sprintf("##FORMAT=<ID=%s,Number=%s,Type=%s,Description=\"%s\">", i.Id, i.Number, i.Type, i.Description)
}

/*
func (h *Header) Validate(verr *VCFError) []error {
	var errs []error
	return errs
}*/

func (h *Header) parseSample(format []string, s string) (*SampleGenotype, []error) {
	values := strings.Split(s, ":")
	if len(format) != len(values) {
		return nil, []error{fmt.Errorf("bad sample string: %s", s)}
	}
	//if geno == nil {
	var value string
	var geno = NewSampleGenotype()
	var errs []error
	//}
	var e error
	for i, field := range format {
		value = values[i]
		switch field {
		case "GT":
			e = h.setSampleGT(geno, value)
		case "DP":
			e = h.setSampleDP(geno, value)
		case "GL":
			e = h.setSampleGL(geno, value)
		case "GQ":
			if format, ok := h.SampleFormats[field]; ok {
				e = h.setSampleGQ(geno, value, format.Type)
			}
		}
		geno.Fields[field] = value
		if e != nil {
			errs = append(errs, e)
		}
	}
	return geno, errs
}

func (h *Header) setSampleDP(geno *SampleGenotype, value string) error {
	var err error
	geno.DP, err = strconv.Atoi(value)
	if err != nil && value == "" || value == "." {
		return nil
	}
	return err
}

func (h *Header) setSampleGQ(geno *SampleGenotype, value string, Type string) error {
	var err error
	if Type == "Integer" {
		geno.GQ, err = strconv.Atoi(value)
	} else if Type == "Float" {
		var v float64
		v, err = strconv.ParseFloat(value, 32)
		if err == nil {
			err = errors.New("setSampleGQ: GQ reported as float. rounding to int")
			geno.GQ = int(math.Floor(v + 0.5))
		}
	}

	if err != nil && (value == "" || value == ".") {
		return nil
	}
	return err
}

func (h *Header) setSampleGL(geno *SampleGenotype, value string) error {
	var err error
	if len(geno.GL) != 0 {
		geno.GL = geno.GL[:0]
	}
	vals := strings.Split(value, ",")
	var v float64
	for _, val := range vals {
		v, err = strconv.ParseFloat(val, 32)
		/*if err != nil {
			return err
		}*/
		geno.GL = append(geno.GL, float32(v))
	}
	return err
}

func (h *Header) setSampleGT(geno *SampleGenotype, value string) error {
	if len(geno.GT) != 0 {
		geno.GT = geno.GT[:0]
	}
	if (value[1] == '/' || value[1] == '|') && len(value) == 3 {
		geno.Phased = value[1] == '|'
		var v int
		if value[0] == '.' {
			v = -1
		} else {
			v = int(value[0] - '0') // convert to int the old-fashioned way
		}
		geno.GT = append(geno.GT, v)
		if value[2] == '.' {
			v = -1
		} else {
			v = int(value[2] - '0') // convert to int
		}
		geno.GT = append(geno.GT, v)
	} else if value == "." {
		geno.GT = []int{-1, -1}
		geno.Phased = false
	} else {
		geno.Phased = strings.Contains(value, "|")
		splitString := "/"
		if geno.Phased {
			splitString = "|"
		}
		alleles := strings.Split(value, splitString)
		var v int
		for _, allele := range alleles {
			if allele == "." {
				v = -1
			} else {
				v = int(v - '0')
			}
			geno.GT = append(geno.GT, v)
		}

	}
	return nil
}

// NewHeader returns a Header with the requisite allocations.
func NewHeader() *Header {
	var h Header
	h.Filters = make(map[string]string)
	h.Infos = make(map[string]*Info)
	h.SampleFormats = make(map[string]*SampleFormat)
	h.SampleNames = make([]string, 0)
	h.Extras = make(map[string]string)
	h.Contigs = make(map[string]map[string]string)
	return &h
}

func parseHeaderInfo(info string) (*Info, error) {
	res := infoRegexp.FindStringSubmatch(info)
	if len(res) != 5 {
		return nil, fmt.Errorf("INFO error: %s, %r", info, res)
	}
	var i Info
	i.Id = res[1]
	i.Number = res[2]
	i.Type = res[3]
	i.Description = res[4]
	return &i, nil
}

func parseHeaderContig(contig string) (map[string]string, error) {
	vmap := make(map[string]string)
	contig = strings.TrimSuffix(strings.TrimPrefix(contig, "##contig=<"), ">")
	rdr := csv.NewReader(strings.NewReader(contig))
	rdr.LazyQuotes = true
	rdr.TrimLeadingSpace = true
	contigs, err := rdr.Read()

	for _, pair := range contigs {
		kv := strings.SplitN(pair, "=", 2)
		vmap[kv[0]] = kv[1]
	}
	return vmap, err
}

func parseHeaderExtraKV(kv string) ([]string, error) {
	kv = strings.TrimLeft(kv, "##")
	kv = strings.TrimLeft(kv, " ")
	kvpair := strings.SplitN(kv, "=", 2)

	if len(kvpair) != 2 {
		return nil, fmt.Errorf("header error in extra field: %s", kv)
	}
	return kvpair, nil
}

func parseHeaderFormat(info string) (*SampleFormat, error) {
	res := formatRegexp.FindStringSubmatch(info)
	if len(res) != 5 {
		return nil, fmt.Errorf("FORMAT error: %s", info)
	}
	var i SampleFormat
	i.Id = res[1]
	i.Number = res[2]
	i.Type = res[3]
	i.Description = res[4]
	return &i, nil
}

func parseHeaderFilter(info string) ([]string, error) {
	res := filterRegexp.FindStringSubmatch(info)
	if len(res) != 3 {
		return nil, fmt.Errorf("FILTER error: %s", info)
	}
	return res[1:3], nil
}

func parseHeaderFileVersion(format string) (string, error) {
	res := fileVersionRegexp.FindStringSubmatch(format)
	if len(res) != 2 {
		return "-1", fmt.Errorf("file format error: %s", format)
	}

	return res[1], nil
}

func parseSampleLine(line string) ([]string, error) {
	fields := strings.Split(line, "\t")
	var samples []string
	if len(fields) > 9 {
		samples = fields[9:]
	} else {
		samples = []string{}
	}
	return samples, nil
}

func parseOne(key, val, itype string) (interface{}, error) {
	var v interface{}
	var err error
	switch itype {
	case "Integer":
		v, err = strconv.Atoi(val)
	case "Float":
		v, err = strconv.ParseFloat(val, 32)
	case "Flag":
		if val != "" {
			err = fmt.Errorf("Info Error: flag field (%s) had value", key)
		}
		v = true
	default:
		v = val
	}
	return v, err
}

func (h *Header) parseInfo(info string) (InfoMap, error) {
	toks := strings.Split(info, ";")
	m := make(InfoMap)
	var order []string // keep order of fields so we print in same order
	var err error
	var v, iv interface{}
	for _, tok := range toks {
		// i from the header can tell us about the type of this
		pair := strings.Split(tok, "=")
		i, ok := h.Infos[pair[0]]
		order = append(order, pair[0])
		if !ok {
			err = fmt.Errorf("Info Error: %s not found in header", tok)
			m[pair[0]] = pair[1]
			continue
		}
		if len(pair) == 1 {
			if i.Type != "Flag" {
				err = fmt.Errorf("Info Error: flag field (%s) had value", pair[0])
			}
			pair = append(pair, "")
		}
		switch i.Number {
		case "0":
			if i.Type != "Flag" {
				err = fmt.Errorf("Info Error: flag field (%s) should have Number=0", pair[0])
			}
			v = true
		case "1":
			v, err = parseOne(pair[0], pair[1], i.Type)
		case "R", "A", "G", "2", "3", ".":
			vals := strings.Split(pair[1], ",")
			v = make([]interface{}, len(vals))
			for j, val := range vals {
				iv, err = parseOne(pair[0], val, i.Type)
				v.([]interface{})[j] = iv
			}
		default:
			panic(fmt.Sprintf("found Number=%s", i.Number))

		}
		m[pair[0]] = v
	}
	m["__order"] = order
	return m, err
}

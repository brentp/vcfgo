package vcfgo

import (
	"fmt"
	"strconv"
)

func SplitAlts(v *Variant) []*Variant {
	vars := make([]*Variant, len(v.Alt()))
	for i := range v.Alt() {
		vars[i] = &Variant{Chromosome: v.Chromosome, Pos: v.Pos, Id_: v.Id_,
			Reference: v.Ref(), Alternate: []string{v.Alt()[i]}, Quality: v.Quality, Filter: v.Filter,
			Info_: v.Info_, Samples: v.Samples, sampleString: v.sampleString,
			Header: v.Header, LineNumber: v.LineNumber}

		split(vars[i], i, len(v.Alt()))
	}
	return vars
}

func split(v *Variant, i int, nAlts int) error {
	// extract the appropriate
	var err error
	for _, k := range v.Info_.Keys() {
		h := v.Header.Infos[k]
		switch h.Number {
		case "A":
			var s interface{}
			s, err = splitA(v, i, nAlts)
			v.Info_.Set(k, s)
		case "G":
			var s []interface{}
			var val interface{}
			val, err = v.Info_.Get(k)
			s, err = splitG(val, i, nAlts)
			v.Info_.Set(k, s)
		case "R":
			var s []interface{}
			var val interface{}
			val, err = v.Info_.Get(k)
			s, err = splitR(val, i, nAlts)
			v.Info_.Set(k, s)
		}
	}
	v.Header.ParseSamples(v)
	for _, h := range v.Header.SampleFormats {
		if h.Number != "A" && h.Number != "G" && h.Number != "R" {
			continue
		}
		/*
			for j, samp := range v.Samples {
				switch h.Number {
				case "A":
					// TODO: should this still be a list?
					var s interface{}
					s, err = splitA(samp.Fields[k], i, nAlts)
					//samp.Fields[k] = s

				case "G":

				case "R":
				}
			}
		*/
	}
	// TODO: v.Samples
	return err
}

func splitG(m interface{}, i int, nAlts int) ([]interface{}, error) {
	pairs := [][]int{{0, 0}, {0, i}, {i, i}}
	G := make([]interface{}, 3)
	for o, jk := range pairs {
		order := (jk[1] * (jk[1] + 1) / 2) + jk[0]
		G[o] = m.([]interface{})[order]
	}
	return G, nil
}

func splitGT(m interface{}, i int) ([]interface{}, error) {
	ml := m.([]interface{})
	out := make([]interface{}, len(ml))
	ai := strconv.Itoa(i)
	for i, allelei := range ml {
		allele := allelei.(string)
		if allele == "0" {
			out[i] = "0"
		}
		if ai == allele {
			out[i] = "1"
		} else {
			out[i] = "."
		}
	}
	return out, nil
}

func splitR(m interface{}, i int, nAlts int) ([]interface{}, error) {
	ml := m.([]interface{})
	if len(ml) != nAlts+1 {
		return nil, fmt.Errorf("incorrect number of alts in splitR: %v", m)
	}
	if i+1 >= len(ml) {
		return nil, fmt.Errorf("requested bad alt in splitR: %v", m)
	}
	out := make([]interface{}, 2)
	out[0] = ml[0]
	out[1] = ml[i+1]
	return out, nil
}

// when we split on alt, it returns a single int since we are decomposing.
func splitA(m interface{}, i int, nAlts int) (interface{}, error) {
	ml := m.([]interface{})
	if len(ml) != nAlts {
		return nil, fmt.Errorf("incorrect number of alts in splitA: %v", m)
	}
	return ml[i], nil
}

package vcfgo

import "fmt"

func SplitAlts(v *Variant) []*Variant {
	vars := make([]*Variant, len(v.Alt))
	for i := range v.Alt {
		vars[i] = &Variant{Chromosome: v.Chromosome, Pos: v.Pos, Id: v.Id,
			Ref: v.Ref, Alt: []string{v.Alt[i]}, Quality: v.Quality, Filter: v.Filter,
			Info: v.Info, Samples: v.Samples, sampleStrings: v.sampleStrings,
			Header: v.Header, LineNumber: v.LineNumber}

		split(vars[i], i, len(v.Alt))
	}
	return vars
}

func split(v *Variant, i int, nAlts int) error {
	// extract the appropriate
	var err error
	for _, k := range v.Info.Keys() {
		h := v.Header.Infos[k]
		switch h.Number {
		case "A":
			var s interface{}
			s, err = splitA(v, i, nAlts)
			v.Info[k] = s
		case "G":
			var s []interface{}
			s, err = splitG(v.Info[k], i, nAlts)
			v.Info[k] = s
		case "R":
			var s []interface{}
			s, err = splitR(v.Info[k], i, nAlts)
			v.Info[k] = s
		}
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
	ai := fmt.Sprintf("%d", i)
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

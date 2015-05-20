package vcfgo

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

func split(v *Variant, i int, nAlts int) {
	// extract the appropirate GLs.
	pairs := [][]int{{0, 0}, {0, i}, {i, i}}
	for _, samp := range v.Samples {
		nGL := make([]float32, 3)
		for o, jk := range pairs {
			order := (jk[1] * (jk[1] + 1) / 2) + jk[0]
			nGL[o] = samp.GL[order]
		}
		samp.GL = nGL
	}

}

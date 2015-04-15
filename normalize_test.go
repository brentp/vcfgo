package vcfgo

import (
	"testing"
)

var leftaligntests = []struct {
	pos int
	ref []byte
	alt []byte
	seq []byte

	outPos int
	outRef string
	outAlt string
}{
	{123, []byte{'C', 'A', 'C'}, []byte{'C'}, []byte("GGGCACACAC"), 118, "GCA", "G"},
	{123, []byte{'C', 'A', 'C'}, []byte{'C'}, []byte("CACACAC"), 119, "CAC", "C"},
	{123, []byte{'C', 'C', 'A'}, []byte{'C', 'A', 'A'}, []byte("ACCCCCCA"), 123, "CC", "CA"},
	{123, []byte{'C'}, []byte{'A'}, []byte("ACCCCCC"), 123, "C", "A"},
}

var lefttrimtests = []struct {
	pos int
	ref []byte
	alt []byte

	outPos int
	outRef string
	outAlt string
}{
	{123, []byte{'C', 'C'}, []byte{'C', 'A'}, 124, "C", "A"},
	{123, []byte{'C', 'C'}, []byte{'C', 'C', 'T'}, 124, "C", "CT"},
	{123, []byte{'C', 'C', 'C'}, []byte{'C', 'C', 'C', 'T'}, 125, "C", "CT"},
	{123, []byte{'C'}, []byte{'T'}, 123, "C", "T"},
}

func TestLeftAlign(t *testing.T) {

	for _, v := range leftaligntests {
		p, r, a, err := leftalign(v.pos, v.ref, v.alt, v.seq)

		if p != v.outPos {
			t.Errorf("position should be %d instead of %d\n", v.outPos, p)
		}
		if string(r) != v.outRef {
			t.Errorf("ref should be %s instead of %s\n", v.outRef, string(r))
		}
		if string(a) != v.outAlt {
			t.Errorf("alt should be '%s' instead of %s", v.outAlt, a)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func TestLeftTrim(t *testing.T) {

	for _, v := range lefttrimtests {
		p, r, a, err := lefttrim(v.pos, v.ref, v.alt)

		if p != v.outPos {
			t.Errorf("position should be %d instead of %d\n", v.outPos, p)
		}
		if string(r) != v.outRef {
			t.Errorf("ref should be %s instead of %s\n", v.outRef, string(r))
		}
		if string(a) != v.outAlt {
			t.Errorf("alt should be '%s' instead of %s", v.outAlt, a)
		}
		if err != nil {
			t.Error(err)
		}

	}
}

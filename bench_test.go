package vcfgo

import (
	"os"
	"testing"
)

func benchmarkReader(lazy bool, b *testing.B) {

	for n := 0; n < b.N; n++ {
		f, err := os.Open("examples/test.query.vcf")
		if err != nil {
			panic(err)
		}
		rdr, err := NewReader(f, lazy)
		if err != nil {
			panic(err)
		}

		j := 0
		for {
			v := rdr.Read()
			if v == nil {
				break
			}
			j++
		}
	}
}

func BenchmarkLazy(b *testing.B)  { benchmarkReader(true, b) }
func BenchmarkEager(b *testing.B) { benchmarkReader(false, b) }

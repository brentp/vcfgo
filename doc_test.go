package vcfgo_test

import (
	"fmt"
	"os"

	"github.com/brentp/vcfgo"
)

func Example() {
	f, _ := os.Open("examples/test.auto_dom.no_parents.vcf")
	rdr, err := vcfgo.NewReader(f, false)
	if err != nil {
		panic(err)
	}
	for {
		variant := rdr.Read().(*vcfgo.Variant)
		if variant == nil {
			break
		}
		fmt.Printf("%s\t%d\t%s\t%s\n", variant.Chromosome, variant.Pos, variant.Ref, variant.Alt)
		dp, _ := variant.Info().Get("DP")
		fmt.Printf("%v", dp.(int) > 10)
		// Output: asdf
	}
	// Print all accumulated errors to stderr
	fmt.Fprintln(os.Stderr, rdr.Error())
}

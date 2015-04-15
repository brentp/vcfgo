package vcfgo_test

import (
	"fmt"
	"github.com/brentp/vcfgo"
	"os"
)

func Example() {
	f, _ := os.Open("examples/test.auto_dom.no_parents.vcf")
	rdr, err := vcfgo.NewReader(f)
	if err != nil {
		panic(err)
	}
	for {
		variant := rdr.Read()
		if variant == nil {
			break
		}
		fmt.Printf("%s\t%d\t%s\t%s\n", variant.Chromosome, variant.Pos, variant.Ref, variant.Alt)
		fmt.Printf("%s", variant.Info["DP"].(int) > 10)
		// Output: asdf
	}
	// Print all accumulated errors to stderr
	fmt.Fprintln(os.Stderr, rdr.Error())
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/brentp/vcfgo"
)

func main() {

	flag.Parse()
	files := flag.Args()
	f, err := os.Open(files[0])

	r := io.Reader(f)
	vr, err := vcfgo.NewReader(r, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", vr)

	variant := vr.Read()
	fmt.Println(vr.Error())
	fmt.Println("variant:", variant)
	if len(variant.Samples) > 0 {
		if _, ok := vr.Header.SampleFormats["PL"]; ok {
			if vr.Header.SampleFormats["PL"].Type == "Integer" {
				fmt.Println(variant.GetGenotypeField(variant.Samples[0], "PL", int(-1)))
			} else {
				fmt.Println(variant.GetGenotypeField(variant.Samples[0], "PL", float32(-1)))
			}
		}
	}
	fmt.Println(vr.Error())
	vr.Clear()
	for {
		variant = vr.Read()
		if variant == nil {
			if e := vr.Error(); e != io.EOF && e != nil {
				vr.Clear()
			}
			break
		}
		if vr.Error() != nil {
			fmt.Println(vr.Error())
		}
		vr.Clear()
		if len(variant.Samples) > 0 {
			var pl interface{}
			if _, ok := vr.Header.SampleFormats["PL"]; ok {
				if vr.Header.SampleFormats["PL"].Type == "Integer" {
					pl, err = variant.GetGenotypeField(variant.Samples[0], "PL", int(-1))
				} else {
					pl, err = variant.GetGenotypeField(variant.Samples[0], "PL", float32(-1))
				}
			}
			fmt.Println("ERR:", err)
			fmt.Println(variant.Samples[0])
			if err != nil && variant.Samples[0] != nil {
				log.Println("BBBBBBBBBBBBBBBBBBB")
				if _, ok := vr.Header.SampleFormats["PL"]; ok {
					fmt.Println("")
					fmt.Println(variant.Samples[0])
					log.Println("DDDDDDDDDDDDDDDDD")
					log.Fatal(err)
				}
			}

			if variant.Samples[0] != nil {
				fmt.Println("PL:", pl, "GQ:", variant.Samples[0].GQ, "DP:", variant.Samples[0].DP)
			}
		}
	}
	fmt.Println("OK")
}

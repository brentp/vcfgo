package vcfgo

import (
	"fmt"
	"io"
	"strings"

	. "gopkg.in/check.v1"
)

var regr1 = `##fileformat=VCFv4.1
##INFO=<ID=AC,Number=A,Type=Integer,Description="Total number of alternate alleles in called genotypes">
##INFO=<ID=AN,Number=1,Type=Integer,Description="Total number of alleles in called genotypes">
##INFO=<ID=AF,Number=A,Type=Float,Description="Estimated allele frequency in the range (0,1]">
##INFO=<ID=AO,Number=A,Type=Integer,Description="Alternate allele observations, with partial observations recorded fractionally">
##INFO=<ID=PAO,Number=A,Type=Float,Description="Alternate allele observations, with partial observations recorded fractionally">
##INFO=<ID=SAF,Number=A,Type=Integer,Description="Number of alternate observations on the forward strand">
##INFO=<ID=SAP,Number=A,Type=Float,Description="Strand balance probability for the alternate allele: Phred-scaled upper-bounds estimate of the probability of observing the deviation between SAF and SAR given E(SAF/SAR) ~ 0.5, derived using Hoeffding's inequality">
##INFO=<ID=AB,Number=A,Type=Float,Description="Allele balance at heterozygous sites: a number between 0 and 1 representing the ratio of reads showing the reference allele to all reads, considering only reads from individuals called as heterozygous">
##INFO=<ID=ABP,Number=A,Type=Float,Description="Allele balance probability at heterozygous sites: Phred-scaled upper-bounds estimate of the probability of observing the deviation between ABR and ABA given E(ABR/ABA) ~ 0.5, derived using Hoeffding's inequality">
##INFO=<ID=XX,Number=2,Type=Float,Description="test mult vals">
##INFO=<ID=TYPE,Number=A,Type=String,Description="The type of allele, either snp, mnp, ins, del, or complex.">
##INFO=<ID=CIGAR,Number=A,Type=String,Description="The extended CIGAR representation of each alternate allele, with the exception that '=' is replaced by 'M' to ease VCF parsing.  Note that INDEL alleles do not have the first matched base (which is provided by default, per the spec) referred to by the CIGAR.">
##INFO=<ID=MEANALT,Number=A,Type=Float,Description="Mean number of unique non-reference allele observations per sample with the corresponding alternate alleles.">
##FORMAT=<ID=AO,Number=A,Type=Integer,Description="Alternate allele observation count">
##INFO=<ID=CSQ,Number=.,Type=String,Description="Consequence type as predicted by VEP. Format: Consequence|Codons|Amino_acids|Gene|SYMBOL|Feature|EXON|PolyPhen|SIFT|Protein_position|BIOTYPE">
##INFO=<ID=OLD_VARIANT,Number=1,Type=String,Description="Original chr:pos:ref:alt encoding">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO
1	98683	.	G	A	610.487	.	AB=0.282443;ABP=56.8661;AC=11;AF=0.34375;AN=32;AO=45;CIGAR=1X;TYPE=snp;XX=0.44,0.88`

type RegressionSuite struct {
	reader io.Reader
	vcfStr string
}

var _ = Suite(&RegressionSuite{vcfStr: regr1})

func (s *RegressionSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(s.vcfStr)
}

func (s *RegressionSuite) TestRegr1(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read()
	snp, ok := v.Info["TYPE"]
	c.Assert(ok, Equals, true)
	c.Assert(snp, DeepEquals, []interface{}{"snp"})

	str := fmt.Sprintf("%s", v)
	c.Assert(str, Equals, "1\t98683\t.\tG\tA\t610.5\t.\tAB=0.2824;ABP=56.8661;AC=11;AF=0.3438;AN=32;AO=45;CIGAR=1X;TYPE=snp;XX=0.44,0.88")
}

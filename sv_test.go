package vcfgo_test

import (
	"io"
	"strings"
	"testing"

	"github.com/brentp/vcfgo"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var cnvStr = `##fileformat=VCFv4.1
##fileDate=20100501
##reference=1000GenomesPilot-NCBI36
##assembly=ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/sv/breakpoint_assemblies.fasta
##INFO=<ID=BKPTID,Number=.,Type=String,Description="ID of the assembled alternate allele in the assembly file">
##INFO=<ID=CIEND,Number=2,Type=Integer,Description="Confidence interval around END for imprecise variants">
##INFO=<ID=CIPOS,Number=2,Type=Integer,Description="Confidence interval around POS for imprecise variants">
##INFO=<ID=END,Number=1,Type=Integer,Description="End position of the variant described in this record">
##INFO=<ID=HOMLEN,Number=.,Type=Integer,Description="Length of base pair identical micro-homology at event breakpoints">
##INFO=<ID=HOMSEQ,Number=.,Type=String,Description="Sequence of base pair identical micro-homology at event breakpoints">
##INFO=<ID=SVLEN,Number=.,Type=Integer,Description="Difference in length between REF and ALT alleles">
##INFO=<ID=SVTYPE,Number=1,Type=String,Description="Type of structural variant">
##ALT=<ID=DEL,Description="Deletion">
##ALT=<ID=DEL:ME:ALU,Description="Deletion of ALU element">
##ALT=<ID=DEL:ME:L1,Description="Deletion of L1 element">
##ALT=<ID=DUP,Description="Duplication">
##ALT=<ID=DUP:TANDEM,Description="Tandem Duplication">
##ALT=<ID=INS,Description="Insertion of novel sequence">
##ALT=<ID=INS:ME:ALU,Description="Insertion of ALU element">
##ALT=<ID=INS:ME:L1,Description="Insertion of L1 element">
##ALT=<ID=INV,Description="Inversion">
##ALT=<ID=CNV,Description="Copy number variable region">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
##FORMAT=<ID=GQ,Number=1,Type=Float,Description="Genotype quality">
##FORMAT=<ID=CN,Number=1,Type=Integer,Description="Copy number genotype for imprecise events">
##FORMAT=<ID=CNQ,Number=1,Type=Float,Description="Copy number genotype quality for imprecise events">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001
2	321682	.	T	<DEL>	6	PASS	SVTYPE=DEL;END=321887;SVLEN=-205;CIPOS=-56,20;CIEND=-10,62	GT:GQ	0/1:12
2	14477084	.	C	<DEL:ME:ALU>	12	PASS	SVTYPE=DEL;END=14477381;CIPOS=-22,18;CIEND=-12,32	GT:GQ	0/1:12
3	9425916	.	C	<INS:ME:L1>	23	PASS	SVTYPE=INS;SVLEN=6027;CIPOS=-16,22	GT:GQ	1/1:15
3	12665100	.	A	<DUP>	14	PASS	SVTYPE=DUP;END=12686200;SVLEN=21100;CIPOS=-500,500;CIEND=-500,500	GT:GQ:CN:CNQ	./.:0:3:16.2
4	18665128	.	T	<DUP:TANDEM>	11	PASS	SVTYPE=DUP;END=18665204;SVLEN=76;CIPOS=-10,10;CIEND=-10,10	GT:GQ:CN:CNQ	./.:0:5:8.3
4	18665128	.	T	<INV>	11	PASS	SVTYPE=INS;END=18665204;SVLEN=76;CIPOS=-10,10;CIEND=-10,10	GT:GQ:CN:CNQ	./.:0:5:8.3
4	18665128	.	T	<CNV>	11	PASS	SVTYPE=DUP;END=18665204;SVLEN=76;CIPOS=-10,10;CIEND=-10,10	GT:GQ:CN:CNQ	./.:0:5:8.3
4	43266825	rs369548269	T	<DEL>	.	.	ON=0;OLD_MULTIALLELIC=4:43266825:T/TAC/<DEL>	GT	0/1
5	755892	175981_1	N	]1:759001]N	112.69	PASS	AC=4;AF=0.01361;ALG=PROD;AN=294;CIEND=-10,29;CIEND95=0,34;CIPOS=-30,29;CIPOS95=28,32;EVENT=175981;EXVAR=0.00508;HW=0.00;HetPct=2.7;HomRefPct=97.3;HomVarPct=0.0;IMPRECISE;MATEID=1	GT	0/1`

//1	2827694	rs2376870	CGTGGATGCGGGGAC	C	.	PASS	SVTYPE=DEL;END=2827762;HOMLEN=1;HOMSEQ=G;SVLEN=-68	GT:GQ	1/1:13.9

type CNVSuite struct {
	reader io.Reader
	vcfStr string
}

var _ = Suite(&CNVSuite{vcfStr: cnvStr})

func (s *CNVSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(s.vcfStr)
}

func (s *CNVSuite) TestDupIns(c *C) {
	r, err := vcfgo.NewReader(s.reader, false)
	c.Assert(err, IsNil)
	var v *vcfgo.Variant

	v = r.Read() //.(*vcfgo.Variant)
	c.Assert(int(v.End()), Equals, 321887)

	left, right, ok := v.CIPos()
	c.Assert(ok, Equals, true)
	c.Assert(int(left), Equals, 321682-56-1)
	c.Assert(int(right), Equals, 321682+20)

	eleft, eright, ok := v.CIEnd()
	c.Assert(ok, Equals, true)
	c.Assert(int(eleft), Equals, 321887-10-1)
	c.Assert(int(eright), Equals, 321887+62)

	v = r.Read() //.(*vcfgo.Variant)
	c.Assert(int(v.End()), Equals, 14477381)

	v = r.Read() //.(*vcfgo.Variant)
	c.Assert(int(v.Start()), Equals, 9425915)
	c.Assert(int(v.End()), Equals, 9431943)

	v = r.Read() //.(*vcfgo.Variant)
	c.Assert(int(v.End()), Equals, 12686200)

	v = r.Read() //.(*vcfgo.Variant)
	c.Assert(int(v.End()), Equals, 18665204)

	v = r.Read() //.(*vcfgo.Variant) // INS
	c.Assert(int(v.End()), Equals, 18665204)

	v = r.Read() //.(*vcfgo.Variant) // CNV
	c.Assert(int(v.End()), Equals, 18665204)

	v = r.Read() //.(*vcfgo.Variant) // CNV
	c.Assert(int(v.End()), Equals, 43266825)

	left, right, ok = v.CIPos()
	c.Assert(left, Equals, v.Start())
	c.Assert(right, Equals, v.Start()+uint32(1))
	c.Assert(ok, Equals, false)

	eleft, eright, ok = v.CIEnd()
	c.Assert(eleft, Equals, v.End()-1)
	c.Assert(eright, Equals, v.End())
	c.Assert(ok, Equals, false)

	v = r.Read() //.(*vcfgo.Variant) // BND
	c.Assert(int(v.Start()), Equals, 755891)
	c.Assert(int(v.End()), Equals, 755891+1)

}

var svlenNumberAStr = `##fileformat=VCFv4.1
##fileDate=20100501
##reference=1000GenomesPilot-NCBI36
##assembly=ftp://ftp-trace.ncbi.nih.gov/1000genomes/ftp/release/sv/breakpoint_assemblies.fasta
##INFO=<ID=SVLEN,Number=A,Type=Integer,Description="Length of structural variant">
##INFO=<ID=SVTYPE,Number=1,Type=String,Description="Type of structural variant">
##ALT=<ID=DEL,Description="Deletion">
##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001
1	1000	.	T	<DEL>	60	PASS	SVTYPE=DEL;SVLEN=100	GT	0/1
1	2000	.	A	<DEL>,<DEL>	60	PASS	SVTYPE=DEL;SVLEN=150,200	GT	1/2
`

type SVLENNumberASuite struct {
	reader io.Reader
	vcfStr string
}

var _ = Suite(&SVLENNumberASuite{vcfStr: svlenNumberAStr})

func (s *SVLENNumberASuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(s.vcfStr)
}

func (s *SVLENNumberASuite) TestSVLENNumberA(c *C) {
	r, err := vcfgo.NewReader(s.reader, false)
	c.Assert(err, IsNil)
	var v *vcfgo.Variant

	// Test single alternative allele with SVLEN Number=A
	v = r.Read()
	c.Assert(v, NotNil)
	c.Assert(int(v.Start()), Equals, 999) // 0-based
	c.Assert(int(v.End()), Equals, 1100)  // 1000 + 100

	// Test multiple alternative alleles with SVLEN Number=A
	v = r.Read()
	c.Assert(v, NotNil)
	c.Assert(int(v.Start()), Equals, 1999) // 0-based
	c.Assert(int(v.End()), Equals, 2150)   // 2000 + 150 (first alt allele)
}

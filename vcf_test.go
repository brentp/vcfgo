package vcfgo

import (
	. "gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type VCFSuite struct{}

var _ = Suite(&VCFSuite{})

var infotests = []struct {
	input string
	exp   *Info
}{
	{`##INFO=<ID=NS,Number=1,Type=Integer,Description="Number of Samples With Data">`,
		&Info{Id: "NS", Number: "1", Type: "Integer", Description: "Number of Samples With Data"}},
	{`##INFO=<ID=DP,Number=1,Type=Integer,Description="Total Depth">`,
		&Info{Id: "DP", Number: "1", Type: "Integer", Description: "Total Depth"}},
	{`##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">`,
		&Info{Id: "AF", Number: "A", Type: "Float", Description: "Allele Frequency"}},
	{`##INFO=<ID=AA,Number=1,Type=String,Description="Ancestral Allele">`,
		&Info{Id: "AA", Number: "1", Type: "String", Description: "Ancestral Allele"}},
	{`##INFO=<ID=DB,Number=0,Type=Flag,Description="dbSNP membership, build 129">`,
		&Info{Id: "DB", Number: "0", Type: "Flag", Description: "dbSNP membership, build 129"}},
	{`##INFO=<ID=H2,Number=0,Type=Flag,Description="HapMap2 membership">`,
		&Info{Id: "H2", Number: "0", Type: "Flag", Description: "HapMap2 membership"}},
}

var formattests = []struct {
	input string
	exp   *SampleFormat
}{
	{`##FORMAT=<ID=GT,Number=1,Type=String,Description="Genotype">`,
		&SampleFormat{Id: "GT", Number: "1", Type: "String", Description: "Genotype"}},
	{`##FORMAT=<ID=GQ,Number=1,Type=Integer,Description="Genotype Quality">`,
		&SampleFormat{Id: "GQ", Number: "1", Type: "Integer", Description: "Genotype Quality"}},
	{`##FORMAT=<ID=HQ,Number=2,Type=Integer,Description="Haplotype Quality">`,
		&SampleFormat{Id: "HQ", Number: "2", Type: "Integer", Description: "Haplotype Quality"}},
	{`##FORMAT=<ID=DP,Number=1,Type=Integer,Description="Read Depth">`,
		&SampleFormat{Id: "DP", Number: "1", Type: "Integer", Description: "Read Depth"}},
}

var filtertests = []struct {
	filter string
	exp    []string
}{
	{`##FILTER=<ID=q10,Description="Quality below 10">`,
		[]string{"q10", "Quality below 10"}},
	{`##FILTER=<ID=s50,Description="Less than 50% of samples have data">`,
		[]string{"s50", "Less than 50% of samples have data"}},
}

var samplelinetests = []struct {
	line string
	exp  []string
}{
	{`#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	NA00001	NA00002	NA00003`, []string{"NA00001", "NA00002", "NA00003"}},
	{`#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT`, []string{}},
}

func (s *VCFSuite) TestHeaderInfoParse(c *C) {
	for _, v := range infotests {
		obs, err := parseHeaderInfo(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		c.Assert(obs.String(), Equals, v.input)
	}
}

func (s *VCFSuite) TestHeaderFormatParse(c *C) {
	for _, v := range formattests {
		obs, err := parseHeaderFormat(v.input)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
		c.Assert(obs.String(), Equals, v.input)

	}
}

func (s *VCFSuite) TestHeaderFilterParse(c *C) {

	for _, v := range filtertests {
		obs, err := parseHeaderFilter(v.filter)
		c.Assert(err, IsNil)
		c.Assert(obs, DeepEquals, v.exp)
	}
}

func (s *VCFSuite) TestHeaderVersionParse(c *C) {
	obs, err := parseHeaderFileVersion(`##fileformat=VCFv4.2`)
	c.Assert(err, IsNil)
	c.Assert(obs, Equals, "4.2")
}

func (s *VCFSuite) TestHeaderBadVersionParse(c *C) {
	_, err := parseHeaderFileVersion(`##fileformat=VFv4.2`)
	c.Assert(err, ErrorMatches, "file format error.*")
}

func (s *VCFSuite) TestHeaderContigParse(c *C) {
	m, err := parseHeaderContig(`##contig=<ID=20,length=62435964,assembly=B36,md5=f126cdf8a6e0c7f379d618ff66beb2da,species="Homo sapiens",taxonomy=x>`)
	c.Assert(err, IsNil)
	c.Assert(m, DeepEquals, map[string]string{"assembly": "B36", "md5": "f126cdf8a6e0c7f379d618ff66beb2da", "species": "\"Homo sapiens\"", "taxonomy": "x", "ID": "20", "length": "62435964"})
}

func (s *VCFSuite) TestHeaderExtra(c *C) {
	obs, err := parseHeaderExtraKV("##key=value")
	c.Assert(err, IsNil)
	c.Assert(obs[0], Equals, "key")
	c.Assert(obs[1], Equals, "value")
}

func (s *VCFSuite) TestHeaderSampleLine(c *C) {

	for _, v := range samplelinetests {
		r, err := parseSampleLine(v.line)
		c.Assert(err, IsNil)
		c.Assert(r, DeepEquals, v.exp)
	}
}

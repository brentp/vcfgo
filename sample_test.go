package vcfgo

import (
	"bytes"
	"io"
	"strings"

	. "gopkg.in/check.v1"
)

var sampleStr = `##fileformat=VCFv4.2
##fileDate=20090805
##source=myImputationProgramV3.1
##reference=file:///seq/references/1000GenomesPilot-NCBI36.fasta
##contig=<ID=20,length=62435964,assembly=B36,md5=f126cdf8a6e0c7f379d618ff66beb2da,species="Homo sapiens",taxonomy=x>
##phasing=partial
##INFO=<ID=AF,Number=A,Type=Float,Description="Allele Frequency">
##SAMPLE=<ID=Blood,Genomes=Germline,Mixture=1.,Description="Patient germline genome">
##SAMPLE=<ID=TissueSample,Genomes=Germline;Tumor,Mixture=.3;.7,Description="Patient germline genome;Patient tumor genome">
##PEDIGREE=<Name_0=G0-ID,Name_1=G1-ID,Name_N=GN-ID>
##PEDIGREE=<Name_0=G1-ID,Name_1=G1-ID,Name_N=GN-ID>
#CHROM	POS	ID	REF	ALT	QUAL	FILTER	INFO	FORMAT	BLOOD	TISSUE`

type SampleSuite struct {
	reader io.Reader
}

var _ = Suite(&SampleSuite{})

func (s *SampleSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(sampleStr)

}

func (s *SampleSuite) TestReaderHeaderSamples(c *C) {
	v, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	c.Assert(v.Header.SampleNames, DeepEquals, []string{"BLOOD", "TISSUE"})
}

func (s *SampleSuite) TestReaderSamples(c *C) {
	v, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	samps := v.Header.Samples
	c.Assert(len(samps), Equals, 2)

	for _, k := range []string{"Blood", "TissueSample"} {
		_, ok := samps[k]
		c.Assert(ok, Equals, true)
	}

}

func (s *SampleSuite) TestWriterSamples(c *C) {
	r, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)

	var wtr bytes.Buffer
	_, err = NewWriter(&wtr, r.Header)
	c.Assert(err, IsNil)

	str := wtr.String()

	c.Assert(strings.Contains(str, "\n##SAMPLE=<ID=Blood,"), Equals, true)
	c.Assert(strings.Contains(str, "\n##SAMPLE=<ID=TissueSample,"), Equals, true)

}

func (s *SampleSuite) TestWriterPedigree(c *C) {
	r, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)

	var wtr bytes.Buffer
	_, err = NewWriter(&wtr, r.Header)
	c.Assert(err, IsNil)

	str := wtr.String()

	c.Assert(strings.Contains(str, "\n##PEDIGREE=<Name_0=G0-ID,"), Equals, true)
	c.Assert(strings.Contains(str, "\n##PEDIGREE=<Name_0=G1-ID,"), Equals, true)

}

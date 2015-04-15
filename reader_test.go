package vcfgo

import (
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

type ReaderSuite struct {
	reader io.Reader
}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(vcfStr)

}

func (s *ReaderSuite) TestReaderHeaderSamples(c *C) {
	v, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	c.Assert(v.Header.SampleNames, DeepEquals, []string{"NA00001", "NA00002", "NA00003"})

}

func (s *ReaderSuite) TestReaderHeaderInfos(c *C) {
	v, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	c.Assert(v.Header.Infos["NS"], DeepEquals, &Info{Id: "NS", Number: "1", Type: "Integer", Description: "Number of Samples With Data"})
	c.Assert(v.Header.Filters["q10"], Equals, "Quality below 10")
	c.Assert(v.Header.SampleFormats["GT"], DeepEquals, &SampleFormat{Id: "GT", Number: "1", Type: "String", Description: "Genotype"})
}

func (s *ReaderSuite) TestReaderHeaderExtras(c *C) {
	v, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	c.Assert(len(v.Header.Extras), Equals, 4)
	c.Assert(v.Header.Extras["phasing"], Equals, "partial")

}

func (s *ReaderSuite) TestReaderRead(c *C) {
	rdr, err := NewReader(s.reader)
	c.Assert(err, IsNil)

	rec := rdr.Read()
	c.Assert(rec.Chromosome, Equals, "20")
	c.Assert(rec.Pos, Equals, uint64(14370))
	c.Assert(rec.Id, Equals, "rs6054257")
	c.Assert(rec.Ref, Equals, "G")
	c.Assert(rec.Alt[0], Equals, "A")
	c.Assert(rec.Quality, Equals, float32(29.0))
	c.Assert(rec.Filter, Equals, "PASS")

	//20	17330	.	T	A	3	q10	NS=3;DP=11;AF=0.017	GT:GQ:DP:HQ	0|0:49:3:58,50	0|1:3:5:65,3	0/0:41:3
	rec0 := rdr.Read()
	c.Assert(rec0.Chromosome, Equals, "20")
	c.Assert(rec0.Pos, Equals, uint64(17330))
	c.Assert(rec0.Id, Equals, ".")
	c.Assert(rec0.Ref, Equals, "T")
	c.Assert(rec0.Alt[0], Equals, "A")
	c.Assert(rec0.Quality, Equals, float32(3))
	c.Assert(rec0.Filter, Equals, "q10")

	//20	1110696	rs6040355	A	G,T	67	PASS	NS=2;DP=10;AF=0.333,0.667;AA=T;DB	GT:GQ:DP:HQ	1|2:21:6:23,27	2|1:2:0:18,2	2/2:35:4
	rec = rdr.Read()
	c.Assert(rec.Chromosome, Equals, "20")
	c.Assert(int(rec.Pos), Equals, 1110696)
	c.Assert(rec.Id, Equals, "rs6040355")
	c.Assert(rec.Ref, Equals, "A")
	c.Assert(rec.Alt, DeepEquals, []string{"G", "T"})
	c.Assert(rec.Quality, Equals, float32(67))
	c.Assert(rec.Filter, Equals, "PASS")

	c.Assert(rec0.Chromosome, Equals, "20")
	c.Assert(rec0.Pos, Equals, uint64(17330))
	c.Assert(rec0.Id, Equals, ".")
	c.Assert(rec0.Ref, Equals, "T")
	c.Assert(rec0.Alt[0], Equals, "A")
	c.Assert(rec0.Quality, Equals, float32(3))
	c.Assert(rec0.Filter, Equals, "q10")
}

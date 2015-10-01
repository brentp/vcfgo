package vcfgo_test

import (
	"fmt"
	"io"
	"strings"

	"github.com/brentp/vcfgo"

	. "gopkg.in/check.v1"
)

type VariantSuite struct {
	reader io.Reader
}

var _ = Suite(&VariantSuite{})

func (s *VariantSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(vcfStr)
}

func (s *VariantSuite) TestVariantGetInt(c *C) {
	rdr, err := vcfgo.NewReader(s.reader, true)
	c.Assert(err, IsNil)
	v := rdr.Read().(*vcfgo.Variant)
	v.Info()

	ns, err := v.Info_.Get("NS")
	c.Assert(err, IsNil)
	c.Assert(ns, Equals, 3)

	dp, err := v.Info_.Get("DP")
	c.Assert(dp, Equals, 14)
	c.Assert(err, IsNil)

	nsf, err := v.Info_.Get("NS")
	c.Assert(err, IsNil)
	c.Assert(nsf, Equals, int(3))

	dpf, err := v.Info_.Get("DP")
	c.Assert(err, IsNil)
	c.Assert(dpf, Equals, int(14))

	hqs, err := v.Info_.Get("AF")
	c.Assert(hqs, DeepEquals, []float32{0.5})
	c.Assert(err, IsNil)

	dpfs, err := v.Info_.Get("DP")
	c.Assert(err, IsNil)
	c.Assert(dpfs, DeepEquals, 14)

}

func (s *VariantSuite) TestInfoField(c *C) {
	rdr, err := vcfgo.NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read().(*vcfgo.Variant)
	vstr := fmt.Sprintf("%s", v.Info())
	c.Assert(vstr, Equals, "NS=3;DP=14;AF=0.5;DB;H2")
}

func (s *VariantSuite) TestInfoMap(c *C) {
	rdr, err := vcfgo.NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read().(*vcfgo.Variant)

	vstr := fmt.Sprintf("%s", v)
	c.Assert(vstr, Equals, "20\t14370\trs6054257\tG\tA\t29.0\tPASS\tNS=3;DP=14;AF=0.5;DB;H2\tGT:GQ:DP:HQ\t0|0:48:1:51,51\t1|0:48:8:51,51\t1/1:43:5:.,.")

	v.Info_.Set("asdf", 123)
	v.Info_.Set("float", 123.2001)
	has, err := v.Info_.Get("asdf")
	c.Assert(has, Equals, 123)
	val, err := v.Info_.Get("float")
	vv, ok := val.(float64)
	c.Assert(ok, Equals, true)
	c.Assert(vv-123.2001 < 1e-4 || 123.2001-vv < 1e-4, Equals, true)
	c.Assert(err, IsNil)

	c.Assert(fmt.Sprintf("%s", v.Info_), Equals, "NS=3;DP=14;AF=0.5;DB;H2;asdf=123;float=123.2001")

	rdr.Clear()

}

func (s *VariantSuite) TestStartEnd(c *C) {
	rdr, err := vcfgo.NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read()

	c.Assert(int(v.Start()), Equals, 14369)
	c.Assert(int(v.End()), Equals, 14370)
}

package vcfgo

import (
	"fmt"
	"io"
	"strings"

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
	rdr, err := NewReader(s.reader, true)
	c.Assert(err, IsNil)
	v := rdr.Read()

	ns, err := v.Info.Get("NS")
	c.Assert(err, IsNil)
	c.Assert(ns, Equals, 3)

	dp, err := v.Info.Get("DP")
	c.Assert(dp, Equals, 14)
	c.Assert(err, IsNil)

	nsf, err := v.Info.Get("NS")
	c.Assert(err, IsNil)
	c.Assert(nsf, Equals, int(3))

	dpf, err := v.Info.Get("DP")
	c.Assert(err, IsNil)
	c.Assert(dpf, Equals, int(14))

	hqs, err := v.Info.Get("AF")
	c.Assert(hqs, DeepEquals, []interface{}{0.5})
	c.Assert(err, IsNil)

	dpfs, err := v.Info.Get("DP")
	c.Assert(err, IsNil)
	c.Assert(dpfs, DeepEquals, 14)

}

func (s *VariantSuite) TestInfoField(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read()
	vstr := fmt.Sprintf("%s", v.Info)
	c.Assert(vstr, Equals, "NS=3;DP=14;AF=0.5;DB;H2")
}

func (s *VariantSuite) TestInfoMap(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read()

	vstr := fmt.Sprintf("%s", v)
	c.Assert(vstr, Equals, "20\t14370\trs6054257\tG\tA\t29.0\tPASS\tNS=3;DP=14;AF=0.5;DB;H2\tGT:GQ:DP:HQ\t0|0:48:1:51,51\t1|0:48:8:51,51\t1/1:43:5:.,.")

	v.Info.Add("asdf", 123)
	v.Info.Add("float", float32(123.2))
	has := v.Info.Contains("asdf")
	c.Assert(has, Equals, true)
	val, err := v.Info.Get("float")
	c.Assert(val, Equals, float32(123.2))
	c.Assert(err, IsNil)

	c.Assert(fmt.Sprintf("%s", v.Info), Equals, "NS=3;DP=14;AF=0.5;DB;H2;asdf=123;float=123.2")

	rdr.Clear()

}

func (s *VariantSuite) TestStartEnd(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v := rdr.Read()

	c.Assert(int(v.Start()), Equals, 14369)
	c.Assert(int(v.End()), Equals, 14370)
}

func (s *VariantSuite) TestIs(c *C) {
	rdr, err := NewReader(s.reader, false)
	c.Assert(err, IsNil)
	v1 := rdr.Read()
	v2 := rdr.Read()
	c.Assert(v1.Is(v1), Equals, true)

	c.Assert(v1.Is(v2), Equals, false)
	c.Assert(v2.Is(v1), Equals, false)
}

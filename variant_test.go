package vcfgo

import (
	"fmt"
	. "gopkg.in/check.v1"
	"io"
	"strings"
)

type VariantSuite struct {
	reader io.Reader
}

var _ = Suite(&VariantSuite{})

func (s *VariantSuite) SetUpTest(c *C) {
	s.reader = strings.NewReader(vcfStr)
}

func (s *VariantSuite) TestVariantGetInt(c *C) {
	rdr, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	v := rdr.Read()

	ns, ok := v.Info["NS"]
	c.Assert(ok, Equals, true)
	c.Assert(ns, Equals, 3)

	dp, ok := v.Info["DP"]
	c.Assert(dp, Equals, 14)
	c.Assert(ok, Equals, true)

	nsf, ok := v.Info["NS"]
	c.Assert(ok, Equals, true)
	c.Assert(nsf, Equals, int(3))

	dpf, ok := v.Info["DP"]
	c.Assert(ok, Equals, true)
	c.Assert(dpf, Equals, int(14))

	hqs, ok := v.Info["AF"]
	c.Assert(hqs, DeepEquals, []interface{}{0.5})
	c.Assert(ok, Equals, true)

	dpfs, ok := v.Info["DP"]
	c.Assert(ok, Equals, true)
	c.Assert(dpfs, DeepEquals, 14)

}

func (s *VariantSuite) TestInfoField(c *C) {
	rdr, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	v := rdr.Read()
	vstr := fmt.Sprintf("%s", v.Info)
	c.Assert(vstr, Equals, "NS=3;DP=14;AF=0.50;DB;H2")
}

func (s *VariantSuite) TestInfoMap(c *C) {
	rdr, err := NewReader(s.reader)
	c.Assert(err, IsNil)
	v := rdr.Read()

	vstr := fmt.Sprintf("%s", v)
	c.Assert(vstr, Equals, "20\t14370\trs6054257\tG\tA\t29.0\tPASS\tNS=3;DP=14;AF=0.50;DB;H2\tGT:GQ:DP:HQ\t0|0:48:1:51,51\t1|0:48:8:51,51\t1/1:43:5:.,.")

}

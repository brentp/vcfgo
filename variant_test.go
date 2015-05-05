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
	has := false
	for _, entry := range v.Info["__order"].([]string) {
		if entry == "asdf" {
			has = true
		}
	}
	c.Assert(has, Equals, true)
	c.Assert(v.Info["float"], Equals, float32(123.2))

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

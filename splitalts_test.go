package vcfgo

import . "gopkg.in/check.v1"

type SplitAltSuite struct {
}

var _ = Suite(&SplitAltSuite{})

func (s *SplitAltSuite) SetUpTest(c *C) {
}

func (s *SplitAltSuite) TestSplitA(c *C) {

	in := []interface{}{"AA", "BB", "CC"}
	out, err := splitA(in, 0, 3)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "AA")

	out, err = splitA(in, 0, 4)
	c.Assert(err, Not(IsNil))

	out, err = splitA(in, 2, 3)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "CC")

	in2 := []interface{}{}
	out, err = splitA(in2, 0, 3)
	c.Assert(err, Not(IsNil))

}

func (s *SplitAltSuite) TestSplitR(c *C) {

	in := []interface{}{"ref", "alt1", "alt2"}
	out, err := splitR(in, 0, 3)
	c.Assert(err, Not(IsNil))
	c.Assert(out, IsNil)

	out, err = splitR(in, 0, 2)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{"ref", "alt1"})

	out, err = splitR(in, 1, 2)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{"ref", "alt2"})

	out, err = splitR(in, 2, 2)
	c.Assert(err, Not(IsNil))
}

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

func (s *SplitAltSuite) TestSplitG(c *C) {

	in := []interface{}{281, 5, 9, 58, 0, 115, 338, 46, 116, 809}
	out, err := splitG(in, 1, 3)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{281, 5, 9})

	out, err = splitG(in, 2, 3)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{281, 58, 115})

	in = []interface{}{0, 30, 323, 31, 365, 483, 38, 291, 325, 567}
	out, err = splitG(in, 3, 3)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{0, 38, 567})

}
func (s *SplitAltSuite) TestSplitGT(c *C) {

	in := []interface{}{"1", "2"}
	out, err := splitGT(in, 3)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{".", "."})

	in = []interface{}{"1", "2"}
	out, err = splitGT(in, 2)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{".", "1"})

}

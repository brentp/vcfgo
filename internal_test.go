package vcfgo

import . "gopkg.in/check.v1"

type InternalSuite struct {
}

var _ = Suite(&InternalSuite{})

func (s *InternalSuite) TestFloatFmt(c *C) {

	fmt := fmtFloat32(-0)
	c.Assert(fmt, Equals, "0")

	fmt = fmtFloat32(-0.0)
	c.Assert(fmt, Equals, "0")

	fmt = fmtFloat64(-0.0)
	c.Assert(fmt, Equals, "0")

	fmt = fmtFloat64(-0)
	c.Assert(fmt, Equals, "0")

}

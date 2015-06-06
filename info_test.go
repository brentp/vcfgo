package vcfgo

import (
	"io"

	. "gopkg.in/check.v1"
)

type InfoSuite struct {
	reader io.Reader
}

var _ = Suite(&InfoSuite{})

func (s *InfoSuite) TestInfoGet(c *C) {
	i := InfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3")
	c.Assert(string(i.Get("as")), Equals, "22")
	c.Assert(string(i.Get("asdf")), Equals, "123")
	c.Assert(string(i.Get("ddd")), Equals, "ddd")
	c.Assert(string(i.Get("dddd")), Equals, "dddd")
	c.Assert(string(i.Get("other")), Equals, "as")
	c.Assert(string(i.Get("FLAG1")), Equals, "FLAG1")
	c.Assert(string(i.Get("FLAG")), Equals, "FLAG")
	c.Assert(string(i.Get("FLAG2")), Equals, "FLAG2")
	c.Assert(string(i.Get("FLAG3")), Equals, "FLAG3")
	c.Assert(string(i.Get("FLAG4")), Equals, "")
	c.Assert(string(i.Get("dd")), Equals, "")
	c.Assert(string(i.Get("")), Equals, "")
	c.Assert(string(i.Get("ddddd")), Equals, "")
	c.Assert(string(i.Get("FLAG33")), Equals, "")
}

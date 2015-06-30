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
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3", nil)
	c.Assert(string(i.SGet("as")), Equals, "22")
	c.Assert(string(i.SGet("asdf")), Equals, "123")
	c.Assert(string(i.SGet("ddd")), Equals, "ddd")
	c.Assert(string(i.SGet("dddd")), Equals, "dddd")
	c.Assert(string(i.SGet("other")), Equals, "as")
	c.Assert(string(i.SGet("FLAG1")), Equals, "FLAG1")
	c.Assert(string(i.SGet("FLAG")), Equals, "FLAG")
	c.Assert(string(i.SGet("FLAG2")), Equals, "FLAG2")
	c.Assert(string(i.SGet("FLAG3")), Equals, "FLAG3")
	c.Assert(string(i.SGet("FLAG4")), Equals, "")
	c.Assert(string(i.SGet("dd")), Equals, "")
	c.Assert(string(i.SGet("")), Equals, "")
	c.Assert(string(i.SGet("ddddd")), Equals, "")
	c.Assert(string(i.SGet("FLAG33")), Equals, "")
	i = NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3;asst=33,44", nil)
	c.Assert(string(i.SGet("t")), Equals, "")
}

func (s *InfoSuite) TestInfoSet(c *C) {
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3", nil)
	i.Set("as", "23")
	c.Assert(string(i.SGet("as")), Equals, "23")
	i.Set("asst", 24)
	c.Assert(string(i.SGet("asst")), Equals, "24")
	i.Set("t", 33)
	a := string(i.SGet("t"))
	c.Assert(a, Equals, "33")
	i.Set("t", []interface{}{93, 44, 55, 66})
	c.Assert(string(i.SGet("t")), Equals, "93,44,55,66")
	i.Set("tt", "asdf")
	c.Assert(string(i.SGet("t")), Equals, "93,44,55,66")
	c.Assert(string(i.SGet("tt")), Equals, "asdf")
	i.Set("ttt", "xxx")
	c.Assert(string(i.SGet("t")), Equals, "93,44,55,66")
	c.Assert(string(i.SGet("tt")), Equals, "asdf")
	c.Assert(string(i.SGet("ttt")), Equals, "xxx")

	i.Set("zzz", "zzz")
	i.Set("zz", "zz")
	i.Set("z", "z")
	c.Assert(string(i.SGet("zzz")), Equals, "zzz")
	c.Assert(string(i.SGet("zz")), Equals, "zz")
	c.Assert(string(i.SGet("z")), Equals, "z")

}

func (s *InfoSuite) TestInfoKeys(c *C) {
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3", nil)
	c.Assert(i.Keys(), DeepEquals, []string{"asdf", "FLAG1", "ddd", "FLAG", "dddd", "as", "FLAG2", "other", "FLAG3"})
}

func (s *InfoSuite) TestInfoDelete(c *C) {
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3", nil)
	i.Delete("asdf")
	c.Assert(i.String(), Equals, "FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3")
	i.Delete("other")
	c.Assert(i.String(), Equals, "FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;FLAG3")

	i.Delete("FLAG1")
	c.Assert(i.String(), Equals, "ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;FLAG3")

	i.Delete("FLAG")
	c.Assert(i.String(), Equals, "ddd=ddd;dddd=dddd;as=22;FLAG2;FLAG3")

	i.Delete("ddd")
	c.Assert(i.String(), Equals, "dddd=dddd;as=22;FLAG2;FLAG3")

	i.Delete("FLAG3")
	c.Assert(i.String(), Equals, "dddd=dddd;as=22;FLAG2")

	i.Delete("FLAG2")
	c.Assert(i.String(), Equals, "dddd=dddd;as=22")

	i.Delete("dddd")
	c.Assert(i.String(), Equals, "as=22")

	i.Delete("as")
	c.Assert(i.String(), Equals, "")

}

func (s *InfoSuite) TestInfoFlag(c *C) {
	i := NewInfoByte("AAA;asdf=123;FLAG1;ddd=123", nil)
	i.Set("ggg", true)
	c.Assert(i.String(), Equals, "AAA;asdf=123;FLAG1;ddd=123;ggg")

	i.Set("ggg", false)
	c.Assert(i.String(), Equals, "AAA;asdf=123;FLAG1;ddd=123")

	i.Set("ggg", true)
	c.Assert(i.String(), Equals, "AAA;asdf=123;FLAG1;ddd=123;ggg")

	i.Set("gga", true)
	c.Assert(i.String(), Equals, "AAA;asdf=123;FLAG1;ddd=123;ggg;gga")

	i.Set("AAA", true)
	c.Assert(i.String(), Equals, "AAA;asdf=123;FLAG1;ddd=123;ggg;gga")

	// NOTE: setting to false removes
	i.Set("AAA", false)
	c.Assert(i.String(), Equals, "asdf=123;FLAG1;ddd=123;ggg;gga")
}

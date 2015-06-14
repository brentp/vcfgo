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
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3")
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
	i = NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3;asst=33,44")
	c.Assert(string(i.Get("t")), Equals, "")
}

func (s *InfoSuite) TestInfoSet(c *C) {
	i := NewInfoByte("asdf=123;FLAG1;ddd=ddd;FLAG;dddd=dddd;as=22;FLAG2;other=as;FLAG3")
	i.Set("as", "23")
	c.Assert(string(i.Get("as")), Equals, "23")
	i.Set("asst", 24)
	c.Assert(string(i.Get("asst")), Equals, "24")
	i.Set("t", 33)
	a := string(i.Get("t"))
	c.Assert(a, Equals, "33")
	i.Set("t", []interface{}{93, 44, 55, 66})
	c.Assert(string(i.Get("t")), Equals, "93,44,55,66")
	i.Set("tt", "asdf")
	c.Assert(string(i.Get("t")), Equals, "93,44,55,66")
	c.Assert(string(i.Get("tt")), Equals, "asdf")
	i.Set("ttt", "xxx")
	c.Assert(string(i.Get("t")), Equals, "93,44,55,66")
	c.Assert(string(i.Get("tt")), Equals, "asdf")
	c.Assert(string(i.Get("ttt")), Equals, "xxx")

	i.Set("zzz", "zzz")
	i.Set("zz", "zz")
	i.Set("z", "z")
	c.Assert(string(i.Get("zzz")), Equals, "zzz")
	c.Assert(string(i.Get("zz")), Equals, "zz")
	c.Assert(string(i.Get("z")), Equals, "z")

}

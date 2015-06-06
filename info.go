package vcfgo

import "bytes"

type InfoByte []byte

func (i *InfoByte) Set(key string, value interface{}) {

}

func (i InfoByte) Get(key string) []byte {
	var sub []byte
	if key == "" {
		return sub
	}
	bkey := []byte(key)
	pos := 0
	for {
		ipos := bytes.Index(i[pos:], bkey)
		if ipos == -1 {
			return []byte{}
		}
		pos += ipos
		eq := pos + bytes.IndexByte(i[pos:], byte('='))
		// at end of field and we found an Flag
		var semi int
		if eq == -1 {
			return i[pos:]
		} else if eq-pos != len(bkey) {
			// found a longer key with same prefix.
			semi = bytes.IndexByte(i[pos:], byte(';'))
			// flag field
			if semi == -1 {
				semi = len(i)
			} else {
				semi += pos
			}
			if semi-pos == len(bkey) {
				return i[pos:semi]
			}
			pos = semi + 1
			continue
		} else {
			semi = bytes.IndexByte(i[pos:], byte(';'))
		}
		if semi > -1 && eq > pos+semi {
			// should be a flag.
			return i[pos : pos+semi]
		}

		// not at end of info field
		if semi != -1 {
			semi += pos
		} else {
			semi = len(i)
		}
		sub = i[eq+1 : semi]
		break
	}
	return sub
}

func (i InfoByte) String() string {
	return string(i)
}

package vcfgo

import "bytes"

type InfoByte struct {
	info   []byte
	parsed map[string]interface{}
}

func NewInfoByte(info string) InfoByte {
	return InfoByte{info: []byte(info), parsed: make(map[string]interface{})}
}

// return the start and end positions of the value.
// for flag the value is the flag.
func getpositions(info []byte, key string) (start, end int) {
	bkey := []byte(key)
	pos := 0
	for {
		ipos := bytes.Index(info[pos:], bkey)
		if ipos == -1 {
			return -1, -1
		}
		pos += ipos
		eq := pos + bytes.IndexByte(info[pos:], byte('='))
		// at end of field and we found an Flag
		var semi int
		if eq == -1 {
			return pos, len(info)
		} else if eq-pos != len(bkey) {
			// found a longer key with same prefix.
			semi = bytes.IndexByte(info[pos:], byte(';'))
			// flag field
			if semi == -1 {
				semi = len(info)
			} else {
				semi += pos
			}
			if semi-pos == len(bkey) {
				return pos, semi
			}
			pos = semi + 1
			continue
		} else {
			semi = bytes.IndexByte(info[pos:], byte(';'))
		}
		if semi > -1 && eq > pos+semi {
			// should be a flag.
			return pos, pos + semi
		}

		// not at end of info field
		if semi != -1 {
			semi += pos
		} else {
			semi = len(info)
		}
		return eq + 1, semi
	}
}

func (i InfoByte) Contains(key string) bool {
	// short-circuit common case.
	if bytes.Contains(i.info, []byte(key+"=")) {
		return false
	}
	s, _ := getpositions(i.info, key)
	return s == -1
}

// TODO: attach to header so we can get type.
func (i InfoByte) Get(key string) []byte {
	if val, ok := i.parsed[key]; ok {
		return val.([]byte)
	}
	var sub []byte
	if key == "" {
		return sub
	}
	start, end := getpositions(i.info, key)
	if start == -1 {
		return sub
	}
	val := i.info[start:end]
	i.parsed[key] = val
	return val
}

func (i InfoByte) String() string {
	return string(i.info)
}

func (i *InfoByte) Set(key string, value interface{}) {

}

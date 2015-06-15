package vcfgo

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type InfoByte struct {
	info   []byte
	header *Header
}

func NewInfoByte(info string, h *Header) *InfoByte {
	return &InfoByte{info: []byte(info), header: h}
}

// return the start and end positions of the value.
// for flag the value is the flag.
func getpositions(info []byte, key string) (start, end int) {
	bkey := []byte(key)
	pos := 0
	for {
		if pos >= len(info) {
			return -1, -1
		}
		ipos := bytes.Index(info[pos:], bkey)
		if ipos == -1 {
			return -1, -1
		}
		pos += ipos
		if pos != 0 && info[pos-1] != ';' {
			pos += 1
			continue
		}
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
				return pos, semi - 1
			}
			pos = semi + 1
			continue
		} else {
			semi = bytes.IndexByte(info[pos+1:], byte(';'))
		}
		if semi > -1 && eq > pos+semi {
			// should be a flag.
			return pos, pos + semi
		}

		// not at end of info field
		if semi != -1 {
			semi += pos
		}
		return eq + 1, semi
	}
}

func (i InfoByte) Contains(key string) bool {
	// short-circuit common case.
	if !bytes.Contains(i.info, []byte(key+"=")) {
		return false
	}
	s, _ := getpositions(i.info, key)
	return s != -1
}

func (i InfoByte) Keys() []string {
	sp := bytes.Split(i.info, []byte{';'})
	keys := make([]string, 0, len(sp))
	for _, pair := range sp {
		key := bytes.SplitN(pair, []byte{'='}, 2)[0]
		keys = append(keys, string(key))
	}
	return keys

}

func (i *InfoByte) Delete(key string) {
	s, e := getpositions(i.info, key)
	if s == -1 {
		return
	}
	// check if it's a flag
	if s != 0 && i.info[s-1] != ';' {
		s -= (len(key) + 1)
	}
	if s < 0 {
		s = 0
	}
	if e == -1 {
		e = len(i.info)
	} else {
		e += 2
	}
	if s == 0 && e == len(i.info) {
		i.info = i.info[:0]
	} else if e < len(i.info) {
		i.info = append(i.info[:s], i.info[e:]...)
	} else {
		i.info = i.info[:s-1]
	}
}

func ItoS(k string, v interface{}) string {
	if b, ok := v.(bool); ok && b {
		return k
	} else {
		switch v.(type) {
		case float32:
			return fmtFloat32(v.(float32))
		case float64:
			return fmtFloat64(v.(float64))
		case int:
			return fmt.Sprintf("%d", v.(int))
		case uint32:
			return fmt.Sprintf("%d", v.(uint32))
		case []interface{}:
			vals := v.([]interface{})
			svals := make([]string, len(vals))
			switch vals[0].(type) {
			case float64:
				for i, val := range vals {
					svals[i] = fmtFloat64(val.(float64))
				}
			case float32:
				for i, val := range vals {
					svals[i] = fmtFloat32(val.(float32))
				}
			case int:
				for i, val := range vals {
					svals[i] = strconv.Itoa(val.(int))
				}
			default:
				for i, val := range vals {
					svals[i] = fmt.Sprintf("%v", val)
				}
			}
			return strings.Join(svals, ",")

		default:
			return v.(string)
		}
	}
}

func (i InfoByte) SGet(key string) []byte {
	var sub []byte
	if key == "" || len(i.info) == 1 {
		return sub
	}
	start, end := getpositions(i.info, key)
	if start == -1 {
		return sub
	}
	if end == -1 {
		end = len(i.info) - 1
	}
	val := i.info[start : end+1]
	return val
}

// Get a value from the bytes typed according to the header.
func (i InfoByte) Get(key string) (interface{}, error) {
	v := string(i.SGet(key))
	skey := string(key)
	hi, ok := i.header.Infos[key]
	if !ok {
		err := fmt.Errorf("Info Error: %s not found in header", skey)
		// flag
		if skey == v {
			return true, err
		}
		return v, err
	}

	if len(v) == 0 {
		var err error
		var val interface{} = nil
		if hi.Type != "Flag" {
			err = fmt.Errorf("Info Error: %s not found in row", skey)
		} else {
			val = false
		}
		return val, err
	}
	if v == key {
		var err error
		if hi.Type != "Flag" {
			err = fmt.Errorf("Info Error: flag field (%s) should be specified as such in the header", skey)
		}
		return true, err
	}

	var err error
	var iv interface{}

	switch hi.Number {

	case "0":
		if hi.Type != "Flag" {
			err = fmt.Errorf("Info Error: flag field (%s) should have Number=0", skey)
		}
		return true, err

	case "1":
		return parseOne(skey, v, hi.Type)

	case "R", "A", "G", "2", "3", ".":
		vals := strings.Split(v, ",")
		var vi interface{} = make([]interface{}, len(vals))
		for j, val := range vals {
			iv, err = parseOne(skey, val, hi.Type)
			vi.([]interface{})[j] = iv
		}
		return vi, err

	default:
		vals := strings.Split(v, ",")
		var vi interface{} = make([]interface{}, len(vals))
		if _, err := strconv.Atoi(hi.Number); err == nil {
			for j, val := range vals {
				iv, err = parseOne(skey, val, hi.Type)
				vi.([]interface{})[j] = iv
			}
		} else {
			panic(fmt.Sprintf("found Number=%s", hi.Number))
		}
		return vi, err

	}
}

func (i InfoByte) String() string {
	return string(i.info)
}

func (i *InfoByte) UpdateHeader(key string, value interface{}) {
	if i.header != nil {
		switch value.(type) {
		case bool:
			i.header.Infos[key] = &Info{Id: key, Description: key, Number: "0", Type: "Flag"}
		case string:
			i.header.Infos[key] = &Info{Id: key, Description: key, Number: "1", Type: "Character"}
		case int, int32, int64, uint32, uint64:
			i.header.Infos[key] = &Info{Id: key, Description: key, Number: "1", Type: "Integer"}
		case float32, float64:
			i.header.Infos[key] = &Info{Id: key, Description: key, Number: "1", Type: "Float"}
		case []interface{}:
			v := value.([]interface{})[0]
			i.UpdateHeader(key, v)
		}
	}
}

func (i *InfoByte) Set(key string, value interface{}) {
	if len(i.info) == 0 {
		i.info = []byte(fmt.Sprintf("%s=%s", key, ItoS(key, value)))
		return
	}
	s, e := getpositions(i.info, key)
	if s == -1 || s == len(i.info) {
		slug := []byte(fmt.Sprintf(";%s=%s", key, ItoS(key, value)))
		i.info = append(i.info, slug...)
		i.UpdateHeader(key, value)
		return
	}
	slug := []byte(ItoS(key, value))
	if e == -1 {
		i.info = append(i.info[:s], slug...)
	} else {
		i.info = append(i.info[:s], append(slug, i.info[e+1:]...)...)
	}
}

func (i *InfoByte) Add(key string, value interface{}) {
	i.Set(key, value)
}

package torrent

import (
	"strconv"
	"strings"
)

//Bencode represents bencoded data
type Bencode interface {
	ToString() string
}

type bInt struct {
	value int
}

func (ben *bInt) ToString() string {
	return strconv.Itoa(ben.value)
}

type bStr struct {
	value string
}

func (ben *bStr) ToString() string {
	return ben.value
}

type bList struct {
	value []Bencode
}

func (ben *bList) ToString() string {
	var valuesStr []string
	for _, v := range ben.value {
		strV := v.ToString()
		valuesStr = append(valuesStr, strV)
	}
	return "[" + strings.Join(valuesStr, ",") + "]"
}

type bDict struct {
	value map[bStr]Bencode
}

func (ben *bDict) get(key string) Bencode {
	return ben.value[bStr{key}]
}

func (ben *bDict) ToString() string {
	lineBreak := "\n"
	var pairs []string
	for k, v := range ben.value {
		key := k.value

		var value string
		if key == "pieces" {
			value = "..."
		} else {
			value = v.ToString()
		}
		pairs = append(pairs, lineBreak+key+"->"+value)
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

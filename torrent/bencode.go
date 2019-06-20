package torrent

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

//Bencode represents bencoded data
type Bencode interface {
	String() string
	PrettyString() string
}

type bInt struct {
	value int
}

func (ben *bInt) PrettyString() string {
	return strconv.Itoa(ben.value)
}

func (ben *bInt) String() string {
	return fmt.Sprint("i", ben.value, "e")
}

type bStr struct {
	value string
}

func (ben *bStr) PrettyString() string {
	return ben.value
}

func (ben *bStr) String() string {
	return fmt.Sprint(len(ben.value), ":", ben.value)
}

type bList struct {
	value []Bencode
}

func (ben *bList) PrettyString() string {
	var valuesStr []string
	for _, v := range ben.value {
		strV := v.PrettyString()
		valuesStr = append(valuesStr, strV)
	}
	return "[" + strings.Join(valuesStr, ",") + "]"
}

func (ben *bList) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("l")
	for _, v := range ben.value {
		buffer.WriteString(v.String())
	}
	buffer.WriteString("e")
	return buffer.String()
}

type bDict struct {
	//	"github.com/emirpasic/gods/maps/linkedhashmap"
	value map[bStr]Bencode
}

func (ben *bDict) get(key string) Bencode {
	return ben.value[bStr{key}]
}

func (ben *bDict) PrettyString() string {
	lineBreak := "\n"
	var pairs []string
	for k, v := range ben.value {
		key := k.value

		var value string
		if key == "pieces" {
			value = "..."
		} else {
			value = v.PrettyString()
		}
		pairs = append(pairs, lineBreak+key+"->"+value)
	}
	return "{" + strings.Join(pairs, ", ") + "}"
}

func (ben *bDict) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("d")
	for k, v := range ben.value {
		buffer.WriteString(k.String())
		buffer.WriteString(v.String())
	}
	buffer.WriteString("e")
	return buffer.String()
}

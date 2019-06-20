package torrent

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
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

func newList() *bList {
	return &bList{}
}

func (ben *bList) append(elem Bencode) {
	ben.value = append(ben.value, elem)
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
	value *linkedhashmap.Map
}

func (ben *bDict) get(key string) Bencode {
	v, ok := ben.value.Get(bStr{key})
	if !ok {
		return nil
	}
	return v.(Bencode)
}

func newDict() *bDict {
	return &bDict{linkedhashmap.New()}
}

func (ben *bDict) put(key bStr, value Bencode) {
	ben.value.Put(key, value)
}

type dictIter struct {
	it linkedhashmap.Iterator
}

func (ben *bDict) iter() dictIter {
	it := ben.value.Iterator()
	return dictIter{it}
}

func (dI *dictIter) next() bool {
	return dI.it.Next()
}

func (dI *dictIter) key() bStr {
	return dI.it.Key().(bStr)
}

func (dI *dictIter) value() Bencode {
	return dI.it.Value().(Bencode)
}

func (ben *bDict) PrettyString() string {
	lineBreak := "\n"
	var pairs []string
	it := ben.iter()
	for it.next() {
		k, v := it.key(), it.value()
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
	it := ben.iter()
	for it.next() {
		key, value := it.key(), it.value()
		buffer.WriteString(key.String())
		buffer.WriteString(value.String())
	}
	buffer.WriteString("e")
	return buffer.String()
}

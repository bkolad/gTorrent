package torrent

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//ParseError error for missformated bencoded data
type ParseError struct {
	benType string
	pos     int
	error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("problem with parsing bencoded %s : %s at %d", e.benType, e.Error(), e.pos)
}

func parseError(benType string, err error, pos int) error {
	return &ParseError{benType, pos, err}
}

//Parser parses bencoded data to Bencode
type Parser interface {
	Parse() (Bencode, error)
}

type parser struct {
	data  string
	index int
}

//NewParser returns default bencode parser
func NewParser(str string) Parser {
	return &parser{str, 0}
}

//example: i87e
func (p *parser) parseInt() (Bencode, error) {
	len := strings.Index(p.data[p.index:], "e")
	value, err := strconv.Atoi(p.data[p.index+1 : p.index+len])
	if err != nil {
		return nil, parseError("int", err, p.index)
	}
	p.index += len + 1
	return &bInt{value}, nil
}

//example: 3:dog
func (p *parser) parseStr() (Bencode, error) {
	strs := strings.SplitN(p.data[p.index:], ":", 2)
	lenStr := strs[0]
	length, err := strconv.Atoi(lenStr)
	if err != nil {
		return nil, parseError("string", err, p.index)
	}
	if len(strs) != 2 {
		err = parseError("string", errors.New("bencoded string should have format: len:str"), p.index)
		return nil, err
	}

	value := strs[1][0:length]
	p.index += len(lenStr) + 1 + length
	return &bStr{value}, nil
}

//example: l...e -> li22ei78e3:dogi21ee -> [22,78,dog,21]
func (p *parser) parseList() (Bencode, error) {
	var list []Bencode
	p.index++
	for p.index < len(p.data) {
		if p.data[p.index] == 'e' {
			p.index++
			break
		}
		bencode, err := p.Parse()
		if err != nil {
			return nil, err
		}

		list = append(list, bencode)
	}
	return &bList{list}, nil
}

//example: dstr:ben str:ben -> d1:xi22e6:animal:3:doge -> {x:22, animal:dog}
func (p *parser) parseDict() (Bencode, error) {
	benMap := make(map[bStr]Bencode)
	p.index++
	for p.index < len(p.data) {
		if p.data[p.index] == 'e' {
			p.index++
			break
		}

		bencodeKey, err := p.Parse()
		if err != nil {
			return nil, err
		}

		bencodeKeyStr := *bencodeKey.(*bStr)
		bencodeValue, err := p.Parse()
		if err != nil {
			return nil, err
		}
		benMap[bencodeKeyStr] = bencodeValue
	}
	return &bDict{benMap}, nil
}

func (p *parser) Parse() (Bencode, error) {
	for p.index < len(p.data) {
		c := p.data[p.index]
		switch {
		case c == 'i':
			return p.parseInt()

		case c == 'l':
			return p.parseList()

		case c == 'd':
			return p.parseDict()

		case c >= '0' && c <= '9':
			return p.parseStr()

		default:
			msg := fmt.Sprintf("can't parse %c at %d", c, p.index)
			return nil, errors.New(msg)
		}
	}
	return nil, nil
}

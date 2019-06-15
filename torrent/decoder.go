package torrent

import (
	"errors"
)

type Info struct {
	announce     string
	announceList [][]string
	pieceSize    int
	length       int
	name         string
}

//Decoder decodeds tracker data
type Decoder interface {
	Decode() (*Info, error)
}

type defaultDec struct {
	str string
}

//NewDecoder returns default decoder, holds the entire torrent file in memory.
func NewDecoder(str string) Decoder {
	return &defaultDec{str}
}

func (dec *defaultDec) Decode() (*Info, error) {
	p := NewParser(dec.str)
	ben, err := p.Parse()
	if err != nil {
		return nil, err
	}

	dict, ok := ben.(*bDict)
	if !ok {
		return nil, wrongTypeError("Torrent content ", "dictionary")
	}

	//==== announce
	announce, err := strValue(dict, "announce")
	if err != nil {
		return nil, err
	}

	//==== announceList
	announceList, err := announceList(dict)
	if err != nil {
		return nil, err
	}

	//==== info
	info, err := fromDict(dict, "info")
	infoDict, ok := info.(*bDict)
	if !ok {
		return nil, wrongTypeError("info", "dictionary")
	}

	if err != nil {
		return nil, err
	}

	//==== pieceLength
	pieceLength, err := intValue(infoDict, "piece length")
	if err != nil {
		return nil, err
	}

	//==== length
	length, err := intValue(infoDict, "length")
	if err != nil {
		return nil, err
	}

	//==== name
	name, err := strValue(infoDict, "name")
	if err != nil {
		return nil, err
	}

	return &Info{
		announce,
		announceList,
		pieceLength,
		length,
		name,
	}, err
}

func announceList(bencs *bDict) ([][]string, error) {
	benValue, err := fromDict(bencs, "announce-list")
	if err != nil {
		return nil, err
	}

	benList, ok := benValue.(*bList)
	if !ok {
		return nil, wrongTypeError("announce-list entry", "list of lists")
	}

	var list [][]string
	for _, ls := range benList.value {
		l, ok := ls.(*bList)
		if !ok {
			return nil, wrongTypeError("announce-list entry", "list")
		}
		var internalList []string
		for _, s := range l.value {
			internalList = append(internalList, s.ToString())
		}
		list = append(list, internalList)
	}
	return list, nil
}

func fromDict(dict *bDict, key string) (Bencode, error) {
	value := dict.get(key)
	if value == nil {
		return nil, errors.New(key + " is missing in the dictionary")
	}
	return dict.get(key), nil
}

func intValue(dict *bDict, key string) (int, error) {
	benLength, err := fromDict(dict, key)
	if err != nil {
		return 0, err
	}

	length, ok := benLength.(*bInt)
	if !ok {
		return 0, wrongTypeError(key, "int")
	}

	return length.value, nil
}

func strValue(dict *bDict, key string) (string, error) {
	benName, err := fromDict(dict, key)
	if err != nil {
		return "", err
	}

	name, ok := benName.(*bStr)
	if !ok {
		return "", wrongTypeError(key, "string")
	}

	return name.value, nil
}

func wrongTypeError(str string, t string) error {
	return errors.New("wrong type, " + "-" + str + "- has to be" + t)
}

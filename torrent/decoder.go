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

//Decoder decodeds bencoded tracker related data
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
		return nil, errors.New("Torret content has to be bencoded dictionary")
	}

	//====
	announce, err := strValue(dict, "announce")
	if err != nil {
		return nil, err
	}

	//====
	announceList, err := getAnnounceList(dict)
	if err != nil {
		return nil, err
	}

	//====
	info, err := getFromDict(dict, "info")
	infoDict, ok := info.(*bDict)
	if !ok {
		return nil, errors.New("-info- has to be bencoded dictionary")
	}

	if err != nil {
		return nil, err
	}

	//====
	pieceLength, err := intValue(infoDict, "piece length")
	if err != nil {
		return nil, err
	}

	//====
	length, err := intValue(infoDict, "length")
	if err != nil {
		return nil, err
	}

	//====
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

func getAnnounceList(bencs *bDict) ([][]string, error) {
	benValue, err := getFromDict(bencs, "announce-list")
	if err != nil {
		return nil, err
	}

	benList, ok := benValue.(*bList)
	if !ok {
		return nil, errors.New("announce-list has to be bencoded list of lists")
	}

	var list [][]string
	for _, ls := range benList.value {
		l, ok := ls.(*bList)
		if !ok {
			return nil, errors.New("announce-list entry has to be a list")
		}
		var internalList []string
		for _, s := range l.value {
			internalList = append(internalList, s.ToString())
		}
		list = append(list, internalList)
	}
	return list, nil
}

func getFromDict(dict *bDict, key string) (Bencode, error) {
	value := dict.get(key)
	if value == nil {
		return nil, errors.New(key + " is missing in the dictionary")
	}
	return dict.get(key), nil
}

func intValue(dict *bDict, key string) (int, error) {
	benLength, err := getFromDict(dict, key)
	if err != nil {
		return 0, err
	}

	length, ok := benLength.(*bInt)
	if !ok {
		return 0, errors.New("-" + key + "- has to be int")
	}

	return length.value, nil
}

func strValue(dict *bDict, key string) (string, error) {
	benName, err := getFromDict(dict, key)
	if err != nil {
		return "", err
	}

	name, ok := benName.(*bStr)
	if !ok {
		return "", errors.New("-" + key + "- has to be string")
	}

	return name.value, nil
}

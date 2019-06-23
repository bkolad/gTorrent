package torrent

import (
	"crypto/sha1"
	"errors"
	"io"
)

type file struct {
	length int
	path   []string
}

type Info struct {
	announce     string
	announceList [][]string
	pieceSize    int
	length       int
	name         string
	files        []file
	pieceHashes  [][]byte
	InfoHash     []byte
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
	pieceLength, _, err := intValue(infoDict, "piece length")
	if err != nil {
		return nil, err
	}

	//==== length
	length, isSingleFile, err := intValue(infoDict, "length")
	//if length is absent then we have multifile torrent
	if isSingleFile && err != nil {
		return nil, err
	}

	//==== name
	name, err := strValue(infoDict, "name")
	if err != nil {
		return nil, err
	}

	//==== files
	files, err := files(infoDict)
	if err != nil {
		return nil, err
	}

	if files == nil && !isSingleFile {
		return nil, errors.New("No files to download in the torrent file")
	}

	pieces, err := strValue(infoDict, "pieces")
	if err != nil {
		return nil, err
	}

	pieceHash, err := pieceHashes(pieces)
	if err != nil {
		return nil, err
	}

	infoDictStr := infoDict.String()
	h := sha1.New()
	io.WriteString(h, infoDictStr)
	infoHash := h.Sum(nil)
	return &Info{
		announce,
		announceList,
		pieceLength,
		length,
		name,
		files,
		pieceHash,
		infoHash,
	}, err
}

func announceList(bencs *bDict) ([][]string, error) {
	benValue, err := fromDict(bencs, "announce-list")
	if err != nil {
		return nil, nil
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
			internalList = append(internalList, s.PrettyString())
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
	return value, nil
}

func intValue(dict *bDict, key string) (int, bool, error) {
	benLength, err := fromDict(dict, key)
	if err != nil {
		return 0, false, err
	}

	length, ok := benLength.(*bInt)
	if !ok {
		return 0, true, wrongTypeError(key, "int")
	}

	return length.value, true, nil
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

func files(infoDict *bDict) ([]file, error) {
	benFiles, err := fromDict(infoDict, "files")
	if err != nil {
		return nil, nil
	}
	benList, ok := benFiles.(*bList)
	if !ok {
		return nil, wrongTypeError("files", "list")
	}

	var files []file
	for _, v := range benList.value {
		d := v.(*bDict)
		length, _, err := intValue(d, "length")
		if err != nil {
			return nil, err
		}

		benPath, err := fromDict(d, "path")
		if err != nil {
			return nil, err
		}

		benPathList, ok := benPath.(*bList)
		if !ok {
			return nil, nil
		}
		var path []string
		for _, v := range benPathList.value {
			path = append(path, v.PrettyString())
		}
		files = append(files, file{length, path})
	}
	return files, nil
}

func pieceHashes(pieces string) ([][]byte, error) {
	len := len(pieces)
	if len%20 != 0 {
		return nil, errors.New("piece hash has to be 20 bytes long")
	}

	pieceHashes := make([][]byte, 0)
	for i := 0; i <= len-20; i += 20 {
		pieceHashes = append(pieceHashes, []byte(pieces[i:i+20]))
	}

	return pieceHashes, nil
}

func wrongTypeError(str string, t string) error {
	return errors.New("wrong type, " + "-" + str + "- has to be" + t)
}

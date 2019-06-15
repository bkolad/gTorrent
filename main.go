package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bkolad/gTorrent/torrent"
)

func main() {
	data, err := ioutil.ReadFile("testData/ubuntu-15.10-desktop-amd64.__iso.torrent")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	dec := torrent.NewDecoder(string(data))
	info, err := dec.Decode()

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)
}

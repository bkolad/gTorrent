package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bkolad/gTorrent/torrent"
)

const torrentContent = "d8:announce39:http://torrent.ubuntu.com:6969/announce" +
	"13:announce-listll39:http://torrent.ubuntu.com:6969/announcee" +
	"l44:http://ipv6.torrent.ubuntu.com:6969/announceee" +
	"7:comment29:Ubuntu CD releases.ubuntu.com13:creation datei1445507299e" +
	"4:info" +
	"d6:lengthi1e" +
	"4:name30:ubuntu-15.10-desktop-amd64.iso" +
	"12:piece lengthi524288e" +
	"6:pieces20:aaaaaaaaaaaaaaaaaaaaee"

func main() {
	data, err := ioutil.ReadFile("testData/Fedora-Live-Cinnamon-i686-23.torrent")
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

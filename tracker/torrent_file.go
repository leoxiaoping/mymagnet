package tracker

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
)

type torrentFile struct {
	length int64    //size of file in bytes
	path   []string //filename
}

type TorrentInfo struct {
	announce     string   // tracker url
	name         string   // filename or dirname depending on length of files
	pieceLength  int64    // size of each piece in bytes
	pieces       []string // checksums for each piece
	numpieces    int
	files        []torrentFile // info for each file
	numfiles     int64         // not part of file, but helpful
	info_hash    string        // sha1 of info dict
	client_id    string        // randomly generated 20 bytes
	our_pieces   *Pieces
	total_length int64
	client_chan  chan Message
}

func ReadTorrentFile(path string, c chan Message) (*TorrentInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	t := ParseTorrentInfo(b)
	t.client_chan = c
	t.our_pieces = CreateNewPieces(t.numpieces, t)
	return t, nil
}

func FindInfoHash(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	t := ParseTorrentInfo(b)
	return t.info_hash
}

func (t *TorrentInfo) add_info_hash(info bItem) {
	h := sha1.New()
	h.Write(info.raw)
	t.info_hash = fmt.Sprintf("%x", h.Sum(nil))
	//t.info_hash = "132131321312313"
}

func (t *TorrentInfo) generate_client_id() {
	t.client_id = string(RandomBytes(20))
}

func ParseTorrentInfo(b []byte) *TorrentInfo {
	bi, _ := Bdecode(b)
	t := new(TorrentInfo)
	t.announce = bi.d["announce"].s

	info := bi.d["info"].d

	t.add_info_hash(bi.d["info"])
	t.generate_client_id()

	t.name = info["name"].s
	t.pieceLength = info["piece length"].i
	for i := 0; i < len(info["pieces"].s); i += 20 {
		t.pieces = append(t.pieces, info["pieces"].s[i:i+20])
	}

	t.numpieces = len(t.pieces)

	t.total_length = 0

	if info["length"].i > 0 {
		f := new(torrentFile)
		f.length = info["length"].i
		f.path = []string{t.name}
		t.files = append(t.files, *f)
		t.numfiles = 1
		t.total_length += f.length
	} else {
		t.numfiles = 0
		for _, v := range info["files"].l {
			f := new(torrentFile)
			f.length = v.d["length"].i
			paths := make([]string, 0)
			for _, t := range v.d["path"].l {
				paths = append(paths, t.s)
			}
			f.path = paths
			t.files = append(t.files, *f)
			t.numfiles++
			t.total_length += f.length
		}
	}

	return t
}

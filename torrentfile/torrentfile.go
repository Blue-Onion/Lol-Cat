package torrentfile

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

const peerSize = 6

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Comment  string      `bencode:"comment"`
	Info     bencodeInfo `bencode:"info"`
}
type Torrentfile struct {
	Announce    string
	Comment     string
	PieceLength int
	Length      int
	Name        string
	InfoHashes  [20]byte
	PieceHashes [][20]byte
}
type Peer struct {
	IP   net.IP
	Port uint16
}

func Open(r io.Reader) (*bencodeTorrent, error) {
	bto := bencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		return nil, err
	}
	return &bto, nil
}
func (t *Torrentfile) buildTracker(peerId [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	params := url.Values{
		"info_hash":  []string{string(t.InfoHashes[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
func (bto bencodeTorrent) toTorrentFile() (*Torrentfile, error) {
	return nil, nil
}
func Unmarshal(peersBin []byte) ([]Peer, error) {
	numsPeer := len(peersBin) / peerSize
	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("malformed compact peer data: length %d is not divisible by %d", len(peersBin), peerSize)
	}
	peers := make([]Peer, numsPeer)
	for i := range numsPeer {
		offset := i * peerSize
		peers[i].IP = net.IP(peersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peersBin[offset+4 : offset+6])
	}
	return peers, nil
}

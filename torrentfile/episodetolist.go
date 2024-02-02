package torrentfile

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
	"github.com/veggiedefender/torrent-client/p2p"
	"os"
)

type files = []file

type file struct {
	length int
	path   string
}

type mbencodeInfo struct {
	Files       files  `bencode:"files"`
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type mbencodeTorrent struct {
	Announce string       `bencode:"announce"`
	Info     mbencodeInfo `bencode:"info"`
}

func OpenMulti(path string, episode int) (TorrentFile, error, int, int) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err, 0, 0
	}
	defer file.Close()

	bto := mbencodeTorrent{}
	err = bencode.Unmarshal(file, &bto)
	if err != nil {
		return TorrentFile{}, err, 0, 0
	}
	return bto.toTorrentFile(episode)
}

func (i *mbencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *mbencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (i *mbencodeInfo) splitEpisodeHash(episodeNum int) ([][20]byte, int, int) {
	episodeLength := i.Files[episodeNum-1].length //episode num 2 length is 2 gb
	interval := episodeLength / i.PieceLength     //piece size is 1 mb, for 2 gb we must download 2000 pure pieces
	inx, beginOff, endOff := i.beginEndIndex(2)
	if endOff != 0 {
		interval += 1
	}
	fullhash, _ := i.splitPieceHashes()
	var partialHash = make([][20]byte, interval)
	copy(partialHash[:], fullhash[inx:inx+interval])
	return partialHash, beginOff, endOff
}

func (i *mbencodeInfo) beginEndIndex(episodeNum int) (int, int, int) {
	beginOffset := 0
	var endOffset int
	start := 0

	for e := 0; e < episodeNum-1; e++ {
		start += i.Files[e].length / i.PieceLength
		offset := (i.Files[e].length - beginOffset) % i.PieceLength //episode length is 2000,5 mb pieceSize is 1 mb, so we have 0,5 mb of the second file in buf
		beginOffset = (beginOffset * 0) + offset
		if offset != 0 {
			start += 1
		}

	}
	endOffset = (i.Files[episodeNum-1].length - beginOffset) % i.PieceLength
	return start, beginOffset, endOffset
}

func (bto *mbencodeTorrent) toTorrentFile(episode int) (TorrentFile, error, int, int) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err, 0, 0
	}
	pieceHashes, beginOff, endOff := bto.Info.splitEpisodeHash(2)

	t := TorrentFile{
		Announce:    bto.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Files[episode].length + endOff,
		Name:        bto.Info.Files[episode].path,
	}
	return t, nil, beginOff, endOff
}

func (t *TorrentFile) DownloadEpisodeToFile(path string, begin int, end int) error {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return err
	}

	peers, err := t.requestPeers(peerID, Port)
	if err != nil {
		return err
	}

	torrent := p2p.Torrent{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
	}
	buf, err := torrent.Download()
	if err != nil {
		return err
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.Write(buf[begin:end])
	if err != nil {
		return err
	}
	return nil
}

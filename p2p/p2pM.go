package p2p

import (
	"log"
	"runtime"
)

const (
	CONTINUE = 1
	STOP     = 0
	ERROR    = -1
)

func MultiDownload(signals []chan int, torrents []*Torrent) [][]byte {
	results := make([][]byte, len(signals))
	for i, torrent := range torrents {
		torrent.DownloadPiece(signals[i])
	}
	index := len(signals)
	for {
		for i, torrent := range torrents {
			buf, err := torrent.DownloadPiece(signals[i])
			if err == nil {
				results[i] = buf
				// TODO Drop Torrents
				// TODO add result to current Torrent
			}
		}
		index += 1
		index %= len(signals)
		signals[index] <- CONTINUE
	}
	return results
}

func (t *Torrent) DownloadPiece(torrent chan int) ([]byte, error) {
	log.Println("Starting downloading for", t.Name)

	// Init queues for workers to retrieve work and send results
	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)

	for index, hash := range t.PieceHashes {
		length := t.calculatePieceSize(index)
		workQueue <- &pieceWork{index, hash, length}
	}

	for _, peer := range t.Peers {
		go t.startDownloadWorker(peer, workQueue, results)
		switch <-torrent {
		case CONTINUE:
		case ERROR:
			return nil, nil
		case STOP:
			break
		}

	}
	// Collect results into a buffer until full
	buf := make([]byte, t.Length)
	donePieces := 0
	for donePieces < len(t.PieceHashes) {
		res := <-results

		begin, end := t.calculateBoundsForPiece(res.index)
		copy(buf[begin:end], res.buf)
		donePieces++

		percent := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		numWorkers := runtime.NumGoroutine() - 1 // subtract 1 for main thread
		log.Printf("(%0.2f%%) Downloaded piece #%d from %d peers\n", percent, res.index, numWorkers)
	}
	close(workQueue)

	return buf, nil
}

package threading

import (
	"github.com/veggiedefender/torrent-client/client"
)

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

type IO struct {
	input chan pieceWork
	out   chan pieceResult
}

func newIO(length int) IO {
	i := make(chan pieceWork, length) //unbuffered
	o := make(chan pieceResult)       //buffered

	return IO{i, o}
}

// new queque with episode pieces at first chan, and unpriorites
func redesignDownloadPriority(beginIndex int, endIndex int, wq chan *pieceWork) (chan *pieceWork, chan *pieceWork) {
	partial := make(chan *pieceWork, len(wq))
	newCh := make(chan *pieceWork, len(wq))

	for dt := range wq {
		if dt.index < beginIndex || endIndex < dt.index {
			newCh <- dt
		}
		partial <- dt
	}

	//wq.assert( partial)

	return wq, partial

}

// now queque is fullfield only with priority connections
func redesignClientsPriority(beginInx int, endInx int, cl []*client.Client) []*client.Client {
	var priorityCl []*client.Client
	for i, v := range cl {
		if v.Bitfield.HasPiece(i+beginInx) && (i+beginInx) < endInx {
			priorityCl = append(priorityCl, v)
		}
	}
	return priorityCl
}

func episodeDownload(obs *observer) { //pause for client
	obs.ctrl <- 2
	//attemptDownload(newobs){<-newobs...onRabbitMqBlock {newobs<-1} onRabbitMqUnblock{<-newobs}newobs<-1
	<-obs.ctrl
}

//go func acceleratedDownload(bufobs *bufobserver, begin,end int) {
//	wq := make(chan *pieceWork, end-begin)
//	prCh, lastCh := redesignDownloadPriority(begin, end, wq)   //exp 3 threads always async download
//
//	//attemptdownload()
//	//attemptdownload(obs){<-bufobs pq<-wq ...result<-&workPiece bufobs<-1  } //1 result recieved => 2 result rec=> 3 res rec=>1
//}

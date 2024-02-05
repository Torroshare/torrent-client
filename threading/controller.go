package threading

import (
	"fmt"
	"github.com/veggiedefender/torrent-client/p2p"
	"github.com/veggiedefender/torrent-client/rabbitmq"
)

const (
	stpS = "1"
	stpM = "2"
	stpA = "3"

	pauseS = "4"
	pauseM = "5"
	pauseA = "6"
)

type observer struct {
	ctrl chan int
}

func newObserver() *observer {
	obs := make(chan int, 0)
	ins := observer{obs}
	return &ins
}

// redesign for select later
func (o *observer) rbmqObserver(uuid string) {
	rbmq := rabbitmq.NewRabbitMQ(uuid)
	rabbitmq.FromQueque(rbmq)

	msg := string(<-rbmq.Messages)
	switch msg {

	case stpS + uuid:
		o.ctrl <- 1
		break

	case "pause" + uuid:
		o.ctrl <- 2
		break

	case "cls" + uuid:
		o.ctrl <- 3
		break

	default:
		fmt.Println("command not found")
	}
}

func singlefileDownload(t *p2p.Torrent,memory []byte, sep chan []byte) error {
	obs := newObserver()
	obs.rbmqObserver(t.Name)

	res, err := t.Download() //sep

	if err != nil {
		return  err
	}
	copy(memory[:],res)
	return nil
}

func multifileDownload(t *p2p.Torrent,memory []byte, sep chan []byte) error{
	obs := newObserver()
	obs.rbmqObserver(t.Name)

	// t.downloadMulti(res)
	return nil
}

func mixedDownload(t *p2p.Torrent,memory []byte, sep chan []byte) ([]byte,error){
	obs := newObserver()
	obs.rbmqObserver(t.Name)
	res, err := t.Download() //sep

	if err != nil {
		return  nil,err
	}

	return res,nil
}

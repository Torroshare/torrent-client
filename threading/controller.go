package threading

type observer struct {
	ctrl chan byte
}

func newObserver() *observer {
	obs := make(chan byte, 1)
	ins := observer{obs}
	return &ins
}

func (o *observer) StopQueque() { //we stop buffered thread until the episode won't be downloaded
	o.ctrl <- 1
}

package barutil

type task struct {
	name     string
	worker   func() error
	done     chan struct{}
	finished bool
	err      error
}

func (t *task) start() {
	defer close(t.done)
	if t.worker == nil {
		return
	}
	t.err = t.worker()
}

func (t *task) wait() {
	<-t.done
	t.finished = true
}

package barutil

type task struct {
	name     string
	worker   func(string) error
	done     chan struct{}
	finished bool
	err      error
}

func (t *task) start() {
	defer close(t.done)
	if t.worker == nil {
		return
	}
	t.err = t.worker(t.name)
}

func (t *task) wait() {
	<-t.done
	t.finished = true
}

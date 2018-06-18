package inmem

import "context"

type incrRequest struct {
	id   string
	by   int64
	done chan struct{}
}

func newIncrRequest(id string) incrRequest {
	return incrRequest{
		id:   id,
		by:   1,
		done: make(chan struct{}),
	}
}

type getRequest struct {
	id string
	rv chan int64
}

func newGetRequest(id string) getRequest {
	return getRequest{
		id: id,
		rv: make(chan int64),
	}
}

type Counter struct {
	c map[string]int64

	incrChan chan incrRequest
	getChan  chan getRequest
	stopChan chan chan struct{}
}

// New creates a in-memory counter instance.
func New() *Counter {
	return &Counter{
		c:        map[string]int64{},
		incrChan: make(chan incrRequest),
		getChan:  make(chan getRequest),
		stopChan: make(chan chan struct{}),
	}
}

// Start starts counter execution (block).
func (c *Counter) Start() error {
	for {
		select {
		case s := <-c.stopChan:
			close(s)
			return nil
		case incr := <-c.incrChan:
			c.incrWith(incr.id, incr.by)
			close(incr.done)
		case get := <-c.getChan:
			get.rv <- c.get(get.id)
		}
	}
}

// Stop terminates executeion.
func (c Counter) Stop() {
	stop := make(chan struct{})
	c.stopChan <- stop
	<-stop
}

func (c *Counter) incrWith(id string, by int64) {
	counter, exists := c.c[id]
	if !exists {
		counter = 0
	}
	c.c[id] = counter + 1
}

func (c *Counter) get(id string) int64 {
	if counter, exists := c.c[id]; exists {
		return counter
	}
	return 0
}

func (c Counter) Incr(ctx context.Context, id string) error {
	req := newIncrRequest(id)

	c.incrChan <- req

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-req.done:
		return nil
	}
}

func (c Counter) Get(ctx context.Context, id string) (int64, error) {
	req := newGetRequest(id)

	c.getChan <- req

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case rv := <-req.rv:
		return rv, nil
	}
}

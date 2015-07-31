package worker

import (
	"log"
	"sync"
	"time"
)

type Pool struct {
	Buffer    int
	queue     chan bool
	wg        sync.WaitGroup
	start     time.Time
	completed int
	total     int
}

func NewPool(buffer int) *Pool {
	return &Pool{
		Buffer: buffer,
		queue:  make(chan bool, buffer),
		start:  time.Now(),
	}
}

func (p *Pool) Task(f func()) {
	p.wg.Add(1)
	p.queue <- true
	go func() {
		f()
		<-p.queue
		p.wg.Done()
		p.completed++
	}()
}

func (p *Pool) Time() {
	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			p.total += p.completed
			log.Println(p.total, "[", p.completed, "/s ]")
			p.completed = 0
		}
	}()
}

func (p *Pool) Wait() int {
	p.wg.Wait()
	log.Println(time.Now().Sub(p.start))
	return p.total
}

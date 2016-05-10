package worker

import (
	"log"
	"sync"
	"time"
)

type Pool struct {
	Buffer    int
	Name      string
	queue     chan bool
	result    chan interface{}
	wg        sync.WaitGroup
	start     time.Time
	completed int
	total     int
}

func NewPool(buffer int, name string) *Pool {
	return &Pool{
		Name:   name,
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
	log.Println(p.Name, time.Now().Sub(p.start))
	return p.total
}

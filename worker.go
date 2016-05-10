package worker

import (
	"log"
	"sync"
	"time"
)

type Worker struct {
	Produce   func(chan []interface{})
	Consume   func([]interface{}) bool
	Silent    bool
	Count     int
	Retry     int
	completed int
	total     int
}

func (worker *Worker) Run() {
	queue := make(chan []interface{}, 1000)
	worker.completed = 0
	worker.total = 0
	ticker := time.NewTicker(time.Second)
	go func() {
		for _ = range ticker.C {
			worker.total += worker.completed
			if !worker.Silent {
				log.Println(worker.total, "[", worker.completed, "/s ]")
			}
			worker.completed = 0
		}
	}()
	var wg sync.WaitGroup
	if worker.Count == 0 {
		worker.Count = 500
	}
	if !worker.Silent {
		log.Println("Spawning", worker.Count, "workers")
	}
	for i := 0; i < worker.Count; i++ {
		wg.Add(1)
		go worker.spin(i, queue, &wg)
	}
	if !worker.Silent {
		log.Println("Working...")
	}
	worker.Produce(queue)
	close(queue)
	wg.Wait()
	ticker.Stop()
	if !worker.Silent {
		log.Println("Completed", worker.total+worker.completed, "tasks")
	}
}

func (worker *Worker) spin(id int, queue chan []interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range queue {
		count := 0
		for {
			success := worker.Consume(line)
			count++
			if success {
				break
			}
			if count > worker.Retry {
				worker.completed--
				break
			}
		}
		worker.completed++
	}

}

package workerpool

import (
	"errors"
	"log"
	"sync"
)

type Worker struct {
	ID     int
	stop   chan bool
	jobs   chan string
	result chan struct{}
}

func newWorker(ID int, stop chan bool, jobs chan string, result chan struct{}) *Worker {
	return &Worker{
		ID:     ID,
		stop:   stop,
		jobs:   jobs,
		result: result,
	}
}

func (worker *Worker) Process() {
	for {
		select {
		case <-worker.stop:
			log.Printf("worker %d done", worker.ID)
			return
		case j, ok := <-worker.jobs:
			if !ok {
				return
			}
			log.Printf("worker %d processed: %s", worker.ID, j)
			worker.result <- struct{}{}
		}
	}
}

type WorkerPool struct {
	workers []*Worker
	wg      *sync.WaitGroup
	jobs    chan string
	result  chan struct{}
}

func NewWorkerPool(workersAmount int, jobs chan string, result chan struct{}) (*WorkerPool, error) {
	if workersAmount < 0 {
		return nil, errors.New("negative worker-pool size")
	}
	workers := make([]*Worker, 0, workersAmount)
	for i := 0; i < workersAmount; i++ {
		stop := make(chan bool)
		workers = append(workers, newWorker(i+1, stop, jobs, result))
	}
	wg := sync.WaitGroup{}
	return &WorkerPool{
		workers: workers,
		wg:      &wg,
		jobs:    jobs,
		result:  result,
	}, nil
}

func (wp *WorkerPool) Add() {
	ID := len(wp.workers) + 1
	stop := make(chan bool)
	worker := newWorker(ID, stop, wp.jobs, wp.result)
	wp.workers = append(wp.workers, worker)
	go worker.Process()
}

func (wp *WorkerPool) Shrink() error {
	if len(wp.workers) == 0 {
		return errors.New("worker-pool is empty")
	}
	worker := wp.workers[0]
	worker.stop <- true
	wp.workers = wp.workers[1:]
	return nil
}

func (wp *WorkerPool) Run() {

	for _, worker := range wp.workers {
		go worker.Process()
	}
}

func (wp *WorkerPool) Results() <-chan struct{} {
	return wp.result
}

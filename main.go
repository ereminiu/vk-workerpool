package main

import (
	"fmt"
	"log"
	"time"

	workerpool "github.com/ereminiu/vk-workerpool/worker-pool"
)

const N = 10

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// rand.Seed(time.Now().Unix())

	ch := make(chan string, N)
	result := make(chan struct{}, N)

	M := 3
	wp, err := workerpool.NewWorkerPool(M, ch, result)
	if err != nil {
		log.Fatalln(err)
	}

	go wp.Run()

	go wp.Add()

	go func() {
		time.Sleep(1 * time.Second)
		log.Printf("second is left")
		wp.Shrink()
	}()

	for j := 0; j < N; j++ {
		log.Printf("add job %d", j)
		job := fmt.Sprintf("job #%d", j)
		ch <- job
		time.Sleep(250 * time.Millisecond)
	}
	close(ch)

	for i := 0; i < N; i++ {
		<-result
	}
	close(result)
}

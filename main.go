package main

import (
	"fmt"
	"log"

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

	for j := 0; j < N; j++ {
		ch <- fmt.Sprintf("job #%d", j)
	}
	close(ch)

	for i := 0; i < N; i++ {
		<-result
	}
}

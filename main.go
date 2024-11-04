package main

import (
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	// rand.Seed(time.Now().Unix())

	const N = 10

	stop := make(chan int)
	ch := make(chan int, N)
	res := make(chan int, N)

	M := 1
	for w := 0; w < M; w++ {
		go process(w, stop, ch, res)
	}
	go process(288, stop, ch, res)

	go func() {
		time.Sleep(1 * time.Second)
		log.Printf("1 sec left")
		stop <- 1
	}()

	for j := 0; j < 10; j++ {
		log.Printf("add %d", j)
		ch <- j
		time.Sleep(250 * time.Millisecond)
	}
	close(ch)

	for j := 0; j < N; j++ {
		<-res
	}
	close(res)
}

func process(w int, stop <-chan int, ch <-chan int, res chan<- int) {
	for {
		select {
		case <-stop:
			log.Printf("done")
			return
		case j, ok := <-ch:
			if !ok {
				break
			}
			log.Printf("worker: %d, processed: %d\n", w, j)
			res <- j
		}
	}
}

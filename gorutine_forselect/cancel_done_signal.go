package main

import "fmt"
import "time"

func do_stuff() int {
	return 1
}

func main() {

	ch := make(chan int, 100)
	//canal para hacer signaling de cuando terminar
	done := make(chan struct{})
	go func() {
		for {
			select {
			case ch <- do_stuff():
			case <-done:
				close(ch)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	//se dispara una go rutina que tres segundos despues acaba con el done
	go func() {
		time.Sleep(3 * time.Second)
		done <- struct{}{}
	}()

	for i := range ch {
		fmt.Println("receive value: ", i)
	}

	fmt.Println("finish")
}

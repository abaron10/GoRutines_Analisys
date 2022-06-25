package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

func pings(urls string) (int, error) {
	resp, err := http.Get(urls)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

func main() {

	urls := []string{
		"https://www.google.com/?hl=es",
		"https://www.facebook.com",
		"https://www.tesla.com",
		"https://www.google.com/?hl=es",
		"https://www.facebook.com",
		"https://www.tesla.com",
		"https://www.teslaNOEXISTO.com",
	}
	var wg sync.WaitGroup
	ch := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, _ := pings(u)
			select {
			case <-ctx.Done():
				return
			case ch <- resp:

			default:

			}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		resp := <-ch
		if resp == 0 {
			fmt.Println("error")
			//cancel()
			continue
		}
		fmt.Println(resp)
	}
	wg.Wait()
	cancel()
	fmt.Println("end")
}

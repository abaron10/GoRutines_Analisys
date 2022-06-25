package main

import (
	"bufio"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// func transform(word *[]string){
// 	(*word)[0] = "ämigo"
// 	(*word) = append((*word), "ssss")
// }

// func main(){
// 	word := []string{"hola"}
// 	transform(&word)
// 	fmt.Println(word)

// }

func get(num int, wg *sync.WaitGroup) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/" + strconv.Itoa(num))
	if err != nil {
		panic(err)
	}
	defer func() {
		resp.Body.Close()
		wg.Done()
	}()
	fmt.Println("Status: ", resp.Status)

	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	//Para poder esperar y no bloquear el hilo principal, de no ser necesario mandar o usar un canal, puedo usar un wait group
	//Con esto adiciono un marcador a cada routina y con el wg.wait() espero a que todas se completen.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go get(i, &wg)
	}
	wg.Wait()

}

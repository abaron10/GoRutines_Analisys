package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

//funcion random donde ejecuto logica
func ping(urls string) *Response {
	resp, err := http.Get(urls)
	if err != nil {
		return &Response{nil, err}

	}
	return &Response{resp, nil}
}

type Response struct {
	resp *http.Response
	err  error
}

func main() {
	//urls con intención de que una no exista y falle con error
	urls := []string{
		"https://www.google.com/?hl=es",
		"https://www.facebook.com",
		"https://www.tesla.com",
		"https://www.google.com/?hl=es",
		"https://www.facebook.com",
		"https://www.tesla.com",
		//"https://www.teslaNOEXISTO.com",
	}
	//canal de transmisión de respuesta
	chPG := make(chan *Response, 2)
	var wg sync.WaitGroup
	//la forma correcta de cerrar el canal es con el contexto
	ctx, cancelWorkers := context.WithCancel(context.Background())
	for _, l := range urls {
		//adición de un waitgroup para que así se cierre el canal, espere a las otras que falten si ya empezarón
		wg.Add(1)
		go func(l string) {
			wg.Done()
			for {
				select {
				//si se llega a terminar el contexto
				case <-ctx.Done():
					return
				case chPG <- ping(l):
				}
			}

		}(l)
	}
	for i := 0; i < len(urls); i++ {
		e := <-chPG
		//si llega a  existir un error cancelo el contexto
		if e.err != nil {
			fmt.Println(e.err.Error())
			cancelWorkers()
			break
		}
		//logica con ese mensaje del canal
		fmt.Println(e.resp.StatusCode)
	}
	//espero a que las go rutinas terminen
	wg.Wait()
	//cierro el canal
	close(chPG)
	fmt.Println("end")
}

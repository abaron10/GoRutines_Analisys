package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	UserId    int32  `json:"userId"`
	Id        int32  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type UserDetail struct {
	Email    string  `json:"email"`
	Username string  `json:"name"`
	Address  address `json:"address"`
}

type address struct {
	Street string `json:"street"`
	Suite  string `json:"suite"`
}

type UserResponse struct {
	Id      int32   `json:"id"`
	Name    string  `json:"name"`
	Title   string  `json:"title"`
	Address address `json:"address"`
}

type processor struct {
	userchan  chan User
	todoschan chan UserDetail
	errs      chan error
}

func getData(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *processor) process(ctx context.Context, id int) {
	go func() {
		var user User
		response, err := p.gather(fmt.Sprintf("https://jsonplaceholder.typicode.com/todos/%d", 1))
		if err != nil {
			p.errs <- err
			return
		}
		json.Unmarshal(response, &user)
		p.userchan <- user
	}()
	go func() {
		var userDetail UserDetail
		response, err := p.gather(fmt.Sprintf("https://jsonplaceholder.typicode.com/users/%d", 1))
		if err != nil {
			p.errs <- err
			return
		}
		json.Unmarshal(response, &userDetail)
		p.todoschan <- userDetail
	}()
}

func (p *processor) waitforTodos(ctx context.Context) (UserResponse, error) {
	//Esta función es necesaria junto con el contador, ya que sin este, el for select al
	//encontrarse con un caso positivo, sale y no evalua el otro, con este contador forzamos
	//al sistema para que espere a 2 mensajes.
	var userResponse UserResponse
	count := 0
	for count < 2 {
		select {
		case <-ctx.Done():
			return UserResponse{}, ctx.Err()
		case userChannel := <-p.todoschan:
			userResponse.Address = userChannel.Address
			userResponse.Name = userChannel.Username
			count++
		case userChannel := <-p.userchan:
			userResponse.Id = userChannel.UserId
			userResponse.Title = userChannel.Title
			count++
		}
	}
	return userResponse, nil
}

func (p *processor) gather(url string) ([]byte, error) {
	response, err := getData(url)
	defer func() {
		response.Body.Close()
	}()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println(err2.Error())
		return nil, err2
	}
	return body, nil
}

func main() {
	//cada canal esta buffereado para que las go rutinas que escriban puedan salir despues de escribir y no
	//esperar a que un read suceda, si no estuviera buffereado, la go rutina se bloquea hasta que se tenga una
	//lectura. El canal de errores puede tener hasta dos errores de escritura en el buffer. Si no estuviera buffereado
	//necesitaria cerrar la gorutina cerrando el canal
	p := &processor{make(chan User, 1), make(chan UserDetail, 1), make(chan error, 2)}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Millisecond)
	defer cancel()
	p.process(ctx, 1)

	res, err := p.waitforTodos(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("El usuario %s escribio %s en la calle %s", res.Name, res.Title, res.Address.Suite)
	fmt.Println("end")
}

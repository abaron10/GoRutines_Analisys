package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type User struct{
	UserId int32 `json:"userId"`
	Id int32 `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}
func get(num int, c chan User){
	var person User

	resp,err := http.Get("https://jsonplaceholder.typicode.com/todos/" + strconv.Itoa(num))
	if err != nil {
		panic(err)
	}
	//No importa la posición del defer, siempre despues de ejecutar este método se va a ejecturar( sirve para cerrar sesiones , etc)
	defer func(){
		resp.Body.Close()
	}()
	//fmt.Println("Status: ",resp.Status)

	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil{
		fmt.Println("mistakes")
	}
	json.Unmarshal(body, &person)
	c<- person
}

func main(){
	response := []User{}
	c:= make(chan User)
	//Manda las 1000 go routines
	for i:=0;i < 1000;i++{
		go get(i,c)
	}
	//Espera a los 1000 mensajes 
	for i:=0; i < 1000;i++{
		response = append(response,<-c)
	}
	fmt.Println(response)
}
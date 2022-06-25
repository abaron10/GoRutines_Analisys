package main

import "fmt"

func main() {

	test := "cambio de valor"
	cambio(&test)
	fmt.Println(test)

}

func cambio(s *string) {
	(*s) = "camb"
}

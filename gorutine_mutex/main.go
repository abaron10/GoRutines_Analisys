package main

import (
	"fmt"
	"sync"
)

var (
	mutex   sync.Mutex
	balance int
)

//Lock(): only one go routine read/write at a time by acquiring the lock.

//RLock(): multiple go routine can read(not write) at a time by acquiring the lock.

func deposit(value int, wg *sync.WaitGroup) {
	mutex.Lock()
	fmt.Printf("Depositing %d to account with balance %d\n", value, balance)
	balance += value
	mutex.Unlock()
	wg.Done()
}
func withdraw(value int, wg *sync.WaitGroup) {
	//mutex ayuda a bloquear la operación es decir las otras go rutinas no pueden acceder y quedan en espera.
	mutex.Lock()
	fmt.Printf("withdrawing %d from account with balance %d\n", value, balance)
	balance -= value
	//se desbloquea la operación
	mutex.Unlock()
	wg.Done()

}
func main() {
	balance = 1000
	//waitgroup para esperar a dos go rutinas
	var wg sync.WaitGroup
	wg.Add(2)
	//independiente de cual acabe primero la operación esta bloqueada
	go withdraw(700, &wg)
	go deposit(500, &wg)
	//espera a terminar la operación
	wg.Wait()
	fmt.Printf("New balance %d\n", balance)
}

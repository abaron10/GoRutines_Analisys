package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
)

func Fibonacci(n int) int {
	if n <= 1 {

		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

type Memory struct {
	f Function
	//aca se cachea e.i {"5":..}
	cache map[int]FunctionResult
}

type Function func(key int) (interface{}, error)

type FunctionResult struct {
	value interface{}
	err   error
}

func NewCache(f Function) *Memory {
	return &Memory{
		f:     f,
		cache: make(map[int]FunctionResult),
	}
}

func (m *Memory) Get(key int) (interface{}, error) {
	result, exists := m.cache[key]

	if !exists {
		mutex.Lock()
		result.value, result.err = m.f(key)
		m.cache[key] = result
		mutex.Unlock()
	}
	return result.value, result.err
}

func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

func main() {
	cache := NewCache(GetFibonacci)
	fibo := []int{42, 40, 41, 42, 32}
	var wg sync.WaitGroup
	for _, n := range fibo {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			start := time.Now()
			value, _ := cache.Get(index)
			fmt.Printf("%d, %s, %d\n", index, time.Since(start), value)
		}(n)

	}
	wg.Wait()
}

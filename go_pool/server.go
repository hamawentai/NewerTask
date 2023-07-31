package main

import (
	"fmt"
	pool "gopool/mygopool"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	p, _ := pool.NewPool(3)
	var wg sync.WaitGroup
	dec := func(f func(http.ResponseWriter, *http.Request, *pool.Pool), p *pool.Pool) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			wg.Add(1)
			f(w, r, p)
			wg.Done()
		}
	}
	wg.Wait()
	http.HandleFunc("/ping", dec(myfunc, p))
	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func myfunc(w http.ResponseWriter, r *http.Request, p *pool.Pool) {
	ch := make(chan string)
	p.Submit(func() {
		fmt.Println("come")
		time.Sleep(2 * time.Second)
		ch <- "ping"
	})
	s := <- ch
	fmt.Fprintf(w, s)
}


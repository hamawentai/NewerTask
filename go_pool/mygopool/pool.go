package mygopool

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type Worker struct {
	pool *Pool
	task chan func()
	name string
}

func (w *Worker) run() {

	go func() {
		for task := range w.task {
			if task != nil {
				fmt.Println(w.name)
				task()
				w.pool.recoveryWorker(w)
			} else {
				w.pool.changeRunWorkersNum(-1)
				return
			}
		}
	}()
}

type Pool struct {
	cap         int32
	num_working int32
	workers     []*Worker
	lock        sync.Mutex
	cnt         int32
}

var (
	ErrInputArg = errors.New("Input parameter error")
)

type Task func(int)

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, ErrInputArg
	}
	p := &Pool{
		cap:         int32(size),
		num_working: 0,
		cnt:         0,
	}
	return p, nil
}

func (p *Pool) FindWorker() *Worker {
	var w *Worker
	waiting := false
	p.lock.Lock()
	idles := p.workers
	n := len(idles)
	if n <= 0 {
		waiting = p.Runs() >= p.Cap()
	} else {
		w = idles[n-1]
		idles[n-1] = nil
		p.workers = idles[:n-1]
	}
	p.lock.Unlock()
	if waiting {
		for {
			p.lock.Lock()
			idles = p.workers
			l := len(idles) - 1
			if l < 0 {
				p.lock.Unlock()
				continue
			}
			w = idles[l]
			idles[l] = nil
			p.workers = idles[:l]
			p.lock.Unlock()
			break
		}
	} else if w == nil {
		w = &Worker{
			pool: p,
			task: make(chan func(), 1),
			name: fmt.Sprintf("worker %d", p.getCnt()),
		}
		p.addcnt()
		w.run()
		p.changeRunWorkersNum(1)
	}
	return w
}

func (p *Pool) Submit(task func()) error {
	w := p.FindWorker()
	if w == nil {
		fmt.Println("w is nil")
		return ErrInputArg
	}
	w.task <- task
	return nil
}

func (p *Pool) recoveryWorker(w *Worker) {
	p.lock.Lock()
	p.workers = append(p.workers, w)
	p.lock.Unlock()
}

func (p *Pool) Cap() int {
	return int(atomic.LoadInt32(&p.cap))
}

func (p *Pool) Runs() int {
	return int(atomic.LoadInt32(&p.num_working))
}

func (p *Pool) getCnt() int {
	return int(atomic.LoadInt32(&p.cnt))
}
func (p *Pool) addcnt() {
	atomic.AddInt32(&p.cnt, int32(1))
}

func (p *Pool) changeRunWorkersNum(n int) {
	atomic.AddInt32(&p.num_working, int32(n))
}

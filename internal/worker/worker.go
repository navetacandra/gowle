package worker

import (
	"sync"
	"time"
)

type Worker struct {
	mu        sync.Mutex
	wg        sync.WaitGroup
	isRunning bool
	onStart   func()
	onStop    func()
	onLoop    func()
}

func NewWorker(onStart func(), onLoop func(), onStop func()) *Worker {
	return &Worker{
		onStart: onStart,
		onStop:  onStop,
		onLoop:  onLoop,
	}
}

func (w *Worker) run() {
	defer w.wg.Done()
	w.onStart()
	for {
		w.mu.Lock()

		if !w.isRunning { // clean-up process
			w.onStop()
			w.mu.Unlock()
			return
		}

		w.onLoop()
		w.mu.Unlock()
		time.Sleep(800 * time.Millisecond)
	}
}

func (w *Worker) Start() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isRunning {
		return
	}

	w.isRunning = true
	w.wg.Add(1)
	go w.run()
}

func (w *Worker) Stop() {
	w.mu.Lock()

	if !w.isRunning {
		w.mu.Unlock()
		return
	}

	w.isRunning = false
	w.mu.Unlock()
	w.wg.Wait()
}

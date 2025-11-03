package worker

import (
	"sync"
	"time"
)

type Worker struct {
	mu        sync.Mutex
	wg        sync.WaitGroup
	isRunning bool
	onStart   func(*Worker)
	onStop    func(*Worker)
	onLoop    func(*Worker)
	debounce  time.Duration
}

func NewWorker(onStart func(*Worker), onLoop func(*Worker), onStop func(*Worker), debounce time.Duration) *Worker {
	return &Worker{
		onStart:  onStart,
		onStop:   onStop,
		onLoop:   onLoop,
		debounce: debounce,
	}
}

func (w *Worker) run() {
	defer w.wg.Done()
	w.onStart(w)
	for {
		time.Sleep(w.debounce)
		w.mu.Lock()

		if !w.isRunning { // clean-up process
			w.onStop(w)
			w.mu.Unlock()
			return
		}

		w.onLoop(w)
		w.mu.Unlock()
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

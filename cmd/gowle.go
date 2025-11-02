package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/navetacandra/gowle/internal/config"
	"github.com/navetacandra/gowle/internal/worker"
)

func main() {
	fmt.Println("[GOWLE] Start.")

	sigs := make(chan os.Signal, 1)
	scanner := bufio.NewScanner(os.Stdin)
	appConfig := config.GowleConfig{}

	appWorker := worker.NewWorker(func() { // on start
		appConfig.Load()
	}, func() { // on loop
	}, func() { // on stop
	})

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // handle external termination
	go func() {
		<-sigs
		handleTermination(appWorker)
	}()

	appWorker.Start()
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "rs":
			fmt.Println("[GOWLE] Reloading.")
			appWorker.Stop()
			appWorker.Start()
		case ".exit":
			handleTermination(appWorker)
		}
	}
}

func handleTermination(w *worker.Worker) {
	w.Stop()
	fmt.Println("[GOWLE] Stoped.")
	os.Exit(0)
}

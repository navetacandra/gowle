package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/navetacandra/gowle/internal/config"
	"github.com/navetacandra/gowle/internal/fsscan"
	"github.com/navetacandra/gowle/internal/spawn"
	"github.com/navetacandra/gowle/internal/worker"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose mode. Print every change")
	flag.Parse()
	fmt.Println("[GOWLE] Start.")

	sigs := make(chan os.Signal, 1)
	scanner := bufio.NewScanner(os.Stdin)

	appConfig := config.GowleConfig{}
	appProcess := spawn.ChildProcess{}
	snapshot := []fsscan.Info{}
	snapshotDiff := []fsscan.DiffInfo{}

	appWorker := worker.NewWorker(func(w *worker.Worker) { // on start
		appConfig.Load()
		fsscan.Scan(&snapshot, &appConfig)

		if err := appProcess.Start(&appConfig); err != nil {
			fmt.Printf("[GOWLE] Failed to start child process: %v\n", err)
			handleTermination(w, &appProcess)
		}
	}, func(w *worker.Worker) { // on loop
		tmpSnapshot := []fsscan.Info{}
		fsscan.Scan(&tmpSnapshot, &appConfig)
		fsscan.DiffSnapshot(&tmpSnapshot, &snapshot, &snapshotDiff)

		if len(snapshotDiff) > 0 {
			snapshot = tmpSnapshot
			if *verbose {
				for _, d := range snapshotDiff {
					mode := "modified"
					switch d.Diff {
					case -1:
						mode = "deleted"
					case 1:
						mode = "created"
					}
					fmt.Printf("[GOWLE] %s was %s\n", d.Path, mode)
				}
			}

			fmt.Println("[GOWLE] Reloading.")
			if err := appProcess.Stop(); err != nil {
				fmt.Printf("[GOWLE] Failed to stop child process: %v\n Try to restart with command 'rs'\n", err)
			}
			if err := appProcess.Start(&appConfig); err != nil {
				fmt.Printf("[GOWLE] Failed to start child process: %v\n Try to restart with command 'rs'\n", err)
			}
		}
	}, func(w *worker.Worker) { // on stop
		if err := appProcess.Stop(); err != nil {
			fmt.Printf("[GOWLE] Failed to stop child process: %v\n Try to restart with command 'rs'\n", err)
		}
	}, 1000*time.Millisecond)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // handle external termination
	go func() {
		<-sigs
		handleTermination(appWorker, &appProcess)
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
			handleTermination(appWorker, &appProcess)
		}
	}

	// clean-up
	handleTermination(appWorker, &appProcess)
}

func handleTermination(w *worker.Worker, cp *spawn.ChildProcess) {
	cp.Stop()
	w.Stop()
	fmt.Println("[GOWLE] Stoped.")
	os.Exit(0)
}

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
	"github.com/navetacandra/gowle/internal/worker"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose mode. Print every change")
	flag.Parse()
	fmt.Println("[GOWLE] Start.")

	sigs := make(chan os.Signal, 1)
	scanner := bufio.NewScanner(os.Stdin)
	appConfig := config.GowleConfig{}
	snapshot := []fsscan.Info{}
	snapshotDiff := []fsscan.DiffInfo{}

	appWorker := worker.NewWorker(func() { // on start
		appConfig.Load()
		fsscan.Scan(&snapshot, &appConfig)
	}, func() { // on loop
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
					fmt.Printf("%s was %s\n", d.Path, mode)
				}
			}

			// TODO: Implement action on change
		}
	}, func() { // on stop
	}, 1000*time.Millisecond)

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

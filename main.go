package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	workDuration  = 2 * time.Minute
	breakDuration = 5 * time.Minute
)

func runSession(duration time.Duration, label string) {
	fmt.Printf("\n--- %s started (%v) ---\n", label, duration)

	bar := progressbar.NewOptions64(
		duration.Milliseconds()/1000, // total seconds
		progressbar.OptionSetDescription(label),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetItsString("sec"),
		progressbar.OptionClearOnFinish(),
	)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	end := time.Now().Add(duration)
	for range ticker.C {
		remaining := time.Until(end)
		if remaining <= 0 {
			bar.Finish()
			fmt.Printf("%s done!\n", label)
			return
		}
		bar.Add(1)
	}
}

func main() {
	// Handle Ctrl+C clean exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nExiting Pomodoro. Stay productive!")
		os.Exit(0)
	}()

	for {
		runSession(workDuration, "Work")
		runSession(breakDuration, "Break")
	}
}

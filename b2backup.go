package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var pathListFile string
var ctx = context.Background()
var run = true

func main() {
	singleRun()
	processArgs()

	if err := loadPaths(); err != nil {
		log.Fatal(err)
	}

	loadSettings()

	// setup signal catching
	sigs := make(chan os.Signal, 1)

	// catch all signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// method invoked upon seeing signal
	go func() {
		s := <-sigs
		fmt.Printf(" RECEIVED SIGNAL: %s\n", s)
		run = false
	}()

	// Startup the watchers
	for i := range paths {
		fmt.Println("Connecting to bucket '" + paths[i].Bucket + "'...")
		paths[i].connect()
		fmt.Println("Starting watcher for '" + paths[i].Source + "'")
		go watchDirectory(&paths[i])
	}

	// TODO: Run signature checks accross current files

	// Stay alive
	for run {
		time.Sleep(5 * time.Second)
	}

	// Wait for all writes to go through
	time.Sleep(5 * time.Second)

	// Save out files
	for i := range paths {
		fmt.Println("Stopping watcher for '" + paths[i].Source + "'...")
		saveSettings(paths[i])
	}
}

func processArgs() {
	init := flag.String("init", "", "path to create a new backup")
	pathList := flag.String("paths", "paths.list", "location of global path keeper")
	flag.Parse()

	pathListFile = *pathList

	//Send off to init promt
	if *init != "" {
		newPath(*init)
	}

}

func singleRun() {
	if _, err := net.Listen("tcp", ":56210"); err != nil {
		log.Fatal("an instance was already running, please stop before starting")
	}
}

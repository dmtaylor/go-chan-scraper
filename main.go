package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dmtaylor/go-chan-scraper/util"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	MaxThreads uint   `short:"j" long:"threads" description:"Max number of downloader threads" default:"10"`
	Directory  string `short:"d" long:"directory" description:"download directory" default:"."`
}

func imgWorker(images, errors chan string, dir string) {
	// TODO
}

func processThread(url string, maxThreads uint, dir string) {
	fmt.Printf("Downloading thread %s\n", url)

	// TODO

}

func main() {
	var opts Options

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	isDir, err := util.DirExists(opts.Directory)
	if err != nil {
		log.Fatalf("Failed to get valid directory: %v\n", err)
	}
	if isDir == false {
		log.Fatalf("Error: %s is not a valid directory", opts.Directory)
	}

	threads := args[1:]

	for _, item := range threads {
		processThread(item, opts.MaxThreads, opts.Directory)
	}

	fmt.Printf("Done\n")

}

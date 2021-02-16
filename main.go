package main

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	MaxThreads uint   `short:"j" long:"threads" description:"Max number of downloader threads" default:"10"`
	Directory  string `short:"d" long:"directory" description:"download directory" default:"."`
}

func main() {
	var opts Options

	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	_ = args // here to silence unused var err

	// TODO scraping here
}

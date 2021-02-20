package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/dmtaylor/go-chan-scraper/extractors"
	"github.com/dmtaylor/go-chan-scraper/util"
	"github.com/jessevdk/go-flags"
)

type options struct {
	MaxThreads uint   `short:"j" long:"threads" description:"Max number of downloader threads" default:"10"`
	Directory  string `short:"d" long:"directory" description:"download directory" default:"."`
	Extractor  string `short:"e" long:"engine" description:"Site engine to use" choice:"4chan" choice:"8kun" default:"4chan"`
}

type extractor func(io.ReadCloser, chan extractors.ImageFile) error

func imgWorker(
	images chan extractors.ImageFile,
	errors chan extractors.ImageError,
	dir string,
	wg *sync.WaitGroup) {

	for img := range images {
		resp, err := http.Get(img.FileURL.String())
		if err != nil {
			errors <- extractors.ImageError{
				Err:     err,
				FileURL: img.FileURL,
			}
			continue
		}
		defer resp.Body.Close()
		err = util.ValidateHttpResponse(resp, nil)
		if err != nil {
			errors <- extractors.ImageError{
				Err:     err,
				FileURL: img.FileURL,
			}
			continue
		}

		filename := filepath.Join(dir, img.FileName)
		err = util.SaveFile(filename, resp.Body)
		if err != nil {
			errors <- extractors.ImageError{
				Err:     err,
				FileURL: img.FileURL,
			}
			continue
		}
		fmt.Printf("Downloaded %s\n", filename)

	}

	wg.Done()
}

func proccessErrors(errChan chan extractors.ImageError, wg *sync.WaitGroup) {
	for err := range errChan {
		fmt.Fprintf(os.Stderr, "Failed to get file %s: %v", err.FileURL.String(), err.Err)
	}
	wg.Done()
}

func processThread(urlStr string, maxThreads uint, dir string, extract extractor) error {
	fmt.Printf("Downloading thread %s\n", urlStr)
	url, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return err
	}
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	requiredType := "text/html"
	err = util.ValidateHttpResponse(resp, &requiredType)
	if err != nil {
		return err
	}

	var workerWg sync.WaitGroup
	var errorWg sync.WaitGroup
	files := make(chan extractors.ImageFile, 200)
	errors := make(chan extractors.ImageError, 200)

	var i uint
	for i = 0; i < maxThreads; i++ {
		go imgWorker(files, errors, dir, &workerWg)
		workerWg.Add(1)
	}
	go proccessErrors(errors, &errorWg)
	errorWg.Add(1)

	err = extract(resp.Body, files)
	if err != nil {
		return err
	}

	close(files)
	workerWg.Wait()
	close(errors)
	errorWg.Wait()

	return nil
}

func main() {
	var opts options

	// TODO change getopt lib for better help management
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

	var extractor extractor
	switch opts.Extractor {
	case "4chan":
		extractor = extractors.Extract4Chan
	case "8kun":
		extractor = extractors.Extract8Kun
	default:
		log.Fatalf("Invalid extractor type: %s\n", opts.Extractor)
	}
	threads := args[1:]

	for _, item := range threads {
		err := processThread(item, opts.MaxThreads, opts.Directory, extractor)
		if err != nil {
			log.Printf("Failed to extract url %s: %v", item, err)
		}
	}

	fmt.Printf("Done\n")

}

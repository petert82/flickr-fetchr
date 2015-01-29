package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/petert82/flickr-fetchr/api"
	"github.com/petert82/flickr-fetchr/savr"
	"os"
)

var (
	apiKey     string
	userId     string
	outputFile string
)

func init() {
	flag.StringVar(&apiKey, "api-key", "", "Flickr API key")
	flag.StringVar(&userId, "user-id", "", "Flickr user ID")
	flag.StringVar(&outputFile, "output-file", "", "File to save photo JSON into")
	flag.Parse()
}

// check just dies on error
func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

// checkOptions checks all required command line options were provided
func checkOptions() error {
	if len(apiKey) < 1 {
		return errors.New("api-key option is missing")
	}

	if len(userId) < 1 {
		return errors.New("user-id option is missing")
	}

	if len(outputFile) < 1 {
		return errors.New("output-file option is missing")
	}

	return nil
}

// getPhotoList gets a list of photo-summary data for the user's whole
// photostream
func getPhotoList() ([]api.PhotoSummary, error) {
	getPage := 1
	var ps []api.PhotoSummary

	for {
		p, err := api.GetPhotoSearchPage(getPage, apiKey, userId)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p.Photos...)
		fmt.Printf("Got page %v of %v photos\n", getPage, len(p.Photos))

		if p.Pages == getPage {
			break
		}
		getPage++
	}

	return ps, nil
}

// saveWorker receives photo info on the given channel and saves it to the given
// file
func saveWorker(photos <-chan api.FullPhotoer, f *os.File, done chan<- bool) {
	for p := range photos {
		err := savr.Save(p, f)
		check(err)
		// Move back to erase newline added by savr.Save
		_, err = f.Seek(-1, os.SEEK_CUR)
		check(err)
		_, err = f.WriteString(",")
		check(err)
		fmt.Println("Saved details for photo: ", p.Id())
	}
	done <- true
}

// photoInfoWorker receives photo summaries on the given channel, gets full
// photo info from Flickr and passes it to the saveQueue
func photoInfoWorker(summaries <-chan api.PhotoSummary, saveQueue chan<- api.FullPhotoer, done chan<- bool) {
	for s := range summaries {
		p, err := api.GetPhotoInfo(s.Id(), s.Secret, apiKey)
		check(err)
		fmt.Println("Got details for photo:", p.Id())

		saveQueue <- p
	}
	done <- true
}

func main() {
	check(checkOptions())

	fmt.Println("Fetching photo data from Flickr")
	photos, err := getPhotoList()
	if len(photos) == 0 {
		fmt.Println("Got no photos")
		os.Exit(0)
	}

	f, err := os.Create(outputFile)
	check(err)
	defer f.Close()

	// Our output file is going to contain a JSON array, so write the opening
	// bracket
	_, err = f.WriteString("[")
	check(err)

	// Create a separate goroutine for writing photo details to our output file
	saveQueue := make(chan api.FullPhotoer, 10)
	saveDone := make(chan bool)
	go saveWorker(saveQueue, f, saveDone)

	// Create some workers for fetching our photo details
	fetcherCount := 20
	fetchQueue := make(chan api.PhotoSummary, 500)
	fetchDone := make(chan bool, fetcherCount)
	for w := 1; w <= fetcherCount; w++ {
		go photoInfoWorker(fetchQueue, saveQueue, fetchDone)
	}

	// Send all our photo summaries to the photoInfoWorkers so they can fetch
	// the full details and pass them to the saveQueue
	for _, p := range photos {
		fetchQueue <- p
	}
	close(fetchQueue)
	// Wait for all our fetchers to finish
	for w := 1; w <= fetcherCount; w++ {
		<-fetchDone
	}
	// Wait for saving to finish
	close(saveQueue)
	<-saveDone

	// Get rid of trailing comma and, close our JSON array and sync to disk
	_, err = f.Seek(-1, os.SEEK_CUR)
	check(err)
	_, err = f.WriteString("]\n")
	check(err)
	err = f.Sync()
	check(err)

	fmt.Printf("Saved details for %v photos\n", len(photos))
}

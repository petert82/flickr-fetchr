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
	apiKey string
	userId string
)

func init() {
	flag.StringVar(&apiKey, "api-key", "", "Flickr API key")
	flag.StringVar(&userId, "user-id", "", "Flickr user ID")
	flag.Parse()
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func checkOptions() error {
	if len(apiKey) < 1 {
		return errors.New("api-key option is missing")
	}

	if len(userId) < 1 {
		return errors.New("user-id option is missing")
	}

	return nil
}

func main() {
	check(checkOptions())

	fmt.Println("Fetching photo data from Flickr")
	getPage := 1
	var ps []api.PhotoSummary

	for {
		p, err := api.GetPhotoSearchPage(getPage, apiKey, userId)
		check(err)
		ps = append(ps, p.Photos...)
		fmt.Printf("Got page %v of %v photos\n", getPage, len(p.Photos))

		if p.Pages == getPage {
			break
		}
		getPage++
	}

	if len(ps) == 0 {
		fmt.Println("Got no photos")
		os.Exit(0)
	}

	f, err := os.Create("photos.json")
	check(err)
	defer f.Close()

	// Our output file is going to contain a JSON array, so write the opening
	// bracket
	_, err = f.WriteString("[")
	check(err)

	// Create a separate goroutine for writing photo details to our output file
	queue := make(chan api.FullPhotoer, 10)
	done := make(chan bool)
	go func() {
		for p := range queue {
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
	}()

	// Get details of all the photos in our photostream and pass them to the
	// queue to be saved
	for i := 0; i < 5; i++ {
		p, err := api.GetPhotoInfo(ps[i].Id(), ps[i].Secret, apiKey)
		check(err)
		fmt.Println("Got details for photo:", p.Id())

		queue <- p
	}

	close(queue)
	<-done

	// Get rid of trailing comma and and close our JSON array
	_, err = f.Seek(-1, os.SEEK_CUR)
	check(err)
	_, err = f.WriteString("]\n")
	check(err)
	err = f.Sync()
	check(err)

	fmt.Printf("Saved details for %v photos\n", len(ps))
}

package main

import (
	"fmt"
	"github.com/petert82/flickr-fetchr/api"
	// "io/ioutil"
	"errors"
	"flag"
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
	var photos []api.PhotoSummary

	for {
		p, err := api.GetPhotoSearchPage(getPage, apiKey, userId)
		check(err)
		photos = append(photos, p.Photos...)
		fmt.Printf("Got page %v of %v photos\n", getPage, len(p.Photos))

		if p.Pages == getPage {
			break
		}
		getPage++
	}

	// for _, v := range photos {
	// 	fmt.Println(v.OriginalUrl())
	// 	fmt.Println(v.LargeThumbnailUrl())
	// 	fmt.Println(v.SmallThumbnailUrl())
	// }
	photo, err := api.GetPhotoInfo(photos[0].Id, photos[0].Secret, apiKey)
	check(err)
	fmt.Printf("%+v", photo)

	fmt.Printf("Got %v photos\n", len(photos))
}

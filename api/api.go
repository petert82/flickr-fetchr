package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ApiUrl     = "https://api.flickr.com/services/rest/"
	statusOk   = "ok"
	statusFail = "fail"
)

func flickrRequest(options map[string]string) (*http.Response, error) {
	url, err := url.Parse(ApiUrl)
	if err != nil {
		return nil, err
	}
	q := url.Query()
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")

	for k, v := range options {
		q.Set(k, v)
	}
	url.RawQuery = q.Encode()

	return http.Get(url.String())
}

// GetPhotoSearchPage a single page of photo search results
// https://www.flickr.com/services/api/flickr.photos.search.html
func GetPhotoSearchPage(page int, apiKey, userId string) (*PhotoSearchPage, error) {
	opts := map[string]string{
		"method":         "flickr.photos.search",
		"api_key":        apiKey,
		"user_id":        userId,
		"page":           strconv.Itoa(page),
		"per_page":       "500",
		"privacy_filter": "1", // Only public
		"content_type":   "1", // Only photos
		"extras":         "original_format",
	}
	resp, err := flickrRequest(opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var photos PhotoSearchResult
	err = dec.Decode(&photos)
	if err != nil {
		return nil, err
	}
	switch photos.Status {
	case statusFail:
		return nil, errors.New(fmt.Sprintf("flickr error: %v", photos.Error))
	default:
		return &photos.Photos, nil
	}
}

// GetPhotoInfo gets info for a single photo
// https://www.flickr.com/services/api/flickr.photos.getInfo.html
func GetPhotoInfo(photoId, secret, apiKey string) (*PhotoInfo, error) {
	opts := map[string]string{
		"method":   "flickr.photos.getInfo",
		"api_key":  apiKey,
		"photo_id": photoId,
		"secret":   secret,
	}

	resp, err := flickrRequest(opts)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var photo PhotoInfoResult
	err = dec.Decode(&photo)
	if err != nil {
		return nil, err
	}
	switch photo.Status {
	case statusFail:
		return nil, errors.New(fmt.Sprintf("flickr error: %v", photo.Error))
	default:
		return &photo.Photo, nil
	}
}

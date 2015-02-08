package savr

import (
	"encoding/json"
	"github.com/petert82/flickr-fetchr/api"
	"io"
)

type saveablePhoto struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Urls        map[string]string `json:"urls"`
}

func Save(p api.FullPhotoer, w io.Writer) error {
	jsonPhoto := saveablePhoto{
		Id:          p.Id(),
		Title:       p.Title(),
		Description: p.Description(),
		Urls: map[string]string{
			"original":   p.OriginalUrl(),
			"thumbnailL": p.LargeThumbnailUrl(),
			"thumbnailS": p.SmallThumbnailUrl(),
		},
	}

	enc := json.NewEncoder(w)
	err := enc.Encode(jsonPhoto)
	if err != nil {
		return err
	}

	return nil
}

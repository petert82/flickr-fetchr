package savr

import (
	"encoding/json"
	"github.com/petert82/flickr-fetchr/api"
	"os"
)

type saveablePhoto struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Urls        map[string]string `json:"urls"`
}

func Save(photos []api.FullPhotoer, file string) error {
	jsonPhotos := make([]saveablePhoto, len(photos))

	for i, p := range photos {
		jsonPhotos[i] = saveablePhoto{
			Title:       p.Title(),
			Description: p.Description(),
			Urls: map[string]string{
				"original":   p.OriginalUrl(),
				"thumbnailL": p.LargeThumbnailUrl(),
				"thumbnailS": p.SmallThumbnailUrl(),
			},
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	err = enc.Encode(jsonPhotos)
	if err != nil {
		return err
	}

	f.Sync()

	return nil
}

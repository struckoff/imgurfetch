package imgurfetch

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const albumBaseTpl = "%s/ajaxalbums/getimages/%s/hit.json?all=true"

//Album - information about images in the album.
type Album struct {
	Data struct {
		Images []Image `json:"images"`
	} `json:"data"`
}

//AlbumMeta downloads information about images in the album by album ID.
func AlbumMeta(host, id string) (data Album, err error) {
	if len(id) == 0 {
		return data, errors.New("empty album id")
	}
	url := fmt.Sprintf(albumBaseTpl, host, id)

	res, err := http.Get(url)
	if err != nil {
		return data, errors.Wrap(err, "unable to fetch album json")
	}
	if res.StatusCode >= 400 {
		return data, errors.New(res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return data, errors.Wrap(err, "unable to read res body")
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, errors.Wrap(err, "unable to unmarshal res body")
	}
	return data, nil
}

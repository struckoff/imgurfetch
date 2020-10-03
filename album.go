package imgurfetch

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const albumBaseTpl = "http://imgur.com/ajaxalbums/getimages/%s/hit.json?all=true"

type Album struct {
	Data struct {
		Images []Image `json:"images"`
	} `json:"data"`
}

func AlbumMeta(id string) (data Album, err error) {
	if len(id) == 0 {
		return data, errors.New("empty album id")
	}
	url := fmt.Sprintf(albumBaseTpl, id)

	response, err := http.Get(url)
	if err != nil {
		return data, errors.Wrap(err, "unable to fetch album json")
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return data, errors.Wrap(err, "unable to read response body")
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, errors.Wrap(err, "unable to unmarshal response body")
	}
	return data, nil
}
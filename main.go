package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	BASE_IMAGE_URL = "http://i.imgur.com/"
)

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
}

func getJSON(url string) Album {
	var data Album
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return data
}

func run(url, dir string, sortByResolution bool) {
	var wg sync.WaitGroup
	imgCh := make(chan ImageItem, 10)
	stateCh := make(chan int)

	urlSplit := strings.Split(url, "/")
	url = "http://imgur.com/ajaxalbums/getimages/" + urlSplit[len(urlSplit)-1] + "/hit.json?all=true"
	var resp = getJSON(url)

	go func(stateCh chan int){
		cnt := 0
		for _ = range stateCh {
			cnt++
			fmt.Printf("[ %d/%d ]: %d%%\n", cnt, len(resp.Data.Images), int(float32(cnt)/float32(len(resp.Data.Images)) * 100))
		}
	}(stateCh)

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(imgCh <-chan ImageItem, stateCh chan int) {
			defer wg.Done()
			for img := range imgCh {
				img.get(dir, sortByResolution)
				stateCh <- 1
			}
		}(imgCh, stateCh)
	}

	for _, img := range resp.Data.Images {
		imgCh <- img
	}

	close(imgCh)
	wg.Wait()
}

func main() {
	sortByResolution := flag.Bool("byres", false, "sort by resolution")
	help := flag.Bool("h", false, "sort by resolution")

	flag.Parse()

	if bool(*help) || (len(flag.Args()) < 1){
		fmt.Println("USAGE: imgurfetch [-byres] <url> [path]")
		fmt.Println("-byres - sort images by resolution")
		fmt.Println("-h - show this message")
	} else if len(flag.Args()) >= 2 {
		run(flag.Args()[0], flag.Args()[1], *sortByResolution)
	} else if len(flag.Args()) == 1 {
		run(flag.Args()[0], ".", *sortByResolution)
	}

}

type ImageItem struct {
	Hash   string `json:"hash"`
	Title  string `json:"title"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Ext    string `json:"ext"`
}

func (img *ImageItem) get(imgPath string, sortByResolution bool) {
	response, err := http.Get(BASE_IMAGE_URL + img.Hash + img.Ext)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if sortByResolution {
		imgPath = path.Join(imgPath, fmt.Sprintf("%dx%d", img.Width, img.Height))
	}
	os.MkdirAll(imgPath, 0777)

	err = ioutil.WriteFile(path.Join(imgPath, img.Hash+img.Ext), body, 0644)
	if err != nil {
		panic(err)
	}
}

type Album struct {
	Data struct {
		Images []ImageItem `json:"images"`
	} `json:"data"`
}

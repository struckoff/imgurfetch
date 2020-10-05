package imgurfetch

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

const imageURLTpl = "%s/%s%s"

//ImageWorker - contains information about how to download images
// and where to store them. Worker recieves tasks from "in" channel.
//when task is done it send signal to "done" channel.
//Before each task it asks limiter to get permission for execution.
type ImageWorker struct {
	hostname string //i.imgur.com
	in       <-chan Image
	done     chan<- struct{}
	path     string
	grByRes  bool
	http     *http.Client
	limit    *rate.Limiter
}

//NewWorker create new worker instance.
//If grByRes is true, it will create sub directory WxH.
func NewWorker(host string, in <-chan Image, done chan<- struct{}, path string, grByRes bool, l *rate.Limiter, hc *http.Client) *ImageWorker {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &ImageWorker{
		host,
		in,
		done,
		path,
		grByRes,
		hc,
		l,
	}
}

//Run loop which waits tasks in "in" channel until ctx signals done.
//Before executing tasks it asks limiter for permission.
func (w *ImageWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case img := <-w.in:
			if err := w.limit.Wait(ctx); err != nil {
				if !errors.Is(err, context.Canceled) {
					log.Err(err).Send()
				}
			}
			err := w.imageDownload(img)
			if err != nil {
				log.Err(err).Send()
			}
			w.done <- struct{}{}
		}
	}
}

//imageDownload downloads image and saves it to w.path
//If path is not exist, function will try to create it.
//If flag grByRes is set, it will create sub directory WxH.
func (w *ImageWorker) imageDownload(img Image) error {
	url := fmt.Sprintf(imageUrlTpl, w.hostname, img.Hash, img.Ext)
	res, err := w.http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return errors.New(res.Status)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	ipath := w.path
	if w.grByRes {
		ipath = path.Join(ipath, fmt.Sprintf("%dx%d", img.Width, img.Height))
	}

	err = os.MkdirAll(ipath, 0777)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(ipath, img.Hash+img.Ext), body, 0644)
	if err != nil {
		return err
	}
	return nil
}

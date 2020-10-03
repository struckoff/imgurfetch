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

const imageUrlTpl = "http://i.imgur.com/%s%s"

type ImageWorker struct {
	in      <-chan Image
	done    chan<- struct{}
	path    string
	grByRes bool
	http    *http.Client
	limit *rate.Limiter
}

func NewWorker(in <-chan Image, done chan<- struct{}, path string, grByRes bool, l *rate.Limiter, hc *http.Client) *ImageWorker {
	if hc == nil {
		hc = http.DefaultClient
	}
	return &ImageWorker{
		in,
		done,
		path,
		grByRes,
		hc,
		l,
	}
}

func (w *ImageWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case img := <-w.in:
			if err := w.limit.Wait(ctx); err != nil {
				if !errors.Is(err, context.Canceled){
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

func (w *ImageWorker) imageDownload(img Image) error {
	url := fmt.Sprintf(imageUrlTpl, img.Hash, img.Ext)
	response, err := w.http.Get(url)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(response.Body)
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

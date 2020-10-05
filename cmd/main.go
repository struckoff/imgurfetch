package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/struckoff/imgurfetch"
	"golang.org/x/time/rate"
	"os"
	"os/signal"
	"strings"
	"time"
)

const imageHost = "http://i.imgur.com"
const albumHost = "http://imgur.com"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [arguments] <url> [path(default: .)]\n", os.Args[0])
		flag.PrintDefaults()
	}
	g := flag.Bool("g", false, "group images by resolution")
	w := flag.Int("w", 10, "number of workers")
	r := flag.Duration("r", 0, "rate limit(how often requests could happen)")
	flag.Parse()

	u := flag.Arg(0) //url of an album
	d := flag.Arg(1) //target directory

	if len(u) == 0 {
		fmt.Println("album url is not specified")
		return
	}

	if len(d) == 0 {
		d = "."
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	if err := run(ctx, u, d, *g, *w, *r); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, url, dir string, grByRes bool, wn int, r time.Duration) error {
	lim := rate.NewLimiter(rate.Every(r), 1)

	iCh := make(chan imgurfetch.Image)
	defer close(iCh)

	done := make(chan struct{})
	defer close(done)

	urlSplit := strings.Split(url, "/")
	if len(urlSplit) == 0 {
		return errors.New("album ID not found")
	}

	albID := urlSplit[len(urlSplit)-1]
	alb, err := imgurfetch.AlbumMeta(albumHost, albID)
	if err != nil {
		return err
	}

	for i := 0; i < wn; i++ {
		w := imgurfetch.NewWorker(imageHost, iCh, done, dir, grByRes, lim, nil)
		go w.Run(ctx)
	}

	go func() {
		for _, img := range alb.Data.Images {
			iCh <- img
		}
	}()

	bar := pb.StartNew(len(alb.Data.Images))
	for i := len(alb.Data.Images) - 1; i >= 0; i-- {
		select {
		case <-ctx.Done():
			return nil
		case <-done:
			bar.Increment()
		}
	}
	bar.Finish()

	return nil
}

[![PkgGoDev](https://pkg.go.dev/badge/github.com/struckoff/imgurfetch)](https://pkg.go.dev/github.com/struckoff/imgurfetch)
![Go](https://github.com/struckoff/imgurfetch/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/struckoff/imgurfetch)](https://goreportcard.com/report/github.com/struckoff/imgurfetch)


# Imgurfetch
Download images from album by album URL.

```sh
imgurfetch -h

Usage of imgurfetch:
imgurfetch [arguments] <url> [path(default: .)]
  -g	group images by resolution
  -r duration
    	rate limit(how often requests could happen)
  -w int
    	number of workers (default 10)

```

## Test and build
```sh
make all
```

Will build binary file into $BINARY_FOLDER(defaut:./bin )

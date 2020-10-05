package imgurfetch

import (
	"crypto/rand"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
)

func TestNewWorker(t *testing.T) {
	type args struct {
		hostname string
		in       <-chan Image
		done     chan<- struct{}
		path     string
		grByRes  bool
		l        *rate.Limiter
		hc       *http.Client
	}
	tests := []struct {
		name string
		args args
		want *ImageWorker
	}{
		{
			name: "test",
			args: args{
				hostname: "test-host",
				in:       nil,
				done:     nil,
				path:     "test-path",
				grByRes:  false,
				l:        nil,
				hc:       nil,
			},
			want: &ImageWorker{
				hostname: "test-host",
				in:       nil,
				done:     nil,
				path:     "test-path",
				grByRes:  false,
				http:     http.DefaultClient,
				limit:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWorker(tt.args.hostname, tt.args.in, tt.args.done, tt.args.path, tt.args.grByRes, tt.args.l, tt.args.hc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestImageWorker_imageDownload(t *testing.T) {
	type fields struct {
		grByRes bool
		limit   *rate.Limiter
	}
	type args struct {
		img Image
	}
	type want struct {
		body     []byte
		status   int
		err      bool
		filename string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "no group, no err",
			fields: fields{
				grByRes: false,
				limit:   nil,
			},
			args: args{
				img: Image{
					Hash:   "test-hash",
					Title:  "test-title",
					Width:  240,
					Height: 480,
					Ext:    ".png",
				},
			},
			want: want{
				body:     make([]byte, 250250),
				status:   http.StatusOK,
				err:      false,
				filename: "test-hash.png",
			},
		},
		{
			name: "group, no err",
			fields: fields{
				grByRes: true,
				limit:   nil,
			},
			args: args{
				img: Image{
					Hash:   "test-hash",
					Title:  "test-title",
					Width:  240,
					Height: 480,
					Ext:    ".png",
				},
			},
			want: want{
				body:     make([]byte, 250250),
				status:   http.StatusOK,
				err:      false,
				filename: path.Join("240x480", "test-hash.png"),
			},
		},
		{
			name: "no group, err",
			fields: fields{
				grByRes: true,
				limit:   nil,
			},
			args: args{
				img: Image{
					Hash:   "test-hash",
					Title:  "test-title",
					Width:  240,
					Height: 480,
					Ext:    ".png",
				},
			},
			want: want{
				body:   make([]byte, 250250),
				status: http.StatusInternalServerError,
				err:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := ioutil.TempDir(path.Join("testdata", "temp"), "temp-*")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			if _, err := rand.Read(tt.want.body); err != nil {
				t.Fatal(err)
			}

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.want.status)
				_, _ = w.Write(tt.want.body)
			}))

			w := &ImageWorker{
				hostname: srv.URL,
				path:     dir,
				http:     srv.Client(),
				grByRes:  tt.fields.grByRes,
				limit:    tt.fields.limit,
			}

			err = w.imageDownload(tt.args.img)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.FileExists(t, path.Join(dir, tt.want.filename))
				f, err := os.Open(path.Join(dir, tt.want.filename))
				if err != nil {
					t.Fatal(err)
				}
				got, err := ioutil.ReadAll(f)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, tt.want.body, got)
			}
		})
	}
}

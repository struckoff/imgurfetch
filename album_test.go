package imgurfetch

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAlbumMeta(t *testing.T) {
	type args struct {
		body   []byte
		status int
		id     string
	}
	type want struct {
		err  bool
		data Album
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no err",
			args: args{
				body: []byte(`{
    "data": {
        "count": 2,
        "images": [{
            "hash": "GL7igry",
            "title": "",
            "description": null,
            "has_sound": false,
            "width": 1080,
            "height": 1920,
            "size": 594584,
            "ext": ".png",
            "animated": false,
            "prefer_video": false,
            "looping": false,
            "datetime": "2014-11-26 15:58:34",
            "edited": "0"
        }, {
            "hash": "tLhJOrE",
            "title": "",
            "description": null,
            "has_sound": false,
            "width": 1080,
            "height": 1920,
            "size": 1993181,
            "ext": ".png",
            "animated": false,
            "prefer_video": false,
            "looping": false,
            "datetime": "2014-11-26 15:58:50",
            "edited": "0"
        }],
        "include_album_ads": true
    },
    "success": true,
    "status": 200
}`),
				status: http.StatusOK,
				id:     "test-id",
			},
			want: want{
				err: false,
				data: Album{
					Data: struct {
						Images []Image `json:"images"`
					}{
						Images: []Image{{
							Hash:   "GL7igry",
							Title:  "",
							Width:  1080,
							Height: 1920,
							Ext:    ".png",
						}, {
							Hash:   "tLhJOrE",
							Title:  "",
							Width:  1080,
							Height: 1920,
							Ext:    ".png",
						}},
					},
				},
			},
		}, {
			name: "http err",
			args: args{
				body:   nil,
				status: http.StatusInternalServerError,
				id:     "test-id",
			},
			want: want{
				err: true,
			},
		},
		{
			name: "unmarshal err",
			args: args{
				body:   []byte("Definitely not a valid JSON"),
				status: http.StatusOK,
				id:     "test-id",
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.args.status)
				_, _ = w.Write(tt.args.body)
			}))

			got, err := AlbumMeta(srv.URL, tt.args.id)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.data, got)
			}
		})
	}
}

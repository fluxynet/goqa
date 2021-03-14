package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
	"github.com/fluxynet/goqa/web"
)

type fakecache struct {
	pkg      string
	coverage goqa.Coverage
	keys     []string
}

func (f fakecache) Reset(covs ...goqa.Coverage) error {
	panic("not implemented")
	return nil
}

func (f fakecache) Get(pkg string) (*goqa.Coverage, bool) {
	if f.pkg == pkg {
		return &f.coverage, true
	}

	return nil, false
}

func (f fakecache) Keys() ([]string, error) {
	return f.keys, nil
}

func (f fakecache) Close() error {
	return nil
}

func TestServer_Get(t *testing.T) {
	type fields struct {
		Cache  goqa.Cache
		Prefix string
	}

	type args struct {
		path string
	}

	type want struct {
		status  int
		headers http.Header
		body    string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "not found",
			fields: fields{
				Cache: fakecache{},
			},
			args: args{
				path: "/api/foo",
			},
			want: want{
				status: http.StatusNotFound,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body: `{"error":"resource not found"}`,
			},
		},
		{
			name: "found",
			fields: fields{
				Cache: fakecache{
					pkg: "foo",
					coverage: goqa.Coverage{
						Pkg:        "foo",
						Percentage: 10,
						Time:       "2000-01-01T00:00:00.673068822Z",
					},
				},
			},
			args: args{
				path: "/api/foo",
			},
			want: want{
				status: http.StatusOK,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body: `{"pkg":"foo","percentage":10,"time":"2000-01-01T00:00:00.673068822Z"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Cache:  tt.fields.Cache,
				Prefix: "/api/",
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.args.path, nil)

			s.Get(w, r)

			internal.AssertHttp(t, w, tt.want.status, tt.want.headers, tt.want.body)
		})
	}

}

func TestServer_Index(t *testing.T) {
	type fields struct {
		IndexHTML []byte
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "empty",
		},
		{
			name: "markup",
			fields: fields{
				IndexHTML: []byte(`<html><head></head><body></body></html>`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				IndexHTML: tt.fields.IndexHTML,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			s.Index(w, r)

			internal.AssertHttp(t, w, http.StatusOK, http.Header{"Content-Type": []string{web.ContentTypeHTML}}, string(tt.fields.IndexHTML))
		})
	}
}

func TestServer_List(t *testing.T) {
	type fields struct {
		Cache goqa.Cache
	}

	type want struct {
		status  int
		headers http.Header
		body    string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "empty",
			fields: fields{
				Cache: fakecache{},
			},
			want: want{
				status: http.StatusOK,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body: `[]`,
			},
		},
		{
			name: "non-empty",
			fields: fields{
				Cache: fakecache{
					keys: []string{"foo", "bar", "baz"},
				},
			},
			want: want{
				status: http.StatusOK,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body: `["foo","bar","baz"]`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Cache: tt.fields.Cache,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			s.List(w, r)

			internal.AssertHttp(t, w, tt.want.status, tt.want.headers, tt.want.body)
		})
	}
}

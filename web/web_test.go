package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fluxynet/goqa/internal"
)

func TestPrint(t *testing.T) {
	type args struct {
		status  int
		ctype   string
		content []byte
	}

	type want struct {
		Header http.Header
		Status int
		Body   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no status, no ctype, no content",
			args: args{
				status:  0,
				ctype:   "",
				content: nil,
			},
			want: want{
				Header: nil,
				Status: 0,
				Body:   "",
			},
		},
		{
			name: "no ctype, no content",
			args: args{
				status:  http.StatusOK,
				ctype:   "",
				content: nil,
			},
			want: want{
				Header: nil,
				Status: http.StatusOK,
				Body:   "",
			},
		},
		{
			name: "no status, no content",
			args: args{
				status:  0,
				ctype:   ContentTypeJSON,
				content: nil,
			},
			want: want{
				Header: http.Header{
					"Content-Type": []string{ContentTypeJSON},
				},
				Status: 0,
				Body:   "",
			},
		},
		{
			name: "no status, no ctype",
			args: args{
				status:  0,
				ctype:   "",
				content: []byte(`hello world`),
			},
			want: want{
				Header: nil,
				Status: 0,
				Body:   "hello world",
			},
		},
		{
			name: "http status ok, json",
			args: args{
				status:  http.StatusOK,
				ctype:   ContentTypeJSON,
				content: []byte(`{"message":"hello world"}`),
			},
			want: want{
				Header: http.Header{
					"Content-Type": []string{ContentTypeJSON},
				},
				Status: http.StatusOK,
				Body:   `{"message":"hello world"}`,
			},
		},
		{
			name: "http internal server error, html",
			args: args{
				status:  http.StatusInternalServerError,
				ctype:   ContentTypeHTML,
				content: []byte(`<html><head><title>Foobar</title></head><body><h1>Oops</h1><p>something went wrong</p></body></html>`),
			},
			want: want{
				Header: http.Header{
					"Content-Type": []string{ContentTypeHTML},
				},
				Status: http.StatusInternalServerError,
				Body:   `<html><head><title>Foobar</title></head><body><h1>Oops</h1><p>something went wrong</p></body></html>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w = httptest.NewRecorder()
			Print(w, tt.args.status, tt.args.ctype, tt.args.content)

			internal.AssertHttp(t, w, tt.want.Status, tt.want.Header, tt.want.Body)
		})
	}

}

func TestJson(t *testing.T) {
	type args struct {
		r interface{}
	}

	type want struct {
		Header http.Header
		Status int
		Body   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil",
			args: args{
				r: nil,
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusOK,
				Body:   `null`,
			},
		},
		{
			name: "unmarshallable",
			args: args{
				r: func() {},
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusInternalServerError,
				Body:   `{"error":"json: unsupported type: func()"}`,
			},
		},
		{
			name: "array",
			args: args{
				r: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusOK,
				Body:   `[1,2,3,4,5,6,7,8,9,10]`,
			},
		},
		{
			name: "object",
			args: args{
				r: struct {
					Name      string `json:"name"`
					Age       int
					EmailAddr string `json:"email"`
				}{Name: "John Doe", Age: 20, EmailAddr: "john@doe.com"},
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusOK,
				Body:   `{"name":"John Doe","Age":20,"email":"john@doe.com"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w = httptest.NewRecorder()
			Json(w, tt.args.r)
			internal.AssertHttp(t, w, tt.want.Status, tt.want.Header, tt.want.Body)
		})
	}
}

func TestJsonError(t *testing.T) {
	type args struct {
		status int
		err    error
	}

	type want struct {
		Header http.Header
		Status int
		Body   string
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "nil err",
			args: args{
				status: http.StatusBadRequest,
				err:    nil,
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusBadRequest,
				Body:   `{"error":""}`,
			},
		},
		{
			name: "error",
			args: args{
				status: http.StatusInternalServerError,
				err:    errors.New("something went wrong"),
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusInternalServerError,
				Body:   `{"error":"something went wrong"}`,
			},
		},
		{
			name: "error with quotes",
			args: args{
				status: http.StatusNotFound,
				err:    errors.New(`the item "foo" could not be found`),
			},
			want: want{
				Header: http.Header{"Content-Type": []string{ContentTypeJSON}},
				Status: http.StatusNotFound,
				Body:   `{"error":"the item \"foo\" could not be found"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w = httptest.NewRecorder()
			JsonError(w, tt.args.status, tt.args.err)
			internal.AssertHttp(t, w, tt.want.Status, tt.want.Header, tt.want.Body)
		})
	}
}

func TestReadBody(t *testing.T) {
	type args struct {
		body string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty body",
			args: args{
				body: ``,
			},
			want: ``,
		},
		{
			name: "plaintext body",
			args: args{
				body: `foo=bar`,
			},
			want: `foo=bar`,
		},
		{
			name: "json body",
			args: args{
				body: `{"name":"john"}`,
			},
			want: `{"name":"john"}`,
		},
	}

	for _, tt := range tests {
		for _, method := range []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		} {
			t.Run(tt.name+"_"+method, func(t *testing.T) {
				r, err := http.NewRequest(method, "/", strings.NewReader(tt.args.body))

				if err != nil {
					t.Errorf("new request: %s", err.Error())
					return
				}

				var wantErr bool
				switch method {
				default:
					wantErr = true
				case http.MethodPost, http.MethodPut, http.MethodPatch:
					wantErr = false
				}

				got, err := ReadBody(r)
				if (err != nil) != wantErr {
					t.Errorf("ReadBody() error = %v, wantErr %v", err, wantErr)
					return
				}

				if wantErr {
					if len(got) != 0 {
						t.Errorf("read body is not empty, got = %s", got)
					}
				} else if !bytes.Equal(got, []byte(tt.want)) {
					var replacer = strings.NewReplacer(
						"\r", "[R]", "\n", "[N]",
					)

					var wb = replacer.Replace(tt.want)
					var gb = replacer.Replace(string(got))

					t.Errorf("body not as wanted:\nwant = %s\ngot  = %s", wb, gb)
				}
			})
		}
	}
}

func TestVerifyBody(t *testing.T) {
	type args struct {
		b   []byte
		sig string
		key string
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "good",
			args: args{
				b:   []byte(`The quick brown fox jumps over the lazy dog`),
				sig: "sha1=de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9",
				key: "key",
			},
			want: nil,
		},
		{
			name: "sha256",
			args: args{
				b:   []byte(`The quick brown fox jumps over the lazy dog`),
				sig: "sha256=f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd9",
				key: "key",
			},
			want: ErrPayloadUnverified,
		},
		{
			name: "bad length",
			args: args{
				b:   []byte(`The quick brown fox jumps over the lazy dog`),
				sig: "sha1=de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d9xxx",
				key: "key",
			},
			want: ErrPayloadUnverified,
		},
		{
			name: "bad value",
			args: args{
				b:   []byte(`The quick brown fox jumps over the lazy dog`),
				sig: "sha1=de7c9b85b8b78aa6bc8a7a36f70a90701c9db4d0",
				key: "key",
			},
			want: ErrPayloadUnverified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := VerifyBody(tt.args.b, tt.args.sig, tt.args.key); err != tt.want {
				t.Errorf("VerifyBody() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

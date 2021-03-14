package hook

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fluxynet/goqa"
	"github.com/fluxynet/goqa/internal"
	"github.com/fluxynet/goqa/web"
)

type fakebroker struct {
	event goqa.Event
}

func (f *fakebroker) Listen(ctx context.Context) (<-chan goqa.Event, error) {
	return nil, nil
}

func (f *fakebroker) Publish(ctx context.Context, event goqa.Event) error {
	f.event = event
	return nil
}

func (f *fakebroker) Close() error {
	return nil
}

func TestHook_Receive(t *testing.T) {
	type fields struct {
		SigKey string
	}

	type args struct {
		headers http.Header
		body    string
	}

	type want struct {
		status  int
		headers http.Header
		body    string
		event   goqa.Event
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "empty body",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{},
				body:    ``,
			},
			want: want{
				status: http.StatusBadRequest,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body:  `{"error":"request incomplete"}`,
				event: nil,
			},
		},
		{
			name: "invalid json body",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{
					githubHeaderSignature: []string{"sha1=6bb71066bae0bfdf23ced85ab0193759a073c0d8"},
				},
				body: `hello=world`,
			},
			want: want{
				status: http.StatusBadRequest,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body:  `{"error":"request incomplete"}`,
				event: nil,
			},
		},
		{
			name: "missing header",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{},
				body:    `{"event":"push","repository":"fluxynet/go-test-example","commit":"1320d4f1cf36041e6d34ff45ed8661d5940806db","ref":"refs/heads/master","head":"","workflow":"Go","data":[{"Time":"2021-03-07T23:09:38.67302542Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Output":"    --- PASS: TestSum/1_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67302942Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.67303332Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Output":"    --- PASS: TestSum/2_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67303712Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673040821Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Output":"    --- PASS: TestSum/1_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673044821Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673048321Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Output":"    --- PASS: TestSum/2_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673052121Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673056222Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Output":"    --- PASS: TestSum/4_Negative_/_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673062122Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673065422Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum","Elapsed":0},{"Time":"2021-03-07T23:09:38.673068822Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"PASS\n"},{"Time":"2021-03-07T23:09:38.673072523Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"coverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675133054Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"ok  \tgithub.com/fluxynet/go-test-example\t0.006s\tcoverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675521879Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Elapsed":0.006}]}`,
			},
			want: want{
				status: http.StatusBadRequest,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body:  `{"error":"request incomplete"}`,
				event: nil,
			},
		},
		{
			name: "bad signature",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{
					githubHeaderSignature: []string{"sha1=4e800a4a7df54f147a4c08b2461a29640a6aa6ez"},
				},
				body: `{"event":"push","repository":"fluxynet/go-test-example","commit":"1320d4f1cf36041e6d34ff45ed8661d5940806db","ref":"refs/heads/master","head":"","workflow":"Go","data":[{"Time":"2021-03-07T23:09:38.67302542Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Output":"    --- PASS: TestSum/1_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67302942Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.67303332Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Output":"    --- PASS: TestSum/2_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67303712Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673040821Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Output":"    --- PASS: TestSum/1_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673044821Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673048321Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Output":"    --- PASS: TestSum/2_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673052121Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673056222Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Output":"    --- PASS: TestSum/4_Negative_/_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673062122Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673065422Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum","Elapsed":0},{"Time":"2021-03-07T23:09:38.673068822Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"PASS\n"},{"Time":"2021-03-07T23:09:38.673072523Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"coverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675133054Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"ok  \tgithub.com/fluxynet/go-test-example\t0.006s\tcoverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675521879Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Elapsed":0.006}]}`,
			},
			want: want{
				status: http.StatusBadRequest,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body:  `{"error":"payload could not be verified"}`,
				event: nil,
			},
		},
		{
			name: "good signature no coverages",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{
					githubHeaderSignature: []string{"sha1=c059bcbae7f4ca03f2313bce5ec7692a7e0cba87"},
				},
				body: `{"event":"push","repository":"fluxynet/go-test-example","commit":"1320d4f1cf36041e6d34ff45ed8661d5940806db","ref":"refs/heads/master","head":"","workflow":"Go","data":[{"Time":"2021-03-07T23:09:38.67302542Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Output":"    --- PASS: TestSum/1_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67302942Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.67303332Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Output":"    --- PASS: TestSum/2_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67303712Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673040821Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Output":"    --- PASS: TestSum/1_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673044821Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673048321Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Output":"    --- PASS: TestSum/2_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673052121Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673056222Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Output":"    --- PASS: TestSum/4_Negative_/_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673062122Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673065422Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum","Elapsed":0},{"Time":"2021-03-07T23:09:38.673068822Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"PASS\n"},{"Time":"2021-03-07T23:09:38.675521879Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Elapsed":0.006}]}`,
			},
			want: want{
				status: http.StatusOK,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body:  `{"message":"web hook was not very interesting"}`,
				event: nil,
			},
		},
		{
			name: "good signature with coverages",
			fields: fields{
				SigKey: "foobar",
			},
			args: args{
				headers: http.Header{
					githubHeaderSignature: []string{"sha1=4e800a4a7df54f147a4c08b2461a29640a6aa6ea"},
				},
				body: `{"event":"push","repository":"fluxynet/go-test-example","commit":"1320d4f1cf36041e6d34ff45ed8661d5940806db","ref":"refs/heads/master","head":"","workflow":"Go","data":[{"Time":"2021-03-07T23:09:38.67302542Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Output":"    --- PASS: TestSum/1_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67302942Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.67303332Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Output":"    --- PASS: TestSum/2_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.67303712Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673040821Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Output":"    --- PASS: TestSum/1_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673044821Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/1_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673048321Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Output":"    --- PASS: TestSum/2_Negative (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673052121Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/2_Negative","Elapsed":0},{"Time":"2021-03-07T23:09:38.673056222Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Output":"    --- PASS: TestSum/4_Negative_/_Positive (0.00s)\n"},{"Time":"2021-03-07T23:09:38.673062122Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum/4_Negative_/_Positive","Elapsed":0},{"Time":"2021-03-07T23:09:38.673065422Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Test":"TestSum","Elapsed":0},{"Time":"2021-03-07T23:09:38.673068822Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"PASS\n"},{"Time":"2021-03-07T23:09:38.673072523Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"coverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675133054Z","Action":"output","Package":"github.com/fluxynet/go-test-example","Output":"ok  \tgithub.com/fluxynet/go-test-example\t0.006s\tcoverage: 83.3% of statements\n"},{"Time":"2021-03-07T23:09:38.675521879Z","Action":"pass","Package":"github.com/fluxynet/go-test-example","Elapsed":0.006}]}`,
			},
			want: want{
				status: http.StatusOK,
				headers: http.Header{
					"Content-Type": []string{web.ContentTypeJSON},
				},
				body: `{"message":"web hook well received"}`,
				event: goqa.GithubEvent{
					Event:      "push",
					Repository: "fluxynet/go-test-example",
					Commit:     "1320d4f1cf36041e6d34ff45ed8661d5940806db",
					Ref:        "refs/heads/master",
					Head:       "",
					Workflow:   "Go",
					Coverage: []goqa.Coverage{
						{
							Pkg:        "github.com/fluxynet/go-test-example",
							Percentage: 83,
							Time:       "2021-03-07T23:09:38.673072523Z",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &fakebroker{}
			h := &Hook{
				Broker: b,
				SigKey: tt.fields.SigKey,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.body))
			r.Header = tt.args.headers

			h.Receive(w, r)

			internal.AssertHttp(t, w, tt.want.status, tt.want.headers, tt.want.body)

			if (b.event == nil) != (tt.want.event == nil) {
				t.Errorf("event nil, want = %t, got = %t", tt.want.event == nil, b.event == nil)
			}

			if b.event == nil {
				return
			}

			var we, ge goqa.GithubEvent

			if v, ok := tt.want.event.(goqa.GithubEvent); ok {
				we = v
			} else {
				t.Errorf("want event is not github event")
				return
			}

			if v, ok := b.event.(*goqa.GithubEvent); ok {
				ge = *v
			} else if v, ok := b.event.(goqa.GithubEvent); ok {
				we = v
			} else {
				t.Errorf("got event is not github event")
				return
			}

			internal.AssertGithubEventsEqual(t, &ge, &we)
		})
	}
}

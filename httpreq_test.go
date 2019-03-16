/*
Package httpreq unmarshals values from http.Request as specified in struct tags.
*/
package httpreq

import (
	"fmt"
	"net/http"
	"testing"
)

func reqWithHeaders(t *testing.T, method, url string, headers map[string]string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}

func req(t *testing.T, method, url string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

type T1 struct {
	Hiihaa string `httpreq:"foo,at=query"`
	Haa    int    `httpreq:"haa,at=query"`
}

func TestRecognizeTag(t *testing.T) {
	var t1 T1
	err := Unmarshal(req(t, "GET", "/bar?foo=value&haa=a10"), t1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", t1)
}

func TestUnmarshal(t *testing.T) {
	var t1 T1
	var t1p *T1
	type args struct {
		req *http.Request
		v   interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			"ok",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				&t1,
			},
			nil,
		},
		{
			"v is nil",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				t1p,
			},
			ErrIsNil,
		},
		{
			"v is not a pointer",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				t1,
			},
			ErrIsNotPointer,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.req, tt.args.v); err != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
Package httpreq unmarshals values from http.Request as specified in struct tags.
*/
package httpreq

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
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
	Hiihaa   string `req:"foo,at=query"`
	Haa      int    `req:"haa,at=query"`
	skipThis int
}
type T2 struct {
	hiihaa string `req:"foo,at=query"`
	Haa    int    `req:"haa,at=query"`
}
type T3 struct {
	Hiihaa string `req:"foo"`
}
type T4 struct {
	Hiihaa string `req:"foo,at=bar"`
}

func TestUnmarshal(t *testing.T) {
	var t1 T1
	var t1p *T1
	var t2 T2
	var t3 T3
	var t4 T4
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
		{
			"cannot set field value",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				&t2,
			},
			ErrCannotSet,
		},
		{
			"invalid tag",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				&t3,
			},
			ErrInvalidTag,
		},
		{
			"invalid at",
			args{
				req(t, "GET", "/?foo=srtval&haa=10"),
				&t4,
			},
			ErrInvalidAt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.req, tt.args.v); errors.Cause(err) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

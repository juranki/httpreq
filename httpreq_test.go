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
	err := Unmarshal(req(t, "GET", "/bar?foo=value&haa=10"), &t1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", t1)
}

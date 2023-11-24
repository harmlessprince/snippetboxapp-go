package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	/// Initialize a new httptest.ResponseRecorder.
	rr := httptest.NewRecorder()

	// / Initialize a new dummy http.Request
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//call the ping handlare function passing in the ResponseRecorder and http.Request
	ping(rr, r)

	// Call the Result() method on the http.ResponseRecorder to get the
	// http.Response generated by the ping handler.
	response := rr.Result()

	if response.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, response.StatusCode)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}

func TestPingHandler(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	StatusCode, _, body := ts.get(t, "/ping")

	if StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, StatusCode)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
func TestShowSnippet(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())

	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("The boy and the king")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := ts.get(t, test.urlPath)
			if code != test.wantCode {
				t.Errorf("want %d; got %d", test.wantCode, code)
			}
			if !bytes.Contains(body, test.wantBody) {
				t.Errorf("want body to contain %q", test.wantBody)
			}
		})
	}
}
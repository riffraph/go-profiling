package demo

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

/**
 * test by invoking the handler method directly
 */
func TestHandleHi_Recorder(t *testing.T) {
	recorder := httptest.NewRecorder()
	handleHi(recorder, createRequest(t, "GET / HTTP/1.0\r\n\r\n"))

	if !strings.Contains(recorder.Body.String(), "visitor number") {
		t.Errorf("Unexpected output: %s", recorder.Body)
	}
}

func createRequest(t *testing.T, reqStr string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(reqStr)))
	if err != nil {
		t.Fatal(err)
	}

	return req
}

/**
 * test by starting a http server and binding the handler to it
 */
func TestHandleHi_TestServer(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handleHi))
	defer testServer.Close()

	res, err := http.Get(testServer.URL)
	if err != nil {
		t.Error(err)
		return
	}

	if contentType, contentStr := res.Header.Get("Content-Type"), "text/html; charset=utf-8"; contentType != contentStr {
		t.Error("Content-Type = %q; want %q", contentType, contentStr)
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("Got: %s", body)
}

func TestHandleHi_TestServer_Parallel(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handleHi))
	defer testServer.Close()

	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			res, err := http.Get(testServer.URL)
			if err != nil {
				t.Error(err)
				return
			}

			if contentType, contentStr := res.Header.Get("Content-Type"), "text/html; charset=utf-8"; contentType != contentStr {
				t.Error("Content-Type = %q; want %q", contentType, contentStr)
			}

			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			if err != nil {
				t.Error(err)
				return
			}

			t.Logf("Got: %s", body)
		}()
	}

	wg.Wait()
}

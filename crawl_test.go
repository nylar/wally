package wally

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestCrawl_GrabUrl(t *testing.T) {
	data := []byte("some data")
	status := 200
	ts := Handler(status, data)
	defer ts.Close()

	d, err := GrabUrl(ts.URL)
	assert.Equal(t, d, data)
	assert.NoError(t, err)
}

func TestCrawl_Crawler(t *testing.T) {
	DatabaseRebuild(session)

	data := []byte("really cool stuff")
	ts := Handler(200, data)
	defer ts.Close()

	err := Crawler(ts.URL, session)
	assert.NoError(t, err)
}

func TestCrawl_CrawlerNoURL(t *testing.T) {
	DatabaseRebuild(session)

	err := Crawler("", session)
	assert.Error(t, err)
}

package wally

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSource struct {
	Data  []byte
	Error error
}

func (ms *MockSource) Grab(resource string) ([]byte, error) {
	return ms.Data, ms.Error
}

func Handler(status int, data []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write(data)
	}))
}

func TestCrawl_GrabResource(t *testing.T) {
	data := []byte("Hello, World")
	ms := MockSource{
		Data:  data,
		Error: nil,
	}

	d, err := GrabResource(&ms, "mydoc.html")
	assert.Equal(t, d, data)
	assert.NoError(t, err)
}

func TestCrawl_GrabResourceWebSource(t *testing.T) {
	data := []byte("Hello, World")
	ts := Handler(200, data)
	defer ts.Close()

	ws := new(WebSource)
	d, err := ws.Grab(ts.URL)
	assert.Equal(t, data, d)
	assert.NoError(t, err)
}

func TestCrawl_GrabResourceWebSourceWithError(t *testing.T) {
	data := []byte("Hello, World")
	ts := Handler(200, data)
	defer ts.Close()

	ws := new(WebSource)
	_, err := ws.Grab("")
	assert.Error(t, err)
}

func TestCrawl_Crawler(t *testing.T) {
	data := []byte("really cool stuff")

	ms := MockSource{
		Data:  data,
		Error: nil,
	}

	err := Crawler("hello.html", session, &ms)
	assert.NoError(t, err)
}

func TestCrawl_CrawlerWithError(t *testing.T) {
	data := []byte("really cool stuff")

	ms := MockSource{
		Data:  data,
		Error: errors.New("Failed to get resource"),
	}

	err := Crawler("hello.html", session, &ms)
	assert.Error(t, err)
}

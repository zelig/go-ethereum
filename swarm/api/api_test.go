package api

import (
	// "bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

func testApi(t *testing.T, f func(*Api)) {
	datadir, err := ioutil.TempDir("", "bzz-test")
	if err != nil {
		t.Fatalf("unable to create temp dir: %v", err)
	}
	os.RemoveAll(datadir)
	defer os.RemoveAll(datadir)
	dpa, err := storage.NewLocalDPA(datadir)
	if err != nil {
		return
	}
	api := NewApi(dpa, nil)
	dpa.Start()
	f(api)
	dpa.Stop()
}

type testResponse struct {
	reader storage.SectionReader
	*Response
}

func checkResponse(t *testing.T, resp *testResponse, exp *Response) {

	if resp.MimeType != exp.MimeType {
		t.Errorf("incorrect mimeType. expected '%s', got '%s'", exp.MimeType, resp.MimeType)
	}
	if resp.Status != exp.Status {
		t.Errorf("incorrect status. expected '%d', got '%d'", exp.Status, resp.Status)
	}
	if resp.Size != exp.Size {
		t.Errorf("incorrect size. expected '%d', got '%d'", exp.Size, resp.Size)
	}
	if resp.reader != nil {
		content := make([]byte, resp.Size)
		read, _ := resp.reader.Read(content)
		if int64(read) != exp.Size {
			t.Errorf("incorrect content length. expected '%s...', got '%s...'", read, exp.Size)
		}
		resp.Content = string(content)
	}
	if resp.Content != exp.Content {
		// if !bytes.Equal(resp.Content, exp.Content) {
		t.Errorf("incorrect content. expected '%s...', got '%s...'", string(exp.Content), string(resp.Content))
	}
}

// func expResponse(content []byte, mimeType string, status int) *Response {
func expResponse(content string, mimeType string, status int) *Response {
	return &Response{mimeType, status, int64(len(content)), content}
}

// func testGet(t *testing.T, api *Api, bzzhash string) *testResponse {
func testGet(t *testing.T, api *Api, bzzhash string) *testResponse {
	reader, mimeType, status, err := api.Get(bzzhash, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return &testResponse{reader, &Response{mimeType, status, reader.Size(), ""}}
	// return &testResponse{reader, &Response{mimeType, status, reader.Size(), nil}}
}

func TestApiPut(t *testing.T) {
	testApi(t, func(api *Api) {
		content := "hello"
		exp := expResponse(content, "text/plain", 0)
		// exp := expResponse([]byte(content), "text/plain", 0)
		bzzhash, err := api.Put(content, exp.MimeType)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp := testGet(t, api, bzzhash)
		checkResponse(t, resp, exp)
	})
}

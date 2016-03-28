package jaak

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	port              = "1322"
	trackID           = "789"
	streamerEtherAddr = "0xabc"
)

func TestPlay(t *testing.T) {
	StartHttpServer(&Jaak{}, port)
	resp, err := http.PostForm(
		"http://localhost:"+port+"/play",
		url.Values{
			"trackID":           {trackID},
			"streamerEtherAddr": {streamerEtherAddr},
		})
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected no error, got '%v'", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("expected no error, got '%v'", err)
	}

	expBody := "\"success\""

	if string(body) != expBody {
		t.Fatalf("expected body \n'%v', got \n'%v'", expBody, string(body))
	}

}

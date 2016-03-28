package jaak

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	// "sync"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/logger"
	"github.com/ethereum/go-ethereum/logger/glog"
)

// browser API for registering bzz url scheme handlers:
// https://developer.mozilla.org/en/docs/Web-based_protocol_handlers
// electron (chromium) api for registering bzz url scheme handlers:
// https://github.com/atom/electron/blob/master/docs/api/protocol.md

// starts up http server
func StartHttpServer(jaak *Jaak, port string) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, jaak)
	})
	go http.ListenAndServe(":"+port, serveMux)
	glog.V(logger.Info).Infof("[JAAK] Jaak HTTP proxy started on localhost:%s", port)
}

func handler(w http.ResponseWriter, r *http.Request, jaak *Jaak) {
	requestURL := r.URL
	// This is wrong
	//  if requestURL.Host == "" {
	//    var err error
	//    requestURL, err = url.Parse(r.Referer() + requestURL.Error())
	//    if err != nil {
	//      http.Error(w, err.Error(), http.StatusBadRequest)
	//      return
	//    }
	//  }
	glog.V(logger.Debug).Infof("[JAAK]: HTTP request URL: '%s', Host: '%s', Path: '%s', Referer: '%s', Accept: '%s'", r.Method, r.RequestURI, requestURL.Host, requestURL.Path, r.Referer(), r.Header.Get("Accept"))
	action := requestURL.Path
	w.Header().Set("Content-Type", "text/json")

	// HTTP-based URL protocol handler
	switch action {

	// POST /upload
	case "/upload":
		if r.Method != "POST" {
			http.Error(w, "Method "+r.Method+" is not supported for /upload.", http.StatusMethodNotAllowed)
		}
		artistName := strings.TrimSpace(r.FormValue("artistName"))
		if len(artistName) == 0 {
			http.Error(w, "Form field artistName cannot be blank", http.StatusBadRequest)
		}
		title := strings.TrimSpace(r.FormValue("title"))
		if len(title) == 0 {
			http.Error(w, "Form field title cannot be blank", http.StatusBadRequest)
		}
		durationV := strings.TrimSpace(r.FormValue("duration"))
		if len(durationV) == 0 {
			http.Error(w, "Form field duration cannot be blank", http.StatusBadRequest)
		}
		duration, _ := strconv.Atoi(durationV)
		etherAddrV := strings.TrimSpace(r.FormValue("etherAddr"))
		if len(etherAddrV) == 0 {
			http.Error(w, "Form field etherAddr cannot be blank", http.StatusBadRequest)
		}
		etherAddr := common.HexToAddress(etherAddrV)

		trackMeta := &TrackMeta{
			ArtistName: artistName,
			Title:      title,
			Duration:   uint(duration),
			EtherAddr:  etherAddr,
		}
		artworkFile, _, err := r.FormFile("artwork")
		if err != nil {
			http.Error(w, "Error parsing artwork from multipart post data: "+err.Error(), http.StatusBadRequest)
		}
		defer artworkFile.Close()

		trackFile, _, err := r.FormFile("track")
		if err != nil {
			http.Error(w, "Error parsing track from multipart post data: "+err.Error(), http.StatusBadRequest)
		}
		defer trackFile.Close()

		jaakReceipt, err := jaak.Upload(trackMeta, artworkFile, trackFile)
		if err != nil {
			http.Error(w, "Unable to upload: "+err.Error(), http.StatusBadRequest)
		}
		data, err := json.Marshal(jaakReceipt)
		if err != nil {
			http.Error(w, "Error encoding jaak receipt json: "+err.Error(), http.StatusBadRequest)
		}
		resp := bytes.NewReader(data)
		http.ServeContent(w, r, "", time.Now(), resp)

		// GET /tracks
	case "/tracks":
		if r.Method != "GET" {
			http.Error(w, "Method "+r.Method+" is not supported for /tracks.", http.StatusMethodNotAllowed)
		}
		tracks := jaak.GetTracks()
		data, err := json.Marshal(tracks)
		if err != nil {
			http.Error(w, "Error encoding tracks listing json: "+err.Error(), http.StatusBadRequest)
		}
		resp := bytes.NewReader(data)
		http.ServeContent(w, r, "", time.Now(), resp)

		// POST /play
	case "/play":
		if r.Method != "POST" {
			http.Error(w, "Method "+r.Method+" is not supported for /play.", http.StatusMethodNotAllowed)
		}

		trackIDV := strings.TrimSpace(r.FormValue("trackID"))
		if len(trackIDV) == 0 {
			http.Error(w, "Form field trackID cannot be blank", http.StatusBadRequest)
		}
		trackID := common.HexToHash(trackIDV)

		streamerEtherAddrV := strings.TrimSpace(r.FormValue("streamerEtherAddr"))
		if len(streamerEtherAddrV) == 0 {
			http.Error(w, "Form field streamerEtherAddr cannot be blank", http.StatusBadRequest)
		}
		streamerEtherAddr := common.HexToAddress(streamerEtherAddr)
		play := &Play{
			TrackID:           trackID,
			StreamerEtherAddr: streamerEtherAddr,
		}
		playResp := jaak.Play(play)

		data, err := json.Marshal(playResp)
		if err != nil {
			http.Error(w, "Error encoding tracks listing json: "+err.Error(), http.StatusBadRequest)
		}
		resp := bytes.NewReader(data)
		http.ServeContent(w, r, "", time.Now(), resp)

	default:
		http.Error(w, "Action  "+action+" is not supported.", http.StatusMethodNotAllowed)
	}

}

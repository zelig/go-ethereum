package jaak

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/logger"
	"github.com/ethereum/go-ethereum/logger/glog"
	swarm "github.com/ethereum/go-ethereum/swarm/api"
)

type Track struct {
	TrackMeta
	ArtworkUrl    string      `json:"artworkUrl"`    //: "https://i1.sndcdn.com/artworks-000151009654-2e7g5m-large.jpg",
	StreamUrl     string      `json:"streamUrl"`     // : "https://api.soundcloud.com/tracks/252003956/stream",
	PlaybackCount uint        `json:"playbackCount"` //: 64372,
	TrackID       common.Hash `json:"trackID"`       //   : "6c0cd61f98e127153c7828eb51a9f5eca13307c1"
}

type TrackMeta struct {
	Title      string         `json:"title"`      //: "Lorem Ipsum",
	ArtistName string         `json:"artistName"` //: "Lorem Ipsum",
	Duration   uint           `json:"duration"`   //: 195988,
	EtherAddr  common.Address `json:"etherAddr"`  //   : "6c0cd61f98e127153c7828eb51a9f5eca13307c1"
}

type Play struct {
	TrackID           common.Hash    `json:"trackID"`           //
	StreamerEtherAddr common.Address `json:"streamerEtherAddr"` //
}

type Jaak struct {
	swarm *swarm.Api
}

func (self *Jaak) Play(p *Play) string {
	// increments playcount
	// sends out money to artist according to splits
	glog.V(logger.Info).Infof("[JAAK] Jaak pay&play:%v ", p)
	return "success"
}

func (self *Jaak) Upload(p *TrackMeta, artwork, track io.Reader) (*Track, error) {
	return nil, nil
}

func (self *Jaak) GetTracks() []*Track {
	return nil
}

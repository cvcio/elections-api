package twitter

import (
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/sirupsen/logrus"
)

// Demux receives channels or interfaces and type switches them to call the appropriate handle function.
type Demux interface {
	Handle(m interface{})
	HandleChan(c <-chan interface{})
}

// StreamDemux receives messages and type switches them to call functions with typed messages.
type StreamDemux struct {
	All        func(m interface{})
	Tweet      func(tweet anaconda.Tweet)
	Event      func(event anaconda.Event)
	EventTweet func(event anaconda.EventTweet)
	Other      func(m interface{})
}

// NewAPI creates a new anaconda instance.
// Anaconda is a Twitter API Drivers (github.com/ChimeraCoder/anaconda).
func NewAPI(accesstoken string, accesstokensecret string, consumerkey string, consumersecret string) (*anaconda.TwitterApi, error) {
	api := anaconda.NewTwitterApiWithCredentials(
		accesstoken,
		accesstokensecret,
		consumerkey,
		consumersecret,
	)
	// log.Print(api.Credentials)
	if _, err := api.VerifyCredentials(); err != nil {
		log.Errorf("Bad Authorization Tokens. Please refer to https://apps.twitter.com/ for your Access Tokens: %s", err)
		return nil, err
	}
	return api, nil
}

// NewStreamDemux initializes a new StreamDemux.
func NewStreamDemux() StreamDemux {
	return StreamDemux{
		All:        func(m interface{}) {},
		Tweet:      func(tweet anaconda.Tweet) {},
		Event:      func(event anaconda.Event) {},
		EventTweet: func(event anaconda.EventTweet) {},
		Other:      func(m interface{}) {},
	}
}

// Handle handles messages.
func (d StreamDemux) Handle(m interface{}) {
	d.All(m)

	switch t := m.(type) {
	case anaconda.Tweet:
		d.Tweet(t)
	case anaconda.Event:
		d.Event(t)
	case anaconda.EventTweet:
		d.EventTweet(t)
	default:
		d.Other(t)
	}
}

// HandleChan handles channels.
func (d StreamDemux) HandleChan(c <-chan interface{}) {
	for m := range c {
		d.Handle(m)
	}
}

// Listen struct.
type Listen struct {
	TwitterAPI *anaconda.TwitterApi
}

// NewListener return a new Listener service, given a twitter api client
func NewListener(tw *anaconda.TwitterApi) *Listen {
	s := new(Listen)
	s.TwitterAPI = tw
	return s
}

// Listen start the listener and send cathed urls to chan
func (s *Listen) Listen(ids []string, track []string, urlChan chan anaconda.Tweet) {
	stream := s.TwitterAPI.PublicStreamFilter(url.Values{
		"track":  track,
		"follow": ids,
	})
	demux := NewStreamDemux()
	demux.Tweet = func(t anaconda.Tweet) {
		urlChan <- t
	}
	go demux.HandleChan(stream.C)
}

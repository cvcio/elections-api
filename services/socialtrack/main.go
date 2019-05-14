package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/elections-api/pkg/config"
	proto "github.com/cvcio/elections-api/pkg/proto"
	"github.com/cvcio/elections-api/pkg/twitter"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

func main() {
	time.Sleep(4 * time.Second)
	// ========================================
	// Configure
	cfg := config.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("main: Error loading config: %s", err.Error())
	}

	// Create the gRPC Service
	// Parse Server Options
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	grpcConnection, err := grpc.Dial("localhost:50050", opts...)
	if err != nil {
		log.Debugf("main: GRPC Streamer did not connect: %v", err)
	}
	defer grpcConnection.Close()

	// Create gRPC Chat Client
	streamer := proto.NewTwitterClient(grpcConnection)

	session, err := streamer.Connect(context.Background(), &proto.Session{Id: primitive.NewObjectID().Hex(), Type: "listener"})
	if err != nil {
		log.Debugf("Can't join streamer %s", err.Error())
	}

	// Connect to Stream
	stream, err := streamer.Stream(context.Background())
	if err != nil {
		log.Debugf("Stream Connection Failed: %v", err)
	}
	api, _ := twitter.NewAPI(
		cfg.Twitter.TwitterAccessToken,
		cfg.Twitter.TwitterAccessTokenSecret,
		cfg.Twitter.TwitterConsumerKey,
		cfg.Twitter.TwitterConsumerSecret,
	)

	// Create a new Listener service, with our twitter stream and the scrape service grpc conn
	svc := twitter.NewListener(api)

	// Create a channel to send catched tweets
	tweetChan := make(chan anaconda.Tweet)

	// start Listening the twitter stream
	go svc.Listen(nil, []string{"#Αλεξη_πες_μας", "Ευρωεκλογές2019", "Τσίπρας", "Μητσοτάκης"}, tweetChan)

	// ========================================
	// Shutdown
	//
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// Listen forever on channels
	for {
		select {
		case t := <-tweetChan:
			log.Debugf("New Tweet: %s | %v", t.User.ScreenName, t.IdStr)
			tweet, _ := json.Marshal(&t)
			go stream.Send(&proto.Message{Session: session, Tweet: string(tweet)})
		case <-osSignals:
			os.Exit(1)
		}
	}
}

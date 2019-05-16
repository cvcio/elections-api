package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/cvcio/elections-api/models/nodes"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/es"
	proto "github.com/cvcio/elections-api/pkg/proto"
	"github.com/cvcio/elections-api/pkg/twitter"
	"github.com/kelseyhightower/envconfig"
	"github.com/olivere/elastic"
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

	// =========================================================================
	// Start elasticsearch
	log.Info("main: Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.Port, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("main: Register Elasticsearch: %v", err)
	}

	log.Info("main: Connected to Elasticsearch")

	// Create the gRPC Service
	// Parse Server Options
	var grpcOptions []grpc.DialOption
	grpcOptions = append(grpcOptions, grpc.WithInsecure())

	// Create gRPC Streamer Connection
	grpcStreamerConnection, err := grpc.Dial("localhost:50050", grpcOptions...)
	if err != nil {
		log.Debugf("main: GRPC Streamer did not connect: %v", err)
	}
	defer grpcStreamerConnection.Close()

	// Create gRPC Streamer Client
	streamer := proto.NewTwitterClient(grpcStreamerConnection)

	session, err := streamer.Connect(context.Background(), &proto.Session{Id: primitive.NewObjectID().Hex(), Type: "listener"})
	if err != nil {
		log.Debugf("Can't join streamer %s", err.Error())
	}

	// Connect to Stream
	stream, err := streamer.Stream(context.Background())
	if err != nil {
		log.Debugf("Stream Connection Failed: %v", err)
	}

	// Create gRPC Classification Connection
	grpcClassificationConnection, err := grpc.Dial("localhost:50051", grpcOptions...)
	if err != nil {
		log.Debugf("main: GRPC Classification did not connect: %v", err)
	}
	defer grpcClassificationConnection.Close()

	// Create gRPC Classification Client
	classification := proto.NewClassificationClient(grpcClassificationConnection)

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
	var ids []string
	for _, username := range cfg.Streamer.Follow {
		if u := svc.GetUsersShow(username); u != nil {
			ids = append(ids, u.IdStr)
		}
	}
	go svc.Listen(ids, cfg.Streamer.Track, tweetChan)

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
			c := classifyNestedTweet(esClient, &t, classification)
			tweet, _ := json.Marshal(&c)

			log.Infof("New Tweet From %s", t.User.ScreenName)

			go SaveTweet(esClient, &t)
			go stream.Send(&proto.Message{Session: session, Tweet: string(tweet)})

		case <-osSignals:
			os.Exit(1)
		}
	}
}

// SaveTweet on ES Index
func SaveTweet(esClient *elastic.Client, t *anaconda.Tweet) {
	_, err := esClient.Index().
		Index("mediawatch_twitter_elections_tweets").
		Type("document").
		Id(t.IdStr).
		BodyJson(t).
		Do(context.Background())
	if err != nil {
		log.Errorf("Can't Save Tweet: %s", err.Error())
		return
	}
}

// SaveUser on ES Index
func SaveUser(esClient *elastic.Client, u *nodes.ESUserObj) {
	u.CrawledAt = time.Now().Format(time.RubyDate)
	_, err := esClient.Index().
		Index("mediawatch_twitter_elections_users").
		Type("document").
		Id(primitive.NewObjectID().Hex()).
		BodyJson(u).
		Do(context.Background())
	if err != nil {
		log.Errorf("Can't Save User: %s", err.Error())
		return
	}
}

func classifyNestedTweet(esClient *elastic.Client, t *anaconda.Tweet, c proto.ClassificationClient) *nodes.UserObj {
	var tuF, quF, ruF *proto.UserFeatures
	var tuC, quC, ruC *proto.UserClass

	user := &nodes.UserObj{
		Id:              t.User.Id,
		IdStr:           t.User.IdStr,
		CreatedAt:       t.User.CreatedAt,
		Verified:        t.User.Verified,
		ScreenName:      t.User.ScreenName,
		Name:            t.User.Name,
		Description:     t.User.Description,
		FollowersCount:  t.User.FollowersCount,
		FriendsCount:    t.User.FriendsCount,
		ListedCount:     t.User.ListedCount,
		StatusesCount:   t.User.StatusesCount,
		FavouritesCount: t.User.FavouritesCount,
		ProfileImage:    t.User.ProfileImageUrlHttps,
		BannerImage:     t.User.ProfileBannerURL,
	}

	// Tweet User
	tuF = getUserFeatures(&t.User)
	tuC, _ = c.Classify(context.Background(), tuF)

	user.Metrics = tuF
	user.UserClass = tuC.GetLabel()
	user.UserClassScore = tuC.GetScore()

	go SaveUser(esClient, &nodes.ESUserObj{
		Id:              t.User.Id,
		IdStr:           t.User.IdStr,
		CreatedAt:       t.User.CreatedAt,
		Verified:        t.User.Verified,
		ScreenName:      t.User.ScreenName,
		Name:            t.User.Name,
		Description:     t.User.Description,
		FollowersCount:  t.User.FollowersCount,
		FriendsCount:    t.User.FriendsCount,
		ListedCount:     t.User.ListedCount,
		StatusesCount:   t.User.StatusesCount,
		FavouritesCount: t.User.FavouritesCount,
		ProfileImage:    t.User.ProfileImageUrlHttps,
		BannerImage:     t.User.ProfileBannerURL,
		UserClass:       user.UserClass,
		UserClassScore:  user.UserClassScore,
	})

	// Quoted User
	if t.QuotedStatus != nil {
		user.QuotedStatus = &nodes.UserObj{
			Id:              t.QuotedStatus.User.Id,
			IdStr:           t.QuotedStatus.User.IdStr,
			CreatedAt:       t.QuotedStatus.User.CreatedAt,
			Verified:        t.QuotedStatus.User.Verified,
			ScreenName:      t.QuotedStatus.User.ScreenName,
			Name:            t.QuotedStatus.User.Name,
			Description:     t.QuotedStatus.User.Description,
			FollowersCount:  t.QuotedStatus.User.FollowersCount,
			FriendsCount:    t.QuotedStatus.User.FriendsCount,
			ListedCount:     t.QuotedStatus.User.ListedCount,
			StatusesCount:   t.QuotedStatus.User.StatusesCount,
			FavouritesCount: t.QuotedStatus.User.FavouritesCount,
			ProfileImage:    t.QuotedStatus.User.ProfileImageUrlHttps,
			BannerImage:     t.QuotedStatus.User.ProfileBannerURL,
		}
		quF = getUserFeatures(&t.QuotedStatus.User)
		quC, _ = c.Classify(context.Background(), quF)

		user.QuotedStatus.Metrics = quF
		user.QuotedStatus.UserClass = quC.GetLabel()
		user.QuotedStatus.UserClassScore = quC.GetScore()

		go SaveUser(esClient, &nodes.ESUserObj{
			Id:              t.QuotedStatus.User.Id,
			IdStr:           t.QuotedStatus.User.IdStr,
			CreatedAt:       t.QuotedStatus.User.CreatedAt,
			Verified:        t.QuotedStatus.User.Verified,
			ScreenName:      t.QuotedStatus.User.ScreenName,
			Name:            t.QuotedStatus.User.Name,
			Description:     t.QuotedStatus.User.Description,
			FollowersCount:  t.QuotedStatus.User.FollowersCount,
			FriendsCount:    t.QuotedStatus.User.FriendsCount,
			ListedCount:     t.QuotedStatus.User.ListedCount,
			StatusesCount:   t.QuotedStatus.User.StatusesCount,
			FavouritesCount: t.QuotedStatus.User.FavouritesCount,
			ProfileImage:    t.QuotedStatus.User.ProfileImageUrlHttps,
			BannerImage:     t.QuotedStatus.User.ProfileBannerURL,
			UserClass:       user.QuotedStatus.UserClass,
			UserClassScore:  user.QuotedStatus.UserClassScore,
		})
	}
	// Retweeted User
	if t.RetweetedStatus != nil {
		user.RetweetedStatus = &nodes.UserObj{
			Id:              t.RetweetedStatus.User.Id,
			IdStr:           t.RetweetedStatus.User.IdStr,
			CreatedAt:       t.RetweetedStatus.User.CreatedAt,
			Verified:        t.RetweetedStatus.User.Verified,
			ScreenName:      t.RetweetedStatus.User.ScreenName,
			Name:            t.RetweetedStatus.User.Name,
			Description:     t.RetweetedStatus.User.Description,
			FollowersCount:  t.RetweetedStatus.User.FollowersCount,
			FriendsCount:    t.RetweetedStatus.User.FriendsCount,
			ListedCount:     t.RetweetedStatus.User.ListedCount,
			StatusesCount:   t.RetweetedStatus.User.StatusesCount,
			FavouritesCount: t.RetweetedStatus.User.FavouritesCount,
			ProfileImage:    t.RetweetedStatus.User.ProfileImageUrlHttps,
			BannerImage:     t.RetweetedStatus.User.ProfileBannerURL,
		}
		ruF = getUserFeatures(&t.RetweetedStatus.User)
		ruC, _ = c.Classify(context.Background(), ruF)

		user.RetweetedStatus.Metrics = ruF
		user.RetweetedStatus.UserClass = ruC.GetLabel()
		user.RetweetedStatus.UserClassScore = ruC.GetScore()

		go SaveUser(esClient, &nodes.ESUserObj{
			Id:              t.RetweetedStatus.User.Id,
			IdStr:           t.RetweetedStatus.User.IdStr,
			CreatedAt:       t.RetweetedStatus.User.CreatedAt,
			Verified:        t.RetweetedStatus.User.Verified,
			ScreenName:      t.RetweetedStatus.User.ScreenName,
			Name:            t.RetweetedStatus.User.Name,
			Description:     t.RetweetedStatus.User.Description,
			FollowersCount:  t.RetweetedStatus.User.FollowersCount,
			FriendsCount:    t.RetweetedStatus.User.FriendsCount,
			ListedCount:     t.RetweetedStatus.User.ListedCount,
			StatusesCount:   t.RetweetedStatus.User.StatusesCount,
			FavouritesCount: t.RetweetedStatus.User.FavouritesCount,
			ProfileImage:    t.RetweetedStatus.User.ProfileImageUrlHttps,
			BannerImage:     t.RetweetedStatus.User.ProfileBannerURL,
			UserClass:       user.RetweetedStatus.UserClass,
			UserClassScore:  user.RetweetedStatus.UserClassScore,
		})
	}
	return user
}

func getUserFeatures(u *anaconda.User) *proto.UserFeatures {
	userCreatedAt, _ := time.Parse(time.RubyDate, u.CreatedAt)
	delta := time.Now().Sub(userCreatedAt)
	dates := delta.Hours() / 24

	user := &proto.UserFeatures{
		Followers: int64(u.FollowersCount),
		Friends:   int64(u.FriendsCount),
		Statuses:  u.StatusesCount,
		Favorites: int64(u.FavouritesCount),
		Lists:     u.ListedCount,
		Dates:     dates,
		Actions:   float64((u.StatusesCount + int64(u.FavouritesCount))) / dates,
		Ffr:       0,
		Stfv:      0,
		Fstfv:     0,
	}

	if u.FriendsCount > 0 {
		user.Ffr = float64(u.FollowersCount / u.FriendsCount)
	}

	if u.FavouritesCount > 0 {
		user.Stfv = float64(u.StatusesCount / int64(u.FavouritesCount))
	}

	if (u.StatusesCount + int64(u.FavouritesCount)) > 0 {
		user.Fstfv = float64(int64(u.FollowersCount) / (u.StatusesCount + int64(u.FavouritesCount)))
	}

	return user
}

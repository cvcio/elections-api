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
			c := classifyNestedTweet(&t, classification)
			tweet, _ := json.Marshal(&c)
			log.Info(c)
			go stream.Send(&proto.Message{Session: session, Tweet: string(tweet)})
			// Save Enriched
		case <-osSignals:
			os.Exit(1)
		}
	}
}

// UserObj ..
type UserObj struct {
	Id              int64               `json:"id"`
	IdStr           string              `json:"id_str"`
	CreatedAt       string              `json:"created_at"`
	Verified        bool                `json:"verified"`
	ScreenName      string              `json:"screen_name"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	ProfileImage    string              `json:"profile_image_url_https"`
	BannerImage     string              `json:"profile_banner_url"`
	UserClass       string              `json:"user_class"`
	UserClassScore  float64             `json:"user_class_score"`
	QuotedStatus    *UserObj            `json:"quoted_status"`
	RetweetedStatus *UserObj            `json:"retweeted_status"`
	Metrics         *proto.UserFeatures `json:"metrics"`
}

func classifyNestedTweet(t *anaconda.Tweet, c proto.ClassificationClient) *UserObj {
	var tuF, quF, ruF *proto.UserFeatures
	var tuC, quC, ruC *proto.UserClass

	user := &UserObj{
		Id:           t.User.Id,
		IdStr:        t.User.IdStr,
		CreatedAt:    t.User.CreatedAt,
		Verified:     t.User.Verified,
		ScreenName:   t.User.ScreenName,
		Name:         t.User.Name,
		Description:  t.User.Description,
		ProfileImage: t.User.ProfileImageUrlHttps,
		BannerImage:  t.User.ProfileBannerURL,
	}

	// Tweet User
	tuF = getUserFeatures(&t.User)
	tuC, _ = c.Classify(context.Background(), tuF)

	user.Metrics = tuF
	user.UserClass = tuC.GetLabel()
	user.UserClassScore = tuC.GetScore()

	// Quoted User
	if t.QuotedStatus != nil {
		user.QuotedStatus = &UserObj{
			Id:           t.QuotedStatus.User.Id,
			IdStr:        t.QuotedStatus.User.IdStr,
			CreatedAt:    t.QuotedStatus.User.CreatedAt,
			Verified:     t.QuotedStatus.User.Verified,
			ScreenName:   t.QuotedStatus.User.ScreenName,
			Name:         t.QuotedStatus.User.Name,
			Description:  t.QuotedStatus.User.Description,
			ProfileImage: t.QuotedStatus.User.ProfileImageUrlHttps,
			BannerImage:  t.QuotedStatus.User.ProfileBannerURL,
		}
		quF = getUserFeatures(&t.QuotedStatus.User)
		quC, _ = c.Classify(context.Background(), quF)

		user.QuotedStatus.Metrics = quF
		user.QuotedStatus.UserClass = quC.GetLabel()
		user.QuotedStatus.UserClassScore = quC.GetScore()
	}
	// Retweeted User
	if t.RetweetedStatus != nil {
		user.RetweetedStatus = &UserObj{
			Id:           t.RetweetedStatus.User.Id,
			IdStr:        t.RetweetedStatus.User.IdStr,
			CreatedAt:    t.RetweetedStatus.User.CreatedAt,
			Verified:     t.RetweetedStatus.User.Verified,
			ScreenName:   t.RetweetedStatus.User.ScreenName,
			Name:         t.RetweetedStatus.User.Name,
			Description:  t.RetweetedStatus.User.Description,
			ProfileImage: t.RetweetedStatus.User.ProfileImageUrlHttps,
			BannerImage:  t.RetweetedStatus.User.ProfileBannerURL,
		}
		ruF = getUserFeatures(&t.RetweetedStatus.User)
		ruC, _ = c.Classify(context.Background(), ruF)

		user.RetweetedStatus.Metrics = ruF
		user.RetweetedStatus.UserClass = ruC.GetLabel()
		user.RetweetedStatus.UserClassScore = ruC.GetScore()
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

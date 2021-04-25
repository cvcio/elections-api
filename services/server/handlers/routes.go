package handlers

import (
	"io"
	"net/http"

	"github.com/cvcio/elections-api/pkg/auth"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/cvcio/elections-api/pkg/mailer"
	"github.com/cvcio/elections-api/pkg/middleware"
	"github.com/olivere/elastic"

	"github.com/cvcio/elections-api/pkg/redis"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitter"
	"gopkg.in/olahol/melody.v1"

	proto "github.com/cvcio/elections-api/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Broadcast Receives Messages from the Streaming Service and Broadcasts back to websocket
func Broadcast(stream proto.Twitter_StreamClient, m *melody.Melody) {
	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Errorf("RECEIVE FROM STREAMER ERROR: %s", err.Error())
			return
		}
		log.Infof("RECEIVE FROM STREAMER TWEET %s", rec.Tweet)
		m.Broadcast([]byte(rec.Tweet))
	}
}

// API : Returns an new API
func API(cfg *config.Config, db *db.DB, es *elastic.Client, authenticator *auth.Authenticator, mail *mailer.Mailer, streamer *redis.Service /*proto.TwitterClient*/) http.Handler {
	m := melody.New()
	/*
		session, err := streamer.Connect(context.Background(), &proto.Session{Id: primitive.NewObjectID().Hex(), Type: "api"})
		if err != nil {
			log.Debugf("Can't join streamer %s", err.Error())
		}

		// Connect to Stream
		stream, err := streamer.Stream(context.Background())
		if err != nil {
			log.Debugf("Stream Connection Failed: %v", err)
		}
	*/
	app := gin.Default()
	app.RedirectTrailingSlash = true
	app.RedirectFixedPath = true

	// authmw is used for authentication/authorization middleware.
	authmw := middleware.Auth{
		Authenticator: authenticator,
	}

	app.Use(middleware.Logger(log.StandardLogger(), true))

	if cfg.Env == "development" {
		app.Use(middleware.EnableCORS("*"))
	} else {
		gin.SetMode(gin.ReleaseMode)
		app.Use(middleware.EnableCORS(" *." + cfg.Web.DomainName))
	}

	goth.UseProviders(
		twitter.New(
			cfg.Twitter.TwitterConsumerKey,
			cfg.Twitter.TwitterConsumerSecret,
			cfg.Twitter.TwitterAuthCallBack,
		),
	)

	users := &Users{
		cfg:           cfg,
		db:            db,
		authenticator: authenticator,
		mail:          mail,
	}

	metrics := &Metrics{
		cfg: cfg,
		db:  db,
		es:  es,
	}

	annotations := &Annotations{
		cfg: cfg,
		db:  db,
		es:  es,
	}

	authRoutes := app.Group("/api/auth")
	{
		authRoutes.GET("/:provider", users.OAuthTwitter)
		authRoutes.GET("/:provider/callback", users.OAuthTwitterCB)
	}
	app.GET("/v2/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	public := app.Group("/v2")
	{
		public.POST("/users", users.Update)
		public.POST("/users/verify", users.Verify)
		public.POST("/users/2fa", users.SendPin)
		public.POST("/users/token", users.Token)

		public.GET("/annotate", annotations.GetRandom)
		public.GET("/metrics/user/:id/volume", metrics.GetVolumeByUser)
		public.GET("/metrics/user/:id/count", metrics.CountByUser)
	}
	private := app.Group("/v2")
	{
		private.Use(authmw.Authenticate())
		private.POST("/annotate", annotations.Create)
	}
	// Forbid Access
	// This is usefull when you combine multiple microservices
	app.NoRoute(func(c *gin.Context) {
		c.String(http.StatusForbidden, "Access Forbidden")
		c.Abort()
	})

	tweets := make(chan []byte)
	err := streamer.Subscribe("tweets", tweets)
	if err != nil {
		log.Fatal(err)
	}

	go func(tweets chan []byte, m *melody.Melody) {
		for {
			tweet := <-tweets
			log.Info(string(tweet))
			m.Broadcast(tweet)
		}
	}(tweets, m)
	/*
		if session != nil {
			go func(stream proto.Twitter_StreamClient, session *proto.Session) {
				err := stream.Send(&proto.Message{Session: session})
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Infof("proto.Twitter_StreamClient -> Error: %s", err.Error())
				}
				return
			}(stream, session)
			go Broadcast(stream, m)
		}
	*/

	return app
}

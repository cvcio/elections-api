package handlers

import (
	"context"
	"net/http"

	"github.com/cvcio/elections-api/pkg/auth"
	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/cvcio/elections-api/pkg/mailer"
	"github.com/cvcio/elections-api/pkg/middleware"
	"github.com/olivere/elastic"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitter"
	"gopkg.in/olahol/melody.v1"

	"go.mongodb.org/mongo-driver/bson/primitive"

	proto "github.com/cvcio/elections-api/pkg/proto"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Broadcast Receives Messages from the Streaming Service and Broadcasts back to websocket
func Broadcast(stream proto.Twitter_StreamClient, m *melody.Melody) {
	for {
		rec, _ := stream.Recv()
		m.Broadcast([]byte(rec.Tweet))
	}
}

// API : Returns an new API
func API(cfg *config.Config, db *db.DB, es *elastic.Client, authenticator *auth.Authenticator, mail *mailer.Mailer, streamer proto.TwitterClient) http.Handler {
	m := melody.New()

	session, err := streamer.Connect(context.Background(), &proto.Session{Id: primitive.NewObjectID().Hex(), Type: "api"})
	if err != nil {
		log.Debugf("Can't join streamer %s", err.Error())
	}

	// Connect to Stream
	stream, err := streamer.Stream(context.Background())
	if err != nil {
		log.Debugf("Stream Connection Failed: %v", err)
	}

	go stream.Send(&proto.Message{Session: session})
	go Broadcast(stream, m)

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
		app.Use(middleware.EnableCORS("*." + cfg.Web.DomainName))
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
		public.GET("/p/annotate", annotations.GetRandom)
		public.GET("/metrics/user/:id/volume", metrics.GetVolumeByUser)
	}
	private := app.Group("/v2")
	{
		private.Use(authmw.Authenticate())
		private.GET("/annotate", annotations.GetRandom)
		private.POST("/annotate", annotations.Create)
	}
	// Forbid Access
	// This is usefull when you combine multiple microservices
	app.NoRoute(func(c *gin.Context) {
		c.String(http.StatusForbidden, "Access Forbidden")
		c.Abort()
	})

	return app
}

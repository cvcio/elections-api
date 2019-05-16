package handlers

import (
	"context"
	"net/http"

	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/db"
	"github.com/cvcio/elections-api/pkg/mailer"
	"github.com/cvcio/elections-api/pkg/middleware"
	"github.com/cvcio/elections-api/pkg/auth"
	"github.com/plagiari-sm/mediawatch/pkg/es"

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
func API(cfg *config.Config, db *db.DB, es *es.ES, authenticator *auth.Authenticator, mail *mailer.Mailer, streamer proto.TwitterClient) http.Handler {
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
		es:            es,
		authenticator: authenticator,
		mail:          mail,
	}

	authRoutes := app.Group("/api/auth")
	{
		authRoutes.GET("/:provider", users.OAuthTwitter)
		authRoutes.GET("/:provider/callback", users.OAuthTwitterCB)
	}
	app.GET("/v2/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})
	app.POST("/v2/users", users.Update)
	app.POST("/v2/users/verify", users.Verify)
	app.POST("/v2/users/2fa", users.SendPin)
	app.POST("/v2/users/token", users.Token)

	app.Use(authmw.Authenticate())
	{
		
	}

	// Forbid Access
	// This is usefull when you combine multiple microservices
	app.NoRoute(func(c *gin.Context) {
		c.String(http.StatusForbidden, "Access Forbidden")
		c.Abort()
	})

	return app
}

package main

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cvcio/elections-api/pkg/config"
	"github.com/cvcio/elections-api/pkg/es"
	"github.com/cvcio/elections-api/services/server/handlers"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"github.com/plagiari-sm/mediawatch/pkg/auth"
	"github.com/plagiari-sm/mediawatch/pkg/db"
	"github.com/plagiari-sm/mediawatch/pkg/mailer"
	"github.com/cvcio/elections-api/pkg/redis"
	log "github.com/sirupsen/logrus"
)

func main() {
	// ========================================
	// Configure
	cfg := config.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("main: Error loading config: %s", err.Error())
	}

	// Configure logger
	// Default level for this example is info, unless debug flag is present
	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = log.InfoLevel
		log.Error(err.Error())
	}
	log.SetLevel(level)

	// Adjust logging format
	log.SetFormatter(&log.JSONFormatter{})
	if cfg.Log.Dev {
		log.SetFormatter(&log.TextFormatter{})
	}

	log.Debug("main: Starting")
	// ============================================== ==============
	// Start Mongo
	log.Debug("main: Initialize Mongo")
	dbConn, err := db.New(cfg.MongoURL(), cfg.Mongo.Path, cfg.Mongo.DialTimeout)
	if err != nil {
		log.Fatalf("main: Register DB: %v", err)
	}
	log.Debug("main: Connected to Mongo")
	defer dbConn.Close()

	// =========================================================================
	// Start elasticsearch
	log.Debug("main: Initialize Elasticsearch")
	esClient, err := es.NewElasticsearch(cfg.Elasticsearch.Host, cfg.Elasticsearch.Port, cfg.Elasticsearch.User, cfg.Elasticsearch.Pass)
	if err != nil {
		log.Fatalf("main: Register Elasticsearch: %v", err)
	}

	log.Debug("main: Connected to Elasticsearch")

	// =========================================================================
	// Find auth keys
	keyContents, err := ioutil.ReadFile(cfg.Auth.PrivateKeyFile)
	if err != nil {
		log.Fatalf("main: Reading auth private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		log.Fatalf("main: Parsing auth private key: %v", err)
	}

	publicKeyLookup := auth.NewSingleKeyFunc(cfg.Auth.KeyID, key.Public().(*rsa.PublicKey))

	authenticator, err := auth.NewAuthenticator(key, cfg.Auth.KeyID, cfg.Auth.Algorithm, publicKeyLookup)
	if err != nil {
		log.Fatalf("main: Constructing authenticator: %v", err)
	}

	log.Debug("main: Created auth keys")

	// Create mail service"github.com/sirupsen/logrus"
	mail := mailer.New(
		cfg.SMTP.Server,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Pass,
		cfg.SMTP.From,
		cfg.SMTP.FromName,
		cfg.SMTP.Reply,
	)

	log.Debug("main: Created mail service")

	log.Debug("main: Connecting to Redis")
	// Queue setup
	pubsub, err := redis.New(cfg.Redis.Host)
	if err != nil {
		log.Fatal("main: Register Redis Pub/Sub: %v", err.Error())
	}
	log.Debug("main: Connected to Redis")
	/*
		// Create the gRPC Service
		// Parse Server Options
		var grpcOptions []grpc.DialOption
		grpcOptions = append(grpcOptions, grpc.WithInsecure())

		grpcConnection, err := grpc.Dial(fmt.Sprintf("%s:%s", cfg.Streamer.Host, cfg.Streamer.Port), grpcOptions...)
		if err != nil {
			log.Debugf("main: GRPC Streamer did not connect: %v", err)
		}
		defer grpcConnection.Close()

		// Create gRPC Chat Client
		streamer := proto.NewTwitterClient(grpcConnection)
	*/

	// ========================================
	// Create a server

	// create the http.Server
	api := http.Server{
		Addr: cfg.Web.Host + cfg.Web.Port,
		Handler: handlers.API(
			cfg,
			dbConn,
			esClient,
			authenticator,
			mail,
			pubsub,
			// streamer,
		),
		ReadTimeout:    cfg.Web.ReadTimeout,
		WriteTimeout:   cfg.Web.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// ========================================
	// Blocking main listening for requests
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	log.Debug("main: Ready to start")
	go func() {
		log.Infof("main: Starting api Listening %s%s", cfg.Web.Host, cfg.Web.Port)
		serverErrors <- api.ListenAndServe()
	}()

	// ========================================
	// Shutdown
	//
	// Listen for os signals
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)

	// =========================================================================
	// Stop API Service
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("main: Error starting server: %v", err)

	case <-osSignals:
		log.Info("main: Start shutdown...")

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		if err := api.Shutdown(ctx); err != nil {
			log.Infof("main: Graceful shutdown did not complete in %v: %v", cfg.Web.ShutdownTimeout, err)
			if err := api.Close(); err != nil {
				log.Fatalf("main: Could not stop http server: %v", err)
			}
		}
	}
}

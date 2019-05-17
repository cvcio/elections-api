package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cvcio/elections-api/pkg/config"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	proto "github.com/cvcio/elections-api/pkg/proto"
)

// Client Connction
type Client struct {
	Type    string
	Id      string
	Channel chan proto.Message
}

// TwitterHandler ...
type TwitterHandler struct {
	cachedSessions map[string]Client
	// A Mutex is a mutual exclusion lock.
	// The zero value for a Mutex is an unlocked mutex.
	mu sync.Mutex
	// connections
	in  int64
	out int64
}

func (s *TwitterHandler) withReadLock(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f()
}

func (s *TwitterHandler) withWriteLock(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f()
}

// Listen ...
func (s *TwitterHandler) Listen(stream proto.Twitter_StreamServer, ch chan<- proto.Message) {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			close(ch)
			return
		}
		if err != nil {
			return
		}
		ch <- *req
	}
}

// Connect to streamer Endpoint
func (s *TwitterHandler) Connect(ctx context.Context, req *proto.Session) (*proto.Session, error) {
	log.Infof("CONNECTION FROM %s %s", req.Type, req.Id)

	s.withWriteLock(func() {
		s.cachedSessions[req.Id] = Client{
			Type:    req.Type,
			Id:      req.Id,
			Channel: make(chan proto.Message),
		}
	})
	s.in++
	return req, nil
}

// Filter ...
func (s *TwitterHandler) Filter(sender string, m proto.Message, stream proto.Twitter_StreamServer) {
	//Lock locks m. If the lock is already in use, the calling goroutine
	// blocks until the mutex is available.
	s.mu.Lock()
	// Unlock unlocks m. It is a run-time error if m is not locked
	// on entry to Unlock.
	defer s.mu.Unlock()

	for _, receiver := range s.cachedSessions {
		if sender != receiver.Type {
			log.Debugf("STREAMING TO %s->%s %v", sender, receiver.Type, receiver)
			receiver.Channel <- m
		}
	}
}

// Stream Endpoint
func (s *TwitterHandler) Stream(stream proto.Twitter_StreamServer) error {
	rec, err := stream.Recv()
	if err != nil {
		log.Debugf("Failed to Stream: %v", err)
		return err
	}

	log.Infof("STREAMING FROM %s %s", rec.Session.Type, rec.Session.Id)
	// s.cachedSessions[rec.Session.Id] = make(chan proto.Message)

	// Non-Blocking Client Messages Channel
	messagesChannel := make(chan proto.Message, 100)
	go s.Listen(stream, messagesChannel)

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case clientMessage := <-messagesChannel:
			go s.Filter(clientMessage.Session.Type, clientMessage, stream)
		case channelMessage := <-s.cachedSessions[rec.Session.Id].Channel:
			stream.Send(&channelMessage)
		}
	}
}

func main() {
	time.Sleep(0 * time.Second)
	// ========================================
	// Configure
	cfg := config.New()

	// Read config from env variables
	err := envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("main: Error loading config: %s", err.Error())
	}

	// Get local network address to listen on
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Streamer.Host, cfg.Streamer.Port))
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	// Init shutdown listener
	ch := make(chan os.Signal, 5)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()

	// Create the gRPC Service
	// Parse Server Options
	// Create grpc server
	var opts []grpc.ServerOption
	svc := grpc.NewServer(opts...)
	// Register Service Handlers
	proto.RegisterTwitterServer(svc, &TwitterHandler{
		cachedSessions: make(map[string]Client),
	})

	log.Printf("Starting gRPC Server on: %s:%s", cfg.Streamer.Host, cfg.Streamer.Port)

	// Register reflection service on gRPC server.
	//
	// gRPC Server Reflection provides information about publicly-accessible
	// gRPC services on a server, and assists clients at runtime
	// to construct RPC requests and responses without precompiled service information.
	// It is used by gRPC CLI, which can be used to introspect server protos
	// and send/receive test RPCs.
	//
	// https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md
	reflection.Register(svc)
	// Serve gRPC Service with Error
	errChanSVC := make(chan error, 10)

	go func() {
		errChanSVC <- svc.Serve(listen)
	}()

	signalChanSVC := make(chan os.Signal, 1)
	signal.Notify(signalChanSVC, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChanSVC:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChanSVC:
			log.Println(fmt.Sprintf("Captured message %v. Exiting...", s))
			os.Exit(0)
		}
	}
}

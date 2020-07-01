package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server - Server parameters
type Server struct {
	router      *mux.Router
	httpServer  http.Server
	httpHandler http.Handler
	appConfig   *AppConfig
	mongoClient *mongo.Client

	killServer chan int
	connClose  chan int
}

// NewServer - Creates a new instance of the Server
func NewServer(appConfig *AppConfig) *Server {
	s := &Server{
		appConfig:  appConfig,
		router:     mux.NewRouter(),
		killServer: make(chan int),
		connClose:  make(chan int),
	}

	s.router.StrictSlash(false)

	return s
}

// Initialize - Initializes the Server
func (s *Server) Initialize() {
	s.connectMongoDB()
	s.addRoutes()
	s.addMiddleware()
	s.addHandlers()

	go func() {
		signit := make(chan os.Signal, 1)
		signal.Notify(signit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

		select {
		case <-signit:
			log.Info("Received System Interrupt")
		case <-s.killServer:
			log.Info("Received Kill Request")

		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		log.Info("Closing open connections and shutting down server")

		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			log.WithError(err).Error("Error shutting down server")
		}

		close(s.connClose)
	}()

}

func (s *Server) connectMongoDB() {
	ctx := context.Background()

	var err error

	s.mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(appConfig.MongoDBUri))
	if err != nil {
		log.WithError(err).Panic("Unable to connect to mongoDB")
	}

	ctxT, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err = s.mongoClient.Ping(ctxT, readpref.Primary()); err != nil {
		log.WithError(err).Panic("Ping mongoDB failed")
	}

	log.Info("Connected to mongoDB")
}

func (s *Server) addRoutes() {
	//s.router.HandleFunc("/")
}

func (s *Server) addMiddleware() {
	//s.router.Use()
}

func (s *Server) addHandlers() {
	s.httpHandler = s.router
}

// Listen - Starts the HTTP Server
func (s *Server) Listen() {
	addr := fmt.Sprintf("%s:%d", "", s.appConfig.PORT)
	s.httpServer = http.Server{
		Addr:         addr,
		Handler:      s.httpHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Infof("Server started at: %s", addr)

	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.WithError(err).Error("error starting server")
	}
}

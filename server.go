package main

import (
	"fmt"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Server struct {
	PORT    int
	Address string
	config  *AppConfig
	router  *mux.Router
	logger  *logrus.Logger
	DB      *gorm.DB
}

func NewServer(logger *logrus.Logger, config *AppConfig) *Server {
	return &Server{
		logger:  logger,
		PORT:    config.PORT,
		Address: config.Address,
		router:  mux.NewRouter(),
		config:  config,
	}
}

func (s *Server) Initialize() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Kolkata",
		s.config.DBHost,
		s.config.DBUser,
		s.config.DBPassword,
		s.config.DBName,
		s.config.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.logger.WithError(err).Error("Failed to connect to database")
	}

	s.logger.Infof("Connected to Database: %s", dsn)

	s.DB = db

	s.addRoute()
}

func (s *Server) addRoute() {
	s.router.HandleFunc("/", hello).Methods(http.MethodGet)
}

func (s *Server) Listen() {
	addr := fmt.Sprintf("%s:%d", s.Address, s.PORT)

	s.logger.Infof("Starting server on: %s", addr)

	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		panic(err)
	}
}

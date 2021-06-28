package main

import (
	"fmt"
	"go-demo/home"
	rt "go-demo/template"
	"go-demo/users"
	"html/template"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type Server struct {
	PORT      int
	Address   string
	config    *AppConfig
	router    *mux.Router
	logger    *logrus.Logger
	DB        *gorm.DB
	templates *template.Template
}

func NewServer(logger *logrus.Logger, config *AppConfig, templates *template.Template) *Server {
	return &Server{
		logger:    logger,
		PORT:      config.PORT,
		Address:   config.Address,
		router:    mux.NewRouter(),
		config:    config,
		templates: templates,
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

	//s.DB.Logger

	err = runMigration(s.logger, s.config)
	if err != nil {
		s.logger.WithError(err).Panic("Failed to run migrations")
	}

	s.addRoute()
}

func (s *Server) addRoute() {
	s.router.StrictSlash(true)

	r := rt.NewRenderer(s.logger, s.templates)

	ur := users.NewRepository(s.DB)
	uh := users.NewHandler(s.logger, ur, r)
	hh := home.NewHandler(s.logger, r)

	s.router.PathPrefix("/static").Handler(http.FileServer(http.FS(staticFiles)))

	s.router.HandleFunc("/", hh.Home).Methods(http.MethodGet)
	s.router.HandleFunc("/submit", hh.Submit).Methods(http.MethodPost)

	s.router.HandleFunc("/user", uh.User).Methods(http.MethodGet)
	s.router.HandleFunc("/users", uh.Create).Methods(http.MethodPost)
	s.router.HandleFunc("/users", uh.GetAll).Methods(http.MethodGet)
	s.router.HandleFunc("/users/{id}", uh.Update).Methods(http.MethodPost)
	s.router.HandleFunc("/users/{id}", uh.Delete).Methods(http.MethodDelete)
	s.router.HandleFunc("/users/{id}", uh.Get).Methods(http.MethodGet)

}

func (s *Server) Listen() {
	addr := fmt.Sprintf("%s:%d", s.Address, s.PORT)

	s.logger.Infof("Starting server on: %s", addr)

	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		panic(err)
	}
}

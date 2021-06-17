package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	config := NewAppConfig()
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)

	s := NewServer(logger, config)
	s.Initialize()
	s.Listen()
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

package main

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

//go:embed views/*
var htmlFiles embed.FS

//go:embed static/*
var staticFiles embed.FS

func main() {
	config := NewAppConfig()
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)

	templates, err := template.New("view/*.html").Funcs(template.FuncMap{
		"upper": func(str string) string { return strings.ToUpper(str) },
	}).ParseFS(htmlFiles, "views/*.html")
	if err != nil {
		logger.WithError(err).Panic("error loading templates")
	}

	s := NewServer(logger, config, templates)
	s.Initialize()
	s.Listen()
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

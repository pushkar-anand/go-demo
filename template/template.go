package template

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Renderer struct {
	logger    *logrus.Logger
	templates *template.Template
}

func NewRenderer(logger *logrus.Logger, templates *template.Template) *Renderer {
	return &Renderer{logger: logger, templates: templates}
}

func (r *Renderer) Render(w http.ResponseWriter, templateName string, data interface{}) {
	w.Header().Add("Content-Type", "text/html")
	err := r.templates.ExecuteTemplate(w, templateName, data)
	if err != nil {
		r.logger.WithError(err).Error("error loading submit.html")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

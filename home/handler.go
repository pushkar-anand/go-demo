package home

import (
	"fmt"
	"go-demo/template"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger   *logrus.Logger
	renderer *template.Renderer
}

func NewHandler(logger *logrus.Logger, renderer *template.Renderer) *Handler {
	return &Handler{logger: logger, renderer: renderer}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"time": time.Now().Format(time.RFC822Z),
	}

	w.Header().Add("Set-Cookie", "user=1")

	h.renderer.Render(w, "home.html", data)
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.WithError(err).Error("error parsing form")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	if name == "" {
		http.Error(w, "Name is empty", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("user")
	if err != nil {
		h.logger.WithError(err).Error("")
	}

	h.logger.Debugf("Cookie: %s", cookie.Value)

	data := map[string]interface{}{
		"name":   name,
		"secret": true,
		"list": []numbers{
			{Value: 1},
			{Value: 2},
			{Value: 3},
			{Value: 4},
		},
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("name=%s", name))

	h.renderer.Render(w, "submit.html", data)
}

type numbers struct {
	Value int
}

package users

import (
	"encoding/json"
	"errors"
	"go-demo/template"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger     *logrus.Logger
	repository *Repository
	renderer   *template.Renderer
}

func NewHandler(logger *logrus.Logger, repository *Repository, renderer *template.Renderer) *Handler {
	return &Handler{logger: logger, repository: repository, renderer: renderer}
}

func (h *Handler) User(w http.ResponseWriter, r *http.Request) {
	h.renderer.Render(w, "user.html", nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.WithError(err).Error("error in reading form data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	email := r.Form.Get("email")
	ageStr := r.Form.Get("age")

	age, err := strconv.Atoi(ageStr)
	if err != nil {
		http.Error(w, "Age must be an integer", http.StatusBadRequest)
		return
	}

	user := &User{
		Name:  &name,
		Email: &email,
		Age:   &age,
	}

	err = h.repository.Create(user)
	if err != nil {
		h.logger.WithError(err).Error("error in creating users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User Created"))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	pathVariables := mux.Vars(r)
	idStr := pathVariables["id"]

	h.logger.Debugf("ID: %v", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}

	user, err := h.repository.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "User doesn't exists", http.StatusNotFound)
		return
	}

	if err != nil {
		h.logger.WithError(err).Error("error in fetching users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {

	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.repository.GetAll()
	if err != nil {
		h.logger.WithError(err).Error("error in fetching users")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Users": users,
	}

	h.renderer.Render(w, "home.html", data)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Read json request to a user variable
	// 2. Call the repository Delete function

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	// TODO:
	// 1. Read json request to a user variable
	// 2. Call the repository Update function

}

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	users []User
}

var logger = log.New(os.Stdout, "car-renting-api ", 0)

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		users:  []User{},
	}
	s.Routes()
	return s
}

func (s *Server) Routes() {
	s.HandleFunc("/add-user", s.addUser).Methods(http.MethodPost)
	s.HandleFunc("/get-user", s.getUser).Methods(http.MethodGet)
}

func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("add user called")
	var temp User
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Println("temp", temp)
	newId := uuid.New()
	temp.ID = newId.String()
	s.users = append(s.users, temp)
	w.Header().Set("Content-Type", "Application/json")
	logger.Println("added user", temp)
	if err := json.NewEncoder(w).Encode(temp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("uid")
	var temp User
	for _, user := range s.users {
		if user.ID == id {
			temp = user
			break
		}
	}
	w.Header().Set("Content-Type", "Application/json")
	if err := json.NewEncoder(w).Encode(temp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

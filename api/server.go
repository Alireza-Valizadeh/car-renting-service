package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)
type Server struct {
	*mux.Router
}
func findOptimalValue(old string, newS string) (output string) {
	if len(newS) > 0 {
		output = newS
	} else {
		output = old
	}
	return
}
var logger = log.New(os.Stdout, "car-renting-api ", 0)
var db = NewDB()
func handleError(w http.ResponseWriter, message string, code int) {
	logger.Println("error", message)
	http.Error(w, "Internal Server Error", code)
}
func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.Routes()
	return s
}
func (s *Server) Routes() {
	s.HandleFunc("/add-user", s.addUser).Methods(http.MethodPost)
	s.HandleFunc("/get-user", s.getUser).Methods(http.MethodGet)
	s.HandleFunc("/update-user", s.updateUser).Methods(http.MethodPut)
	s.HandleFunc("/add-vehicle", s.addVehicleToUser).Methods(http.MethodPost)
	s.HandleFunc("/get-user-vehicles", s.getUserVehicles).Methods(http.MethodGet)
}
func (s *Server) addUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("add user called")
	var temp User
	if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Println("temp", temp)
	id := 0
	err := db.QueryRow(`INSERT INTO cUsers(username,
		email,
		password,
		phone,
		location)
		VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING ID
		`, temp.Username, temp.Email, temp.Password, temp.Location, temp.Phone).Scan(&id)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.ID = id
	logger.Println("Created user with id", id)
	w.Header().Set("Content-Type", "Application/json")
	if err := json.NewEncoder(w).Encode(temp); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Created User with ID= %d successfully", id)
	logger.Println(response)
}
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("updateUser called")
	uid := r.URL.Query().Get("uid")
	var user User
	err := db.QueryRow(`SELECT password, location, phone FROM cUsers WHERE id=$1`, uid).Scan(&user.Password, &user.Location, &user.Phone)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var newInfo UpdateUserInfo
	if err := json.NewDecoder(r.Body).Decode(&newInfo); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, error := db.Exec(`UPDATE cUsers SET 
		password=$1,
		location=$2,
		phone=$3
	  WHERE id=$4
	`, findOptimalValue(user.Password, newInfo.Password),
		findOptimalValue(user.Location, newInfo.Location),
		findOptimalValue(user.Phone, newInfo.Phone), uid)
	if error != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "Application/json")
	if err := json.NewEncoder(w).Encode(newInfo); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("uid")
	var user User
	err := db.QueryRow(`SELECT * FROM cUsers WHERE id = $1`, id).Scan(&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.Location)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.Println("user", user)
	w.Header().Set("Content-Type", "Application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *Server) addVehicleToUser(w http.ResponseWriter, r *http.Request) {
	logger.Println("addVehicleToUser called")
	var vehicle Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehicleId := 0
	err := db.QueryRow(`INSERT INTO cVehicles (uid,
		make,
		model,
		year,
		color_exterior,
		color_interior,
		vin
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	)
	RETURNING ID
	`, vehicle.Uid, vehicle.Make, vehicle.Model, vehicle.Year, vehicle.ExteriorColor, vehicle.InteriorColor, vehicle.Vin).Scan(&vehicleId)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vehicle.ID = vehicleId
	logger.Println("Created vehicle with id", vehicleId)
	w.Header().Set("Content-Type", "Application/json")
	if err := json.NewEncoder(w).Encode(vehicle); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Created User with ID= %d successfully", vehicleId)
	logger.Println(response)
}
func (s *Server) getUserVehicles(w http.ResponseWriter, r *http.Request) {
	logger.Println("getUserVehicles called")
	var vehicles []Vehicle
	uid := r.URL.Query().Get("uid")
	rows, err := db.Query(`SELECT * FROM cVehicles WHERE uid = $1`, uid)
	if err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var v Vehicle
		rows.Scan(&v.ID, &v.Uid, &v.Make, &v.Model, &v.Year, &v.ExteriorColor, &v.InteriorColor, &v.Vin)
		vehicles = append(vehicles, v)
	}
	if err := json.NewEncoder(w).Encode(vehicles); err != nil {
		handleError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

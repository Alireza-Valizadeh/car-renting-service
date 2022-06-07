package api

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}
type UpdateUserInfo struct {
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Location string `json:"location"`
}

type Vehicle struct {
	ID            int    `json:"id"`
	Uid           int    `json:"uid"`
	Make          string `json:"make"`
	Model         string `json:"model"`
	Year          int    `json:"year"`
	InteriorColor string `json:"color_interior"`
	ExteriorColor string `json:"color_exterior"`
	Vin           string `json:"vin"`
}
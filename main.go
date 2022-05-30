package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alireza-valizadeh/car-renting-service/api"
)

func main() {
	s := api.NewServer()
	fmt.Printf("Starting server at port 127.0.0.1:4040\n")
	if err := http.ListenAndServe("127.0.0.1:4040", s); err != nil {
		log.Fatal(err)
	}
}

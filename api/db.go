package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "139142625"
	dbname   = "car-renting-service"
)

func NewDB() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgreSQL!")
	_,error := db.Exec(`CREATE TABLE IF NOT EXISTS cUsers (
		ID SERIAL PRIMARY KEY,
		username varchar(70) not null,
		email varchar(100) not null,
		password varchar(100) not null,
		phone varchar(15)  null,
		location varchar(30) null
	);`)
	if error != nil {
		panic(error)
	}
	_, carError := db.Exec(`CREATE TABLE IF NOT EXISTS cVehicles (
		ID SERIAL PRIMARY KEY,
		uid int,
		make varchar(70) not null,
		model varchar(100) not null,
		year varchar(100) not null,
		color_exterior varchar(30)  null,
		color_interior varchar(30) null,
		vin varchar(17) null,
		CONSTRAINT fk_user_vehicle
			FOREIGN KEY (uid)
				REFERENCES cUsers(ID)
	);`)
	if carError != nil {
		panic(carError)
	}
	return db
}

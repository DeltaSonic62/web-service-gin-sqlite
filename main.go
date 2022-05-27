package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// Represents data about a car
type car struct {
	ID    string `json:"id"`
	Year  int    `json:"year"`
	Make  string `json:"make"`
	Model string `json:"model"`
}

// cars slice to seed record car data
var cars = []car{}

// Database connection
var db = dbGetDB()

func main() {
	dbInit()

	router := gin.Default()

	router.GET("/cars", getCars)
	router.GET("/cars/:id", getCarById)
	router.GET("/cars/year/:year", getCarsByYear)
	router.GET("/cars/make/:make", getCarsByMake)
	router.GET("/cars/model/:model", getCarsByModel)
	router.POST("/cars", postCar)
	router.DELETE("/cars/:id", deleteCarById)

	router.Run("localhost:1997")

	defer db.Close()
}

// Responds with the list of cars as JSON
func getCars(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, cars)
}

// Responds with a car with a matching id as JSON
func getCarById(ctx *gin.Context) {
	id := ctx.Param("id")

	for _, c := range cars {
		if c.ID == id {
			ctx.IndentedJSON(http.StatusOK, c)
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found."})
}

// Responds with cars with matching years as JSON
func getCarsByYear(ctx *gin.Context) {
	year, err := strconv.Atoi(ctx.Param("year"))
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid year."})
		return
	}

	var carsByYear []car

	for _, c := range cars {
		if c.Year == year {
			carsByYear = append(carsByYear, c)
		}
	}

	if len(carsByYear) > 0 {
		ctx.IndentedJSON(http.StatusOK, carsByYear)
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "year not found."})
}

// Responds with cars with matching makes as JSON
func getCarsByMake(ctx *gin.Context) {
	make := ctx.Param("make")

	var carsByMake []car

	for _, c := range cars {
		if c.Make == make {
			carsByMake = append(carsByMake, c)
		}
	}

	if len(carsByMake) > 0 {
		ctx.IndentedJSON(http.StatusOK, carsByMake)
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "make not found."})
}

// Responds with cars with matching models as JSON
func getCarsByModel(ctx *gin.Context) {
	model := ctx.Param("model")

	var carsByModel []car

	for _, c := range cars {
		if c.Model == model {
			carsByModel = append(carsByModel, c)
		}
	}

	if len(carsByModel) > 0 {
		ctx.IndentedJSON(http.StatusOK, carsByModel)
		return
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "model not found."})
}

// Adds a car from JSON received in request body
func postCar(ctx *gin.Context) {
	var newCar car

	if err := ctx.BindJSON(&newCar); err != nil {
		return
	}

	cars = append(cars, newCar)
	dbAddCar(newCar)
	ctx.IndentedJSON(http.StatusCreated, newCar)
}

// Delete a car with a matching id
func deleteCarById(ctx *gin.Context) {
	id := ctx.Param("id")

	for i, c := range cars {
		if c.ID == id {
			cars = append(cars[:i], cars[i+1:]...)
			dbDeleteCar(c.ID)
			ctx.IndentedJSON(http.StatusOK, c)
			return
		}
	}

	ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "car not found."})
}

// Return database connection
func dbGetDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./cars.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Initalize the database
func dbInit() {
	carRows := dbGetRows()
	cars = carRows

	if len(carRows) == 0 {
		sqlStatement := `
		drop table if exists cars;
		create table cars (id text primary key, year integer not null, make text not null, model text not null);
		`
		_, err := db.Exec(sqlStatement)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Return number of rows in db
func dbGetRows() []car {
	var carRows = []car{}

	rows, err := db.Query("select * from cars;")
	if err != nil {
		return carRows
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		var year int
		var make string
		var model string

		err = rows.Scan(&id, &year, &make, &model)
		if err != nil {
			return carRows
		}

		carRows = append(carRows, car{ID: id, Year: year, Make: make, Model: model})
	}

	return carRows
}

func dbAddCar(newCar car) {
	sqlStatement := `
	insert into cars (id, year, make, model) values ($1, $2, $3, $4);
	`
	_, err := db.Exec(sqlStatement, newCar.ID, newCar.Year, newCar.Make, newCar.Model)
	if err != nil {
		log.Fatal(err)
	}
}

func dbDeleteCar(id string) {
	sqlStatement := `
	delete from cars where id = $1;
	`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatal(err)
	}
}

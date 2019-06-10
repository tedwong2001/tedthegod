package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

const (
	DbHost        = "db"
	DbUser        = "ted"
	DbPassword    = "ted"
	DbName        = "dev"
	DbCreateTable = `CREATE TABLE IF NOT EXISTS OrderList (
id serial PRIMARY KEY,
distance int NOT NULL,
status text NOT NULL) `
)

type Order struct {
	Id       int     `json:"id" binding:"required"`
	Distance float64 `json:"distance" binding:"required"`
	Status   string  `json:"status" binding:"required"`
}

type Locations struct {
	Origin      [2]string `json:"origin" binding:"required"`
	Destination [2]string `json:"destination" binding:"required"`
}

var db *sql.DB

func getOrder(id int) (Order, error) {
	const query = `SELECT id,distance,status FROM OrderList WHERE id = $1`
	row := db.QueryRow(query, id)
	var distance float64
	var status string
	err := row.Scan(&id, &distance, &status)
	if err != nil {
		return Order{}, err
	}
	result := Order{id, distance, status}
	return result, nil
}

func createOrder(locations Locations) (Order, error) {
	var id int
	distance := 2145
	status := "UNSIGNNED"
	const query = `INSERT INTO OrderList ("distance","status") VALUES ($1,$2) RETURNING id`
	err := db.QueryRow(query, distance, status).Scan(&id)
	if err != nil {
		return Order{}, err
	} else {
		var order, err = getOrder(id)
		if err != nil {
			return order, err
		} else {
			return order, err
		}
	}
}

func main() {
	var err error

	r := gin.Default()

	r.POST("/orders", func(context *gin.Context) {
		var location Locations
		if context.Bind(&location) == nil {
			result, err := createOrder(location)
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"Status": "Internal Error:" + err.Error()})
				return
			}
			context.JSON(http.StatusOK, result)
		}
	})

	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Query(DbCreateTable)
	if err != nil {
		log.Println("Failed to Run Table Creation", err.Error())
		return
	}

	log.Println("Running...")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}

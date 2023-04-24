package main

import (
	"go-mygram/database"
	"go-mygram/routers"
)

func main() {

	database.StartDB()
	r := routers.StartApp()
	r.Run(":8080")
}

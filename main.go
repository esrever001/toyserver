package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/esrever001/toyserver/db"
	"github.com/esrever001/toyserver/handlers"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	fmt.Printf("Initializing routing\n")
	router := httprouter.New()

	fmt.Printf("Initializing database\n")
	database := db.Database{
		Type:     "sqlite3",
		Filename: "test.db",
	}
	database.Init()

	handlers := handlers.CreateHandlers(&database)
	for _, handler := range handlers {
		fmt.Printf("Adding handler %s\n", handler.Path())
		router.GET(handler.Path(), handler.Handle)
	}

	log.Fatal(http.ListenAndServe(":4190", router))
}

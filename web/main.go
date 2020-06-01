package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/web/controller"
	"github.com/eldelto/solvent/web/persistence"
	serv "github.com/eldelto/solvent/web/service"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Controller defines the base methods any controller should implement
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

var repository = persistence.NewInMemoryRepository()
var service = serv.NewService(&repository)
var mainController = controller.NewMainController(&service)

func main() {
	port := 8080

	// TODO: Remove afterwards
	notebook, _ := solvent.NewNotebook()
	notebook.ID = uuid.Nil

	list, _ := notebook.AddList("My Server Side List")
	list.AddItem("Item0")
	list.AddItem("Item1")

	repository.Store(notebook)

	r := mux.NewRouter()
	mainController.RegisterRoutes(r)

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/").Handler(fs)

	http.Handle("/", r)

	log.Printf("Listening on localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

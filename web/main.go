package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/internal/conf"
	serv "github.com/eldelto/solvent/service"
	"github.com/eldelto/solvent/web/controller"
	"github.com/eldelto/solvent/web/persistence"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Controller defines the base methods any controller should implement
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

var simCp = conf.NewFileConfigProvider("conf/sim.properties")
var prodCp = conf.NewFileConfigProvider("conf/prod.properties")
var secretsCp = conf.NewFileConfigProvider("secrets/prod.properties")
var config = conf.NewChainConfigProvider([]conf.ConfigProvider{simCp, prodCp, secretsCp})

//var repository = persistence.NewInMemoryRepository()
var repository, postgresRepositoryErr = persistence.NewPostgresRepository(
	config.GetString("postgres.host"),
	config.GetString("postgres.port"),
	config.GetString("postgres.user"),
	config.GetString("postgres.password"),
)

var service = serv.NewService(repository)
var mainController = controller.NewMainController(&service)

func main() {
	// TODO: Where to handle re-connection?
	if postgresRepositoryErr != nil {
		panic(postgresRepositoryErr)
	}

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
	r.PathPrefix("/").Handler(handlers.CompressHandler(fs))

	http.Handle("/", r)

	log.Printf("Listening on localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

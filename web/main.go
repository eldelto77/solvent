package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eldelto/solvent/web/dto"
	"github.com/eldelto/solvent/web/persistence"
	serv "github.com/eldelto/solvent/web/service"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Controller struct {
	Handler http.Handler
}

func wireMainController() *Controller {
	r := mux.NewRouter()
	r.Handle("/health", baseMiddleWare(fetchHealth)).Methods("GET")
	r.Handle("/api/to-do-list", baseMiddleWare(fetchToDoLists)).Methods("GET")
	r.Handle("/api/to-do-list/{id}", baseMiddleWare(fetchToDoList)).Methods("GET")
	r.Handle("/api/to-do-list", baseMiddleWare(createToDoList)).Methods("POST")
	r.Handle("/api/to-do-list", baseMiddleWare(updateToDoList)).Methods("PUT")

	return &Controller{
		Handler: r,
	}
}

var repository = persistence.NewInMemoryRepository()
var service = serv.NewService(&repository)
var MainController = wireMainController()

func main() {
	port := 8080

	http.Handle("/", MainController.Handler)
	log.Printf("Listening on localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func fetchHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func fetchToDoLists(w http.ResponseWriter, r *http.Request) {
	toDoLists := service.FetchAll()

	dtos := make([]dto.ToDoListDto, len(toDoLists))
	for i, toDoList := range toDoLists {
		dtos[i] = dto.ToDoListToDto(&toDoList)
	}

	response := map[string][]dto.ToDoListDto{"toDoLists": dtos}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func fetchToDoList(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := service.Fetch(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	dto := dto.ToDoListToDto(list)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

type CreateRequest struct {
	Title string
}

func createToDoList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request CreateRequest
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := service.Create(request.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := dto.ToDoListToDto(list)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func updateToDoList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request dto.ToDoListDto
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newList := dto.ToDoListFromDto(&request)

	mergedList, err := service.Update(&newList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := dto.ToDoListToDto(mergedList)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func baseMiddleWare(nextFunc http.HandlerFunc) http.Handler {
	next := http.Handler(nextFunc)
	next = handlers.CombinedLoggingHandler(os.Stdout, next)
	next = handlers.ContentTypeHandler(next, "application/json")
	next = responseContentTypeHandler(next, "application/json")

	return next
}

func responseContentTypeHandler(next http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		next.ServeHTTP(w, r)
	})
}

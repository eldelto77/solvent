package controller

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/web/dto"
	"github.com/eldelto/solvent/web/service"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type MainController struct {
	service service.Service
}

func NewMainController(service service.Service) MainController {
	return MainController{
		service: service,
	}
}

func (c *MainController) RegisterRoutes(r *mux.Router) {
	r.Handle("/api/health", baseMiddleWare(c.fetchHealth)).Methods("GET")
	r.Handle("/api/to-do-list", baseMiddleWare(c.fetchToDoLists)).Methods("GET")
	r.Handle("/api/to-do-list/{id}", baseMiddleWare(c.fetchToDoList)).Methods("GET")
	r.Handle("/api/to-do-list", baseMiddleWare(c.createToDoList)).Methods("POST")
	r.Handle("/api/to-do-list", baseMiddleWare(c.updateToDoList)).Methods("PUT")
	r.Handle("/api/to-do-list/bulk", baseMiddleWare(c.pushToDoLists)).Methods("POST")
}

func (c *MainController) fetchHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func (c *MainController) fetchToDoLists(w http.ResponseWriter, r *http.Request) {
	toDoLists := c.service.FetchAll()

	dtos := make([]dto.ToDoListDto, len(toDoLists))
	for i, toDoList := range toDoLists {
		dtos[i] = dto.ToDoListToDto(&toDoList)
	}

	type response struct {
		ToDoLists []dto.ToDoListDto `json:"toDoLists"`
	}
	responseBody := response{ToDoLists: dtos}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
}

func (c *MainController) fetchToDoList(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := c.service.Fetch(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	dto := dto.ToDoListToDto(list)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func (c *MainController) createToDoList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// TODO: Benchmark if returning a closure over the type is faster than creating it
	type request struct {
		Title string
	}

	var requestBody request
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	list, err := c.service.Create(requestBody.Title)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := dto.ToDoListToDto(list)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func (c *MainController) updateToDoList(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var request dto.ToDoListDto
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newList := dto.ToDoListFromDto(&request)

	mergedList, err := c.service.Update(&newList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dto := dto.ToDoListToDto(mergedList)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func (c *MainController) pushToDoLists(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	type request struct {
		ToDoLists []dto.ToDoListDto
	}
	var requestBody request
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	toDoLists := make([]solvent.ToDoList, len(requestBody.ToDoLists))
	for i, toDoListDto := range requestBody.ToDoLists {
		toDoLists[i] = dto.ToDoListFromDto(&toDoListDto)
	}

	mergedList, err := c.service.BulkUpdate(toDoLists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dtos := make([]dto.ToDoListDto, len(mergedList))
	for i, toDoList := range toDoLists {
		dtos[i] = dto.ToDoListToDto(&toDoList)
	}

	type response struct {
		ToDoLists []dto.ToDoListDto
	}
	responseBody := response{ToDoLists: dtos}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseBody)
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

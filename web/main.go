package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/eldelto/solvent"
	"github.com/eldelto/solvent/web/dto"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := 8080
	r := mux.NewRouter()

	r.Handle("/health", baseMiddleWare(fetchHealth)).Methods("GET")
	r.Handle("/api/to-do-list", baseMiddleWare(fetchToDoList)).Methods("GET")

	http.Handle("/", r)
	log.Printf("Listening on localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func fetchHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func fetchToDoList(w http.ResponseWriter, r *http.Request) {
	list, _ := solvent.NewToDoList("test-list")
	dto := dto.ToDoListToDto(&list)

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

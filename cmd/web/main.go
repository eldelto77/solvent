package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := 8080
	r := mux.NewRouter()

	r.Handle("/health", baseMiddleWare(http.HandlerFunc(handleHealth))).Methods("GET")

	http.Handle("/", r)
	log.Printf("Listening on localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func baseMiddleWare(next http.Handler) http.Handler {
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

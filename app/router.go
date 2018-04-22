package app

import (
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sirupsen/logrus"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Use(loggingMiddleware)
	return router
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	logrus.WithFields(httpLogFields(r, 0, 404, 0)).Info()
}

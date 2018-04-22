package main

import (
	"github.com/sahilm/go-twitter/app"
	"github.com/gorilla/mux"
	"log"
	"github.com/pkg/errors"
)

func main() {
	a, err := app.New("twitter", mux.NewRouter())
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to create application"))
	}
	a.Run("8080")
}

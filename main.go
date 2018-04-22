package main

import (
	"github.com/pkg/errors"
	"github.com/sahilm/go-twitter/app"
	"github.com/sirupsen/logrus"
)

func main() {
	a, err := app.New("twitter", app.NewRouter())
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create application"))
	}
	a.Run("8080")
}

package main

import (
	"github.com/pkg/errors"
	"github.com/sahilm/go-twitter/app"
	"github.com/sahilm/go-twitter/log"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(log.UTCFormatter{Formatter: &logrus.TextFormatter{}})
	a, err := app.New("twitter", app.NewRouter())
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to create application"))
	}
	a.Run("8080")
}

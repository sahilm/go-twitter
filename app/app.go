package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func New(dbName string, router *mux.Router) (App, error) {
	app := App{}
	db, err := sqlx.Open("postgres", "dbname="+dbName)
	if err != nil {
		return app, err
	}
	app.DB = db
	app.Router = router
	return app, nil
}

func (app App) Run(port string) {
	addr := "0.0.0.0:" + port
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      app.Router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	logrus.Info("server running on " + addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	wait := time.Second * 15
	logrus.Info("received shutdown signal. Draining connections for a maximum of " + wait.String())
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "failed to shutdown server"))
	}
	logrus.Info("server going down")
}

package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
)

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
	Box    packr.Box
}

func New(dbName string, router *mux.Router) (App, error) {
	app := App{}
	db, err := sqlx.Open("postgres", fmt.Sprintf("sslmode=disable dbname=%s", dbName))
	if err != nil {
		return app, err
	}
	err = db.Ping()
	if err != nil {
		return app, err
	}
	app.DB = db
	app.Router = router
	app.Box = packr.NewBox("../db_migrations")
	return app, nil
}

func (a App) Run(port string) {
	a.migrateDB()
	addr := "0.0.0.0:" + port
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
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

func (a App) migrateDB() {
	migrations := &migrate.PackrMigrationSource{
		Box: a.Box,
	}
	n, err := migrate.Exec(a.DB.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		logrus.Fatal(errors.Wrap(err, "db migrations failed"))
	}
	logrus.Info(fmt.Sprintf("applied %d migrations", n))
}

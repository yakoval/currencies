package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"immo-currencies/project"
)

func main() {
	conf := project.NewConfig()
	log := logrus.New()

	db, err := project.Open(&conf.Database)
	project.Check(log, err, "starting DB")
	defer func() {
		err = db.Close()
		project.Check(log, err, "closing DB")
	}()

	restServer, err := project.StartWebServer(log, &conf.Web, db)
	project.Check(log, err, "starting webserver")
	defer func() {
		err = restServer.Close()
		project.Check(log, err, "closing webserver")
	}()

	waitForInterrupt(log)
}

func waitForInterrupt(log *logrus.Logger) {
	log.Info("waiting for interrupt")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Info("exiting by interrupt signal")
	os.Exit(1)
}

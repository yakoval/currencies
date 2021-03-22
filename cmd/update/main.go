package main

import (
	"github.com/sirupsen/logrus"

	"immo-currencies/project"
	"immo-currencies/project/update"
)

func main() {
	conf := project.NewConfig()
	log := logrus.New()

	log.WithFields(logrus.Fields{
		"stage": "updating currencies",
	}).Info("start...")

	db, err := project.Open(&conf.Database)
	project.Check(log, err, "updating currencies")
	defer func() {
		err = db.Close()
		project.Check(log, err, "updating currencies")
	}()

	updater, err := update.NewUpdater(db, &conf.Updater)
	project.Check(log, err, "updating currencies")

	err = updater.Work()
	project.Check(log, err, "updating currencies")

	log.WithFields(logrus.Fields{
		"stage":    "updating currencies",
		"read":     updater.ReadCounter,
		"inserted": updater.InsertCounter,
		"updated":  updater.UpdateCounter,
		"deleted":  updater.DeleteCounter,
	}).Info("finished!")
}

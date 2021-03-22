package project

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Check(log *logrus.Logger, err error, stage string) {
	if err != nil {
		log.WithFields(logrus.Fields{
			"stage": stage,
		}).Error(err)
		os.Exit(1)
	}
}

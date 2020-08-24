package main

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/acoustic/csv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {


	config,err := csv.InitConfig("config.yaml")
	log.Info(config)
	log.Error(err)
	reeceContent,err := config.Get("Moodboard Reece Product")
	log.Info(reeceContent)
	log.Error(err)
}
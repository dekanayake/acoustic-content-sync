package main

import (
	"github.com/dekanayake/acoustic-content-sync/pkg/command"
	"github.com/dekanayake/acoustic-content-sync/pkg/logrus"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	acousticSyncLog, err := os.OpenFile("acoustic_sync.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Panic("Error in writing to log file", err)
	}
	log.SetFormatter(Formatter)
	mw := io.MultiWriter(os.Stdout, acousticSyncLog)
	log.SetOutput(mw)
	log.SetLevel(log.InfoLevel)
}

func execute() {
	log.Info("--------------Running Synky CLI----------------")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cmdProcessor := command.NewCommandActionProcessor()
	errorHandling := logrus.PkgErrorEntry{Entry: log.WithField("", "")}
	err = cmdProcessor.Execute()
	if err != nil {
		errorHandling.WithError(err).Panic(err)
	}
}

func main() {
	execute()
}

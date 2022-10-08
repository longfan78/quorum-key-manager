package main

import (
	"github.com/longfan78/quorum-key-manager/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	command := cmd.NewCommand()
	if err := command.Execute(); err != nil {
		log.WithError(err).Fatalf("execution failed")
	}
}

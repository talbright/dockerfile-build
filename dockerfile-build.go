package main

import (
	log "github.com/Sirupsen/logrus"
	commands "github.com/talbright/dockerfile-build/commands"
)

func main() {
	log.SetLevel(log.DebugLevel)
	commands.RootCmd.Execute()
}

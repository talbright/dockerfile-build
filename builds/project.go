package build

import (
	log "github.com/Sirupsen/logrus"
)

type Project struct {
	Config *Build
	Flags *BuildFlags
}

func (b *Project) Build() error {
	log.WithFields(log.Fields{
		"build" : "Project",
	}).Debug("build")
	delegate := Base{Config:b.Config, Flags: b.Flags}
	return delegate.Build()
}


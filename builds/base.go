package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	log "github.com/Sirupsen/logrus"
	shell "github.com/talbright/dockerfile-build/shell"
)

type Base struct {
	Config *Build
	Flags *BuildFlags
}

func (b *Base) Build() error {
	log.WithFields(log.Fields{
		"phase" : "start",
	}).Info("build")

	versions := strings.Split(b.Flags.WithVersion,",")
	if len(versions) == 0 {
		versions = append(versions,"latest")
	}

	//phase: pull
	//get versions that might already exist
	for _,version := range versions {
		tag := fmt.Sprintf("registry.mycompany.com/%s:%s",b.Config.Tag,version)
		cmd,_ := shell.Pull(tag)
		log.WithFields(log.Fields{
			"phase" : "pull",
			"tag" : tag,
		}).Info("pulling image")
		if cmd.ProcessState.Success() {
			log.WithFields(log.Fields{
				"phase" : "pull",
				"tag" : tag,
			}).Warn("image already exists")
		}
	}

	//phase: pre-pull
	//pull pre-requisites (images in Dockerfile FROM tag,etc)
	for _, image := range b.Config.Prepull {
		log.WithFields(log.Fields{
			"phase" : "pre-pull",
			"tag" : image,
		}).Info("pre-pulling image")
		cmd,err := shell.Pull(image)
		if !cmd.ProcessState.Success() {
			log.WithFields(log.Fields{
				"phase" : "pre-pull",
				"tag" : image,
				"error" : err,
			}).Warn("failed to pull image")
		}
	}

	//phase: build
	//build it finally
	buildDir := filepath.Join(filepath.Dir(b.Flags.Config),b.Config.Directory)
	log.WithFields(log.Fields{
		"key"   : "buildDir",
		"value" : buildDir,
	}).Debug("build")

	if exists,err := shell.IsDirectory(buildDir); err!=nil || !exists {
		log.WithFields(log.Fields{
			"phase" : "build",
			"error" : err,
		}).Error("invalid path to Dockerfile")
		os.Exit(1)
	}
	
	if cmd,err := shell.Build(b.Config.Tag,buildDir); err!=nil || !cmd.ProcessState.Success() {
		log.WithFields(log.Fields{
			"phase" : "build",
			"tag" : b.Config.Tag,
			"error" : err,
		}).Error("failed to build image")
		os.Exit(1)
	}

	//phase: tag
	for _,version := range versions {
		versionTag := fmt.Sprintf("registry.mycompany.com/%s:%s",b.Config.Tag,version)
		log.WithFields(log.Fields{
			"phase" : "tag",
			"tag" : versionTag,
		}).Info("build")
		shell.Tag(b.Config.Tag,versionTag)
	}

	//phase: push
	log.WithFields(log.Fields{
		"phase" : "push",
		"enabled" : b.Flags.EnablePush,
	}).Info("build")
	
	if b.Flags.EnablePush {
		for _,version := range versions {
			versionTag := fmt.Sprintf("registry.mycompany.com/%s:%s",b.Config.Tag,version)
			log.WithFields(log.Fields{
				"phase" : "tag",
				"tag" : versionTag,
			}).Info("build")
			shell.Push(versionTag)
		}
	}

	return nil
}


package commands

import (
	"strings"
	"io/ioutil"
	"encoding/json"
	"os"
	cobra "github.com/spf13/cobra"
	log "github.com/Sirupsen/logrus"
	build "github.com/talbright/dockerfile-build/builds"
)

var BuildFlags build.BuildFlags

var RootCmd = &cobra.Command {
	Use: "dockerfile-build",
	Run: func(cmd *cobra.Command, args []string) {
		Initialize(cmd, args)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&BuildFlags.Config, "config", "desk-build.json", "build configuration")
	RootCmd.PersistentFlags().StringVar(&BuildFlags.WithVersion, "with-version","","tag build(s) with provided version")
	RootCmd.PersistentFlags().StringVar(&BuildFlags.WithTag, "with-tag","","build only this tag (for configs with multiple build definitions)")
	RootCmd.PersistentFlags().BoolVar(&BuildFlags.EnableRelease, "enable-release",false,"release build")
	RootCmd.PersistentFlags().BoolVar(&BuildFlags.EnablePush, "enable-push",false,"push build")
	RootCmd.PersistentFlags().BoolVar(&BuildFlags.EnableForce, "enable-force",false,"force build")
}

func Initialize(cmd *cobra.Command, args []string) {
	log.WithFields(log.Fields{
		"command" : "dockerfile-build",
		"args" : strings.Join(args, " "),
		"flags" : BuildFlags,
	}).Debug("command")

	if BuildFlags.Config == "" {
		log.Error("no build config provided")
		os.Exit(1)
	}

	configData, err := ioutil.ReadFile(BuildFlags.Config)
	if err != nil {
		panic(err)
	}
	var buildConfig build.BuildConfig
	if err := json.Unmarshal(configData, &buildConfig); err != nil {
		panic(err)
	}

	if BuildFlags.WithTag != "" {
		b := buildConfig.FindByTag(BuildFlags.WithTag)
		if b == nil {
			log.Errorf("build tag not found in config %s",BuildFlags.WithTag)
			os.Exit(1)
		}
		InvokeBuild(b)
	} else {
		for _,b := range buildConfig.Builds {
			log.WithFields(log.Fields{
				"config" : b,
			}).Debug("build")
			InvokeBuild(&b)
		}
	}
}

func InvokeBuild(b *build.Build) {
	switch b.Type {
	case "base":
		buildType := build.Base{Config: b, Flags: &BuildFlags }
		buildType.Build()
	case "project":
		buildType := build.Project{Config: b, Flags: &BuildFlags }
		buildType.Build()
	default:
		log.Errorf("undefined build type: %s",b.Type)
		os.Exit(1)
	}
}


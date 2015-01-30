package build

import (
	. "github.com/talbright/dockerfile-build/utils"
)

type BuildInterface interface {
	Build() error
}

type BuildConfig struct {
	Builds []Build           `json:"build,omitempty"`
}

func (s BuildConfig) String() string {
	return Stringify(s)
}

func (config *BuildConfig) FindByTag(tag string) (build *Build) {
	build = nil
	for _,b := range config.Builds {
		if b.Tag == tag {
			build = &b
		}
	}
	return
}

type Build struct {
	Directory string         `json:"directory,omitempty"`
	Tag string               `json:"tag,omitempty"`
	Type string              `json:"type,omitempty"`
	Prepull []string         `json:"prepull,omitempty"`
	Touch []string           `json:"touch,omitempty"`
}

func (s Build) String() string {
	return Stringify(s)
}


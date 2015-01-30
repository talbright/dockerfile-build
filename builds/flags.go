package build

import (
	. "github.com/talbright/dockerfile-build/utils"
)

type BuildFlags struct {
	Config string
	WithVersion string
	WithTag string
	EnablePush bool
	EnableForce bool
	EnableRelease bool
}

func (s BuildFlags) String() string {
	return Stringify(s)
}


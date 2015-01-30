package shell

import (
	"os"
	"os/exec"
	"strings"
	"bytes"
	log "github.com/Sirupsen/logrus"
)

func LogStderr(cmd *exec.Cmd,err error) {
	if err != nil {
		out := cmd.Stderr.(*bytes.Buffer)
		log.Error(strings.Join(cmd.Args," "),":\n",out.String())
		log.Error(err)
	}
}

func Pull(tag string) (*exec.Cmd,error) {
	shellString := []string{tag}
	return Docker("pull",shellString...)
}

func Push(tag string) (*exec.Cmd,error) {
	shellString := []string{tag}
	return Docker("push",shellString...)
}

func Build(tag string,dir string) (*exec.Cmd,error) {
	shellString := []string{"-t",tag,dir}
	return Docker("build",shellString...)
}

func Tag(tag1 string, tag2 string) (*exec.Cmd,error) {
	shellString := []string{tag1,tag2}
	return Docker("tag",shellString...)
}

func Images() (*exec.Cmd,error) {
	return Docker("images")
}

func Docker(subcmd string,arg ...string) (*exec.Cmd,error) {
	withSubcmd := append([]string{subcmd},arg...)
	log.WithFields(log.Fields{
		"command" : "docker",
		"args" : strings.Join(arg, " "),
	}).Debug("shell")
	cmd := exec.Command("docker",withSubcmd...)
	cmd.Stdout = NewOutputStreamer()
	cmd.Stderr = NewOutputStreamer()
	err := cmd.Run()
	log.WithFields(log.Fields{
		"exit" : err,
	}).Debug("shell")
	return cmd,err
}

func Sh(shell string) (*exec.Cmd,error) {
	log.WithFields(log.Fields{
		"command" : "sh",
		"args" : shell,
	}).Debug("shell")
	shellString := []string{"-c",shell}
	cmd := exec.Command("sh",shellString...)
	cmd.Stdout = NewOutputStreamer()
	cmd.Stderr = NewOutputStreamer()
	err := cmd.Run()
	log.WithFields(log.Fields{
		"exit" : err,
	}).Debug("shell")
	return cmd,err
}

func IsDirectory(path string) (bool, error) {
    fileInfo, err := os.Stat(path)
    return fileInfo.IsDir(), err
}

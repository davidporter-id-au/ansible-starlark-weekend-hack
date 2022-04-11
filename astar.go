package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/davidporter-id-au/ansible-starlark-weekend-hack/larker"
)

const timeout = 20 * time.Second

var starlarkExtension = regexp.MustCompile(".star$")

// converts starlark files in roles and group_vars
// to yaml before execution
func convertConfig() error {
	err := convertStarlarkFile("group_vars/all.star")
	if err != nil {
		return err
	}
	return filepath.WalkDir("roles", func(path string, d os.DirEntry, err error) error {
		if starlarkExtension.MatchString(path) {
			err := convertStarlarkFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func convertStarlarkFile(file string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	l := larker.New()

	res, err := l.Main(ctx, file)
	if err != nil {
		return err
	}
	log.Println("writing ", starlarkExtension.ReplaceAllString(file, ".yaml"))
	return ioutil.WriteFile(starlarkExtension.ReplaceAllString(file, ".yaml"), []byte(res.YAMLConfig), 0600)
}

func execAnsible(args []string) error {
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	err := convertConfig()
	if err != nil {
		log.Fatal("couldn't convert config: ", err)
	}
	err = execAnsible(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}

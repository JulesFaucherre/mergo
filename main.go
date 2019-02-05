package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/run"
)

func localRemote(remote string) (string, error) {
	return run.Run([][]string{
		{"git", "remote", "-v"},
		{"grep", remote},
		{"awk", "{print $2}"},
		{"head", "-n", "1"},
	})
}

func main() {
	opts := &models.Opts{}
	_, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(opts.Host) == 0 || len(opts.Owner) == 0 || len(opts.Repo) == 0 {
		remoteString, err := localRemote(opts.Remote)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(opts.Host) == 0 {
			opts.Host = hosts.GetHostNameFromRemoteString(remoteString)
		}
	}
}

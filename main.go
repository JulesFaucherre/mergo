package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

func main() {
	opts := &models.Opts{}

	_, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		return
	}
	var host models.Host
	var remoteString string

	if !tools.IsEmpty(opts.Delete) {
		if err = tools.DeleteHostConfig(opts.Delete); err != nil {
			fmt.Println(err)
		}
		return
	}

	if tools.IsEmpty(opts.Host) {
		remoteString, err = git.
			LocalRepository().
			Remote(opts.Remote).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		opts.Host = hosts.GetHostNameFromRemoteString(remoteString)
	}

	if host, err = hosts.GetHost(opts.Host); err != nil {
		fmt.Println(err)
		return
	}

	if tools.IsEmpty(opts.Owner) || tools.IsEmpty(opts.Repo) {
		if tools.IsEmpty(remoteString) {
			remoteString, err = git.
				LocalRepository().
				Remote(opts.Remote).
				Do(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		opts.Owner, opts.Repo = host.GetOwnerAndRepo(remoteString)
	}

	if tools.IsEmpty(opts.Head) {
		if opts.Head, err = git.
			LocalRepository().
			Branch().
			Do(context.Background()); err != nil {
			fmt.Println(err)
			return
		}
		opts.Head = strings.Trim(opts.Head, "\n")
	}

	if err = host.SubmitPr(opts); err != nil {
		fmt.Println(err)
		return
	}
}

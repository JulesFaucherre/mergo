package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

var (
	httpsR = regexp.MustCompile(`https://(\w+\.\w+)/([\w-]+)/([\w-]+).git`)
	sshR   = regexp.MustCompile(`git@(\w+\.\w+):([\w-]+)/([\w-]+).git`)
)

func getRemoteInformations(opts *models.Opts) error {
	remoteString, err := git.
		LocalRepository().
		Remote(opts.Remote).
		Do(context.Background())
	if err != nil {
		return err
	}

	matches := httpsR.FindStringSubmatch(remoteString)
	if matches == nil {
		matches = sshR.FindStringSubmatch(remoteString)
		if matches == nil {
			return fmt.Errorf("Unable to extract informations from remote string %s", remoteString)
		}
	}

	if tools.IsEmpty(opts.Host) {
		opts.Host = matches[1]
	}
	if tools.IsEmpty(opts.Owner) {
		opts.Owner = matches[2]
	}
	if tools.IsEmpty(opts.Repo) {
		opts.Repo = matches[3]
	}

	return nil
}

func main() {
	opts := &models.Opts{}

	_, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		return
	}
	var host models.Host

	if !tools.IsEmpty(opts.Delete) {
		if err = tools.DeleteHostConfig(opts.Delete); err != nil {
			fmt.Println(err)
		}
		return
	}

	if tools.IsEmpty(opts.Host) || tools.IsEmpty(opts.Owner) || tools.IsEmpty(opts.Repo) {
		getRemoteInformations(opts)
	}

	if host, err = hosts.GetHost(opts.Host); err != nil {
		fmt.Println(err)
		return
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

	commits, err := git.LocalRepository().GetDifferenceCommits(opts.Head, opts.Base)
	if err != nil {
		fmt.Println(err)
		return
	}
	opts.Commits = commits

	if err = host.SubmitPr(opts); err != nil {
		fmt.Println(err)
		return
	}
}

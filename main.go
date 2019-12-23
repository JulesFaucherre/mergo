package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/atotto/clipboard"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

var (
	httpsR = regexp.MustCompile(`https://(\w+\.\w+)/([\w-]+)/([\w-]+).git`)
	sshR   = regexp.MustCompile(`git@(\w+\.\w+):([\w-]+)/([\w-]+).git`)
)

func main() {
	opts, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	opts.Repo = git.LocalRepository()

	if err = fillDefaultValues(opts); err != nil {
		fmt.Println(err)
		return
	}

	var host models.Host

	if !tools.IsEmpty(opts.Delete) {
		if err = tools.DeleteHostConfig(opts.Delete); err != nil {
			fmt.Println(err)
		}
		return
	}

	if host, err = hosts.GetHost(opts.Host); err != nil {
		fmt.Println(err)
		return
	}

	if tools.Verbose {
		if ok, err := tools.AskYesNo("Do you still want to submit the pr ?"); err != nil {
			fmt.Println(err)
			return
		} else if !ok {
			return
		}
	}

	var mrInfo *models.MRInfo
	if mrInfo, err = host.SubmitPr(opts); err != nil {
		fmt.Println(err)
		return
	}
	if opts.Clipboard {
		clipboard.WriteAll(mrInfo.URL)
	}
	fmt.Printf("Your request is available at the following URL:\n%s", mrInfo.URL)
}

func fillDefaultValues(opts *models.Opts) error {
	var err error

	if tools.IsEmpty(opts.Head) {
		if opts.Head, err = opts.
			Repo.
			Branch().
			Do(context.Background()); err != nil {
			return err
		}
		opts.Head = strings.Trim(opts.Head, "\n")
	}

	if !tools.IsEmpty(opts.Host) && !tools.IsEmpty(opts.Owner) && !tools.IsEmpty(opts.Repository) {
		return nil
	}

	remoteString, err := opts.
		Repo.
		Remote(opts.Remote).
		Do(context.Background())
	if err != nil {
		return err
	}

	var matches []string
	if m := httpsR.FindStringSubmatch(remoteString); m != nil {
		matches = m
	} else if m := sshR.FindStringSubmatch(remoteString); m != nil {
		matches = m
	} else {
		return fmt.Errorf("Unable to extract informations from remote string %s", remoteString)
	}

	if tools.IsEmpty(opts.Host) {
		opts.Host = matches[1]
	}
	if tools.IsEmpty(opts.Owner) {
		opts.Owner = matches[2]
	}
	if tools.IsEmpty(opts.Repository) {
		opts.Repository = matches[3]
	}

	return nil
}

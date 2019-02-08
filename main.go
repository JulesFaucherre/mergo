package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

func localRemote(remote string) (string, error) {
	remoteString, err := tools.Run([][]string{
		{"git", "remote", "-v"},
		{"grep", remote},
		{"awk", "{print $2}"},
		{"head", "-n", "1"},
	})
	if err != nil {
		return "", err
	}
	if isEmpty(remoteString) {
		return "", fmt.Errorf("Unable to get remote informations. Aborting")
	}
	return remoteString, nil
}

func main() {
	opts := &models.Opts{}
	_, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	var host hosts.Host
	var remoteString string

	if isEmpty(opts.Host) {
		remoteString, err = localRemote(opts.Remote)
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

	if isEmpty(opts.Owner) || isEmpty(opts.Repo) {
		if isEmpty(remoteString) {
			remoteString, err = localRemote(opts.Remote)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		opts.Owner, opts.Repo = host.GetOwnerAndRepo(remoteString)
	}

	if err = handleMissingPrInformations(opts); err != nil {
		fmt.Println(err)
		return
	}

	if err = host.SubmitPr(opts); err != nil {
		fmt.Println(err)
		return
	}
}

func handleMissingPrInformations(opts *models.Opts) error {
	var err error
	stdin := bufio.NewReader(os.Stdin)

	if isEmpty(opts.Head) {
		if opts.Head, err = tools.Run([][]string{
			{"git", "branch"},
			{"grep", "*"},
			{"awk", "{print $2}"},
		}); err != nil {
			return err
		}
		opts.Head = strings.Trim(opts.Head, "\n")
	}

	if isEmpty(opts.Title) {
		fmt.Println("Enter the pull request's title:")
		if opts.Title, err = stdin.ReadString('\n'); err != nil {
			return err
		}
	}

	if isEmpty(opts.Body) {
		fmt.Println("Enter the pull request's body:")
		if opts.Body, _ = stdin.ReadString('\n'); err != nil {
			return nil
		}
	}

	return nil
}

func isEmpty(s string) bool {
	return len(s) == 0
}

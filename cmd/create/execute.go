package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

func getRemoteInformations(opts *models.CreateOptions) error {
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

func run(_ *cobra.Command, _ []string) {
	var err error

	if tools.Verbose {
		fmt.Printf("options = %+v\n", options)
	}
	var host models.Host

	if tools.IsEmpty(options.Host) || tools.IsEmpty(options.Owner) || tools.IsEmpty(options.Repo) {
		getRemoteInformations(options)
	}

	if host, err = hosts.GetHost(options.Host); err != nil {
		fmt.Println(err)
		return
	}

	if tools.IsEmpty(options.Head) {
		if options.Head, err = git.
			LocalRepository().
			Branch().
			Do(context.Background()); err != nil {
			fmt.Println(err)
			return
		}
		options.Head = strings.Trim(options.Head, "\n")
	}

	if err != nil {
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

	if tools.Verbose {
		fmt.Printf("options = %+v\n", options)
		fmt.Println("Press enter to continue")
		fmt.Scanf("%s")
	}

	var mrInfo *models.MRInfo
	if mrInfo, err = host.SubmitPr(options, git.LocalRepository()); err != nil {
		fmt.Println(err)
		return
	}
	if options.Clipboard {
		clipboard.WriteAll(mrInfo.URL)
	}
	fmt.Printf("Your request is available at the following URL:\n%s", mrInfo.URL)
}

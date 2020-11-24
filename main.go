package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/credentials"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/logger"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

func main() {
	opts := &models.Opts{}

	repo, err := git.FromPwd()
	tools.CheckError(err)

	err = repo.LoadConfig(&struct {
		Mergo *models.Opts `ini:"mergo"`
	}{Mergo: opts})
	tools.CheckError(err)

	_, err = flags.ParseArgs(opts, os.Args)
	if err != nil {
		return
	}
	models.SetDefaultOptions(opts)

	logger.Verbosity = len(opts.Verbose)
	logger.Info("verbosity = %d\n", logger.Verbosity)
	logger.Debug("opts = %+v\n", opts)

	err = setDefaultMessage(opts, repo)
	tools.CheckError(err)

	if opts.Delete != "" {
		credentials.DeleteHostConfig(opts.Delete)
		return
	}

	err = setHeadBranch(opts, repo)
	tools.CheckError(err)

	err = checkProblems(opts, repo)
	tools.CheckError(err)

	err = FormalizeMessage(opts, repo)
	tools.CheckError(err)

	mrParams := &models.MRParams{
		Head:    opts.Head,
		Base:    opts.Base,
		Message: opts.Message,
	}

	rmts, err := getRemotes(opts, repo)
	tools.CheckError(err)

	mrUrls := SubmitPr(mrParams, rmts)
	urls := strings.Join(mrUrls, "\n")

	if opts.Clipboard {
		clipboard.WriteAll(urls)
	}
}

func setHeadBranch(opts *models.Opts, repo *git.Repo) error {
	var err error
	if len(opts.Head) != 0 {
		return nil
	}
	opts.Head, err = repo.GetActualBranch()
	if err == git.ErrNotABranch {
		return fmt.Errorf("we can not take your branch as a source branch because you are not on a branch")
	}
	if err != nil {
		return err
	}
	return nil
}

// checkProblems checks for different problems that might hold you back from
// making your pull request
func checkProblems(opts *models.Opts, repo *git.Repo) error {
	if opts.Force {
		return nil
	}

	problems := []string{}

	// check if the repo has unstaged changes
	if repo.HasChanges() {
		problems = append(problems, "you have unstaged changes")
	}

	// check if head branch is up to date
	isHdUp, err := repo.IsBranchUpToDate(opts.Head)
	if err != nil {
		return err
	}
	if !isHdUp {
		problems = append(problems, "your head branch is not up to date")
	}

	// check if base branch is up to date
	isBsUp, err := repo.IsBranchUpToDate(opts.Base)
	if err != nil {
		return err
	}
	if !isBsUp {
		problems = append(problems, "your base branch is not up to date")
	}

	if len(problems) == 0 {
		return nil
	}

	// ask you if you still wish to continue
	problems = append([]string{"you have the following problems:"}, problems...)
	msg := strings.Join(problems, "\n\t- ") + "\ndo you wish to ignore them and continue ?"
	agreed, err := tools.AskYesNo(msg)
	if err != nil {
		return err
	}
	if !agreed {
		return errors.New("cancelled due to user input")
	}
	return nil
}

func getRemotes(opts *models.Opts, repo *git.Repo) (remotes []string, err error) {
	if len(opts.RemoteURLs) != 0 {
		remotes = opts.RemoteURLs
		return
	}
	remotes, err = repo.GetRemoteURLs(opts.Remote)
	return
}

func setDefaultMessage(opts *models.Opts, repo *git.Repo) error {
	head, err := repo.GetBranchCommit(opts.Head)
	if err != nil {
		return err
	}

	base, err := repo.GetBranchCommit(opts.Base)
	if err != nil {
		return err
	}

	// If there was only one commit on branch head from base, we set this commit
	// message as the default's message of the pull request's
	if git.IsDirectChild(base, head) {
		opts.Message = head.Message
	}

	return nil
}

package git

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
	Comment string
}

// This function parses commit formatted from the command
// `git log --date=iso`
func parseOneCommit(content []string) (*Commit, error) {
	commit := Commit{}
	var sp []string
	var err error

	if len(content) < 5 {
		return nil, fmt.Errorf("Commit is not long enough")
	}

	commitLine := content[0]
	sp = strings.Split(commitLine, " ")
	if len(sp) < 2 {
		return nil, fmt.Errorf("Commit line do not have hash")
	}
	commit.Hash = sp[1]

	authorLine := content[1]
	sp = strings.Split(authorLine, " ")
	if len(sp) < 2 {
		return nil, fmt.Errorf("Author line do not have author")
	}
	if sp[0] != "Author:" {
		return nil, fmt.Errorf("Author line has wrong format")
	}
	commit.Author = strings.Join(sp[1:], " ")

	dateLine := content[2]
	sp = strings.Split(dateLine, " ")
	if len(sp) < 2 {
		return nil, fmt.Errorf("Date line do not have date")
	}
	if sp[0] != "Date:" {
		return nil, fmt.Errorf("Author line has wrong format")
	}
	date := strings.Join(sp[3:], " ")
	if commit.Date, err = time.Parse("2006-01-02 15:04:05 -0700", date); err != nil {
		return nil, fmt.Errorf("While parsing date : \"%s\"\nError : %+v", date, err)
	}

	for i, l := range content[4:] {
		content[4+i] = strings.TrimSpace(l)
	}
	commit.Message = content[4]
	commit.Comment = strings.TrimSpace(strings.Join(content[5:], "\n"))

	return &commit, nil
}

func parseCommitList(content string) ([]*Commit, error) {
	commits := []*Commit{}
	sp := strings.Split(content, "\n")
	start := 0

	addCommit := func(start, end int) error {
		commitContent := sp[start:end]
		commit, err := parseOneCommit(commitContent)
		if err != nil {
			return fmt.Errorf("While parsing commit got :\nCommit :\n%s\nError :\n%+v", strings.Join(commitContent, "\n"), err)
		}
		commits = append(commits, commit)
		return nil
	}

	for i, line := range sp {
		if i != start && strings.HasPrefix(line, "commit") {
			if err := addCommit(start, i); err != nil {
				return nil, err
			}
			start = i
		}
	}
	addCommit(start, len(sp))
	return commits, nil
}

func (me *Repo) GetDifferenceCommits(head, base string) ([]*Commit, error) {
	return me.GetDifferenceCommitsWithContext(context.Background(), head, base)
}

func (me *Repo) GetDifferenceCommitsWithContext(ctx context.Context, head, base string) ([]*Commit, error) {
	cmd := GitCmd{
		repo: me,
		cmd: [][]string{
			{"git", "log", "--date=iso", base + ".." + head},
		},
	}
	r, err := cmd.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("While running command : %s\nError : \n%+v", strings.Join(cmd.cmd[0], " "), err)
	}
	return parseCommitList(r)
}

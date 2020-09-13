package hostTools

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
)

var (
	baseContent = `# Enter the content of your merge request
# Please enter your request's title and body. Lines starting with '#' will be
# ignored, and an empty content aborts the request.

# Title:

{{.Title}}

# Body:

{{.Body}}

`
)

type UserInfo struct {
	Title string
	Body  string
}

func defaultValues(commits []*git.Commit) (string, string) {
	if len(commits) != 1 {
		return "Enter merge request title", "Enter merge request body"
	}
	return commits[0].Message, commits[0].Comment
}

const (
	none = 0 + iota
	title
	body
)

func parseContent(content []byte) (*UserInfo, error) {
	str := strings.TrimLeft(string(content), " \n\t\r")
	if tools.IsEmpty(str) {
		return nil, fmt.Errorf("Aborted request")
	}
	lines := []string{}

	for _, l := range strings.Split(str, "\n") {
		if !strings.HasPrefix(l, "#") && !tools.IsEmpty(l) {
			lines = append(lines, l)
		}
	}
	if tools.Verbose {
		fmt.Printf("user input lines = %+v\n", strings.Join(lines, "\n"))
	}
	return &UserInfo{
		Title: lines[0],
		Body:  strings.Join(lines[1:], "\n"),
	}, nil
}

func DefaultGetUserInfo(opts *models.Opts) (*UserInfo, error) {
	userInfo := &UserInfo{}
	commits, err := opts.
		Repo.
		GetDifferenceCommits(opts.Head, opts.Base)
	if err != nil {
		return nil, err
	}
	userInfo.Title, userInfo.Body = defaultValues(commits)
	templ := template.Must(template.New("User info").Parse(baseContent))

	var tpl bytes.Buffer
	err = templ.Execute(&tpl, userInfo)
	content, err := git.EditText([]byte(tpl.String()))
	if err != nil {
		return nil, fmt.Errorf("While getting your input got error :\n%+v", err)
	}
	userInfo, err = parseContent(content)
	if err != nil {
		return nil, fmt.Errorf("While parsing your input got error :\n%+v", err)
	}
	if tools.Verbose {
		fmt.Printf("userInfo = %+v\n", userInfo)
	}
	return userInfo, nil
}

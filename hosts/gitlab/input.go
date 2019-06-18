package gitlab

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
	baseContent = `#  Enter the content of your merge request
#  Every line starting with one only '#' as well as empty lines will be considered as a comment and not treated
#  Do not change lines starting with '##' as they are used for parsing
#  Please note that new lines in gitlab's descriptions are <br />

## Title
#  Note that only the first line will be taken since merge request titles are monoline

{{.Title}}

## Body

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

func parseContent(content []byte) (string, string) {
	parsing := none
	titleC := ""
	bodyC := []string{}
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)

		// Swap to title parsing
		if line == "## Title" && tools.IsEmpty(titleC) {
			parsing = title
			continue
		}
		// Swap to body parsing
		if line == "## Body" && len(bodyC) == 0 {
			parsing = body
			continue
		}
		// Drop useless lines
		if (line[0] == '#' && line[1] != '#') || tools.IsEmpty(line) {
			continue
		}

		switch parsing {
		case none:
			break
		case title:
			if tools.IsEmpty(titleC) {
				titleC = line
			}
			break
		case body:
			bodyC = append(bodyC, line)
			break
		}
	}
	return titleC, strings.Join(bodyC, "\n")
}

func getTitleAndDescription(opts *models.Opts) (*UserInfo, error) {
	userInfo := &UserInfo{}
	userInfo.Title, userInfo.Body = defaultValues(opts.Commits)
	templ := template.Must(template.New("User info").Parse(baseContent))

	var tpl bytes.Buffer
	err := templ.Execute(&tpl, userInfo)
	content, err := git.EditText([]byte(tpl.String()))
	if err != nil {
		return nil, fmt.Errorf("While getting your input got error :\n%+v", err)
	}
	userInfo.Title, userInfo.Body = parseContent(content)
	if err != nil {
		return nil, fmt.Errorf("While parsing your input got error :\n%+v", err)
	}
	return userInfo, nil
}

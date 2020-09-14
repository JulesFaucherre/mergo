package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/models"
)

const (
	none = 0 + iota
	title
	body
)

var (
	baseContent = `#  Enter your pull request's message
#  Note that this has the same format as a commit message, implying the first line is the title and the others are descriptions
#  Every line starting with one only '#' are considered as a comments

{{.Message}}

`
	templ = template.Must(template.New("user-input").Parse(baseContent))

	ErrEmptyMessage = errors.New("invalid empty message")
)

func FormalizeMessage(opts *models.Opts, repo *git.Repo) error {
	// get user's message if needed
	err := getUserInput(opts, repo)
	if err != nil {
		return err
	}

	// format the message
	opts.Message = formatMessage(opts.Message)

	// if the message is empty returns an error
	if len(opts.Message) == 0 {
		return ErrEmptyMessage
	}

	return nil
}

func getUserInput(opts *models.Opts, repo *git.Repo) error {
	// if there is already a message and we are not forced to edit it, we return
	if len(opts.Message) > 0 && !opts.ForceEdit {
		return nil
	}

	// get the base content of the editor
	buf := new(bytes.Buffer)
	err := templ.Execute(buf, struct{ Message string }{opts.Message})
	base := buf.Bytes()

	// get user's input through its editor
	content, err := getInputFromEditor(repo, base)
	if err != nil {
		return err
	}
	opts.Message = string(content)

	return nil
}

func getInputFromEditor(repo *git.Repo, base []byte) ([]byte, error) {
	editor, err := repo.GetEditor()
	if err != nil {
		return nil, err
	}
	editor = strings.TrimSpace(editor)
	editionLine := strings.Split(editor, "\n")

	tmpFile, err := ioutil.TempFile("", "mergo-*")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write(base); err != nil {
		return nil, err
	}

	editionLine = append(editionLine, tmpFile.Name())

	cmd := exec.Command(editionLine[0], editionLine[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	var content []byte
	if content, err = ioutil.ReadFile(tmpFile.Name()); err != nil {
		return nil, err
	}

	return content, nil
}

func formatMessage(msg string) string {
	// remove unused spaces
	msg = strings.TrimSpace(msg)

	tmp := []string{}
	for _, line := range strings.Split(msg, "\n") {
		// if this is a comment skip it ...
		trimed := strings.TrimSpace(line)
		if len(trimed) > 0 && trimed[0] == '#' {
			continue
		}

		// ... otherwise add it
		tmp = append(tmp, line)
	}

	return strings.TrimSpace(strings.Join(tmp, "\n"))
}

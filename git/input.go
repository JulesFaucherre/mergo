package git

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"gitlab.com/jfaucherre/mergo/tools"
)

// GetEditor returns the GitCmd to get the user's git editor
func (me *Repo) GetEditor() *GitCmd {
	next := func(editor string, err error) (string, error) {
		if !tools.IsEmpty(editor) {
			return editor, nil
		}

		if editor := os.Getenv("EDITOR"); editor != "" {
			return editor, nil
		}
		if editor := os.Getenv("VISUAL"); editor != "" {
			return editor, nil
		}

		return "vi", nil
	}

	return &GitCmd{
		repo: me,
		cmd: [][]string{
			{"git", "config", "--get", "core.editor"},
		},
		next: next,
	}
}

// EditText launches the user's configured editor and returns the text it has
// written in
func EditText(baseContent []byte) ([]byte, error) {
	ctx := context.Background()
	rawEditor, err := LocalRepository().GetEditor().Do(ctx)
	if err != nil {
		return nil, err
	}
	rawEditor = strings.TrimSpace(rawEditor)
	editor := strings.Split(rawEditor, "\n")

	tmpfile, err := ioutil.TempFile("", "mergo-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write(baseContent); err != nil {
		return nil, err
	}
	editor = append(editor, tmpfile.Name())

	cmd := exec.CommandContext(ctx, editor[0], editor[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return nil, err
	}

	var content []byte
	if content, err = ioutil.ReadFile(tmpfile.Name()); err != nil {
		return nil, err
	}

	return content, nil
}

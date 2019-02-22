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
// Note that lines starting with '#' will be considered as comments
func EditText(baseContent []byte) (string, error) {
	ctx := context.Background()
	rawEditor, err := LocalRepository().GetEditor().Do(ctx)
	if err != nil {
		return "", err
	}
	rawEditor = strings.TrimSpace(rawEditor)
	editor := strings.Split(rawEditor, "\n")

	tmpfile, err := ioutil.TempFile("", "mergo-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write(baseContent); err != nil {
		return "", err
	}
	editor = append(editor, tmpfile.Name())

	cmd := exec.CommandContext(ctx, editor[0], editor[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return "", err
	}

	var content []byte
	if content, err = ioutil.ReadFile(tmpfile.Name()); err != nil {
		return "", err
	}

	res := removeComments(string(content))

	if err = tmpfile.Close(); err != nil {
		return "", err
	}

	return res, nil
}

func removeComments(src string) string {
	dst := make([]string, 0)

	for _, s := range strings.Split(src, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(s), "#") {
			dst = append(dst, s)
		}
	}
	return strings.Join(dst, "\r\n")
}

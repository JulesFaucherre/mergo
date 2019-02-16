package git

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func GetEditor(ctx context.Context) (string, error) {
	cmd := [][]string{
		{"git", "config", "--global", "--get", "core.editor"},
	}
	if editor, _ := run(ctx, cmd); editor != "" {
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

func EditText(baseContent []byte) (string, error) {
	ctx := context.Background()
	editor, err := GetEditor(ctx)
	if err != nil {
		return "", err
	}
	editor = strings.TrimSpace(editor)

	tmpfile, err := ioutil.TempFile("", "mergo-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpfile.Name())

	if _, err = tmpfile.Write(baseContent); err != nil {
		return "", err
	}

	cmd := exec.CommandContext(ctx, editor, tmpfile.Name())
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

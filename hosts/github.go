package hosts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.com/jfaucherre/mergo/models"
)

type github struct{}

func newGithub() Host {
	return github{}
}

func (me github) GetOwnerAndRepo(remote string) (string, string) {
	return ownerAndRepo(remote)
}

func (me github) SubmitPr(opts models.Opts) error {
	body := struct {
		Head  string `json:"head"`
		Base  string `json:"base"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		opts.Head,
		opts.Base,
		opts.Title,
		opts.Body,
	}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/pulls",
		opts.Owner,
		opts.Repo,
	)

	marshaled, err := json.Marshal(body)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(marshaled)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Request failed with status %s", resp.Status)
	}
	return nil
}

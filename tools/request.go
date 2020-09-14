package tools

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.com/jfaucherre/mergo/logger"
)

type RequestParams struct {
	URL     string
	Method  string
	Body    interface{}
	Headers map[string]string
	Result  interface{}
}

// Request is just a helper function for the hosts to request their server
func Request(params *RequestParams) (int, error) {
	logger.Debug("request\nmethod: %s\nURL: %s\n", params.Method, params.URL)
	logger.Silly("body: %+v\nheaders: %+v\n", params.Body, params.Headers)
	marshaled, err := json.Marshal(params.Body)
	if err != nil {
		return 0, err
	}
	reader := bytes.NewReader(marshaled)

	client := &http.Client{}
	logger.Silly("sending request\n")
	req, err := http.NewRequest(params.Method, params.URL, reader)
	if err != nil {
		return 0, err
	}

	if len(params.Headers) == 0 {
		params.Headers = make(map[string]string)
	}
	params.Headers["Content-Type"] = "application/json"
	params.Headers["Accept"] = "application/json"
	for k, v := range params.Headers {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	logger.Info("request response %s\n", b)
	if err != nil {
		return 0, err
	}

	if err = json.Unmarshal(b, params.Result); err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

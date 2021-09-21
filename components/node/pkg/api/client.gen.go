// Code generated by oto; DO NOT EDIT.

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type nodeClient struct {
	username string
	password string
	endpoint string
	cl       *http.Client
}

func NewNodeClient(endpoint, username, password string, timeout time.Duration) Node {
	return &nodeClient{
		username: username,
		password: password,
		endpoint: endpoint,
		cl:       &http.Client{Timeout: timeout},
	}
}

func (c *nodeClient) ProbeReady(ctx context.Context, req ProbeReadyRequest) (*ProbeReadyResponse, error) {
	body, err := json.Marshal(&req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}
	request, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/oto/Node.ProbeReady", c.endpoint),
		bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "request")
	}
	if c.username != "" {
		request.SetBasicAuth(c.username, c.password)
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.cl.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "body")
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			return nil, errors.New(string(body))
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}
	response := &ProbeReadyResponse{}
	if err := json.Unmarshal(body, response); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}
	return response, nil
}

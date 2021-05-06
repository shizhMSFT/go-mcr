package mcr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
)

const endpoint = "https://mcr.microsoft.com/v2"

type client struct {
	base *http.Client
}

func NewClient(base *http.Client) Client {
	if base == nil {
		base = http.DefaultClient
	}
	return &client{
		base: base,
	}
}

func (c *client) Repositories(ctx context.Context) ([]string, error) {
	url := endpoint + "/_catalog"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.base.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var result struct {
		Repositories []string `json:"repositories"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	sort.Strings(result.Repositories)
	return result.Repositories, nil
}

func (c *client) Tags(ctx context.Context, repo string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/tags/list", endpoint, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.base.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var result struct {
		Tags []string `json:"tags"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	sort.Strings(result.Tags)
	return result.Tags, nil
}

func (c *client) Manifest(ctx context.Context, repo, ref string) (string, []byte, error) {
	url := fmt.Sprintf("%s/%s/manifests/%s", endpoint, repo, ref)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", nil, err
	}
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	resp, err := c.base.Do(req)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New(resp.Status)
	}

	mediaType := resp.Header.Get("Content-Type")
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}
	return mediaType, result, nil
}

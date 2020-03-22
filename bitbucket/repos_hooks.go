package bitbucket

import (
	"context"
	"fmt"
)

type WebHook struct {
	ID            int                  `json:"id,omitempty"`
	Name          string               `json:"name,omitempty"`
	CreatedDate   Time                 `json:"createdDate,omitempty"`
	UpdatedDate   Time                 `json:"updatedDate,omitempty"`
	Events        []string             `json:"events,omitempty"`
	Configuration WebHookConfiguration `json:"configuration,omitempty"`
	Url           string               `json:"url,omitempty"`
	Active        bool                 `json:"active,omitempty"`
}

type WebHookConfiguration struct {
	Secret string `json:"secret,omitempty"`
}

type WebHookListOptions struct {
	Event      string `url:"event,omitempty"`
	Statistics bool   `url:"statistics,omitempty"`
}

func (s *RepositoriesService) CreateWebHooks(ctx context.Context, projectKey, repositorySlug string, hook *WebHook) (*WebHook, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s/webhooks", projectKey, repositorySlug)

	req, err := s.client.NewRequest("POST", u, hook)
	if err != nil {
		return nil, nil, err
	}

	v := new(WebHook)
	resp, err := s.client.Do(ctx, req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}

// ListWebHooks finds web hooks in a repository.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp365
func (s *RepositoriesService) ListWebHooks(ctx context.Context, projectKey, repositorySlug string, opts *WebHookListOptions) ([]*WebHook, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s/webhooks", projectKey, repositorySlug)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var hooks []*WebHook
	page := &pagedResponse{
		Values: &hooks,
	}
	resp, err := s.client.Do(ctx, req, page)
	if err != nil {
		return nil, resp, err
	}

	return hooks, resp, nil
}

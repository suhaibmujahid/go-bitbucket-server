package bitbucket

import (
	"context"
	"fmt"
)

// RepositoriesService handles communication with the repository related
// methods of the  Bitbucket Server API.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp167
type RepositoriesService service

// Repository represents a Bitbucket Server repository.
type Repository struct {
	Slug          string           `json:"slug,omitempty"`
	Id            int              `json:"id,omitempty"`
	Name          string           `json:"name,omitempty"`
	Description   string           `json:"description,omitempty"`
	HierarchyId   string           `json:"hierarchyId,omitempty"`
	ScmId         string           `json:"scmId,omitempty"`
	State         string           `json:"state,omitempty"`
	StatusMessage string           `json:"statusMessage,omitempty"`
	Forkable      bool             `json:"forkable,omitempty"`
	Project       *Project         `json:"project,omitempty"`
	Public        bool             `json:"public,omitempty"`
	Links         *RepositoryLinks `json:"links,omitempty"`
}

type RepositoryLinks struct {
	Self  *NamelessLink `json:"self,omitempty"`
	Clone []*Link       `json:"clone,omitempty"`
}

func (s *RepositoriesService) List(ctx context.Context, project string, opts *ListOptions) ([]*Repository, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos", project)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	page := &pagedResponse{
		Values: &repos,
	}
	resp, err := s.client.Do(ctx, req, page)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

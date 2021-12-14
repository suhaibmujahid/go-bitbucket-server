package bitbucket

import (
	"context"
	"fmt"
)

const (
	PermissionRepoRead  = "REPO_READ"
	PermissionRepoWrite = "REPO_WRITE"
	PermissionRepoAdmin = "REPO_ADMIN"
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
	Origin        *Repository      `json:"origin,omitempty"` // this populated only for forked repositories
	Project       *Project         `json:"project,omitempty"`
	Public        bool             `json:"public,omitempty"`
	Links         *RepositoryLinks `json:"links,omitempty"`
}

func (r *Repository) String() string {
	return fmt.Sprintf("<Repository: %s/%s>", r.Project.Key, r.Slug)
}

type RepositoryLinks struct {
	Self  []NamelessLink `json:"self,omitempty"`
	Clone []Link         `json:"clone,omitempty"`
}

//todo: this could be merged with Ref, PullRequestTarget and PullRequestRef
type Branch struct {
	ID              string `json:"id"`
	DisplayID       string `json:"displayId"`
	Type            string `json:"type"`
	LatestCommit    string `json:"latestCommit"`
	LatestChangeset string `json:"latestChangeset"`
	IsDefault       bool   `json:"isDefault,omitempty"`
}

//todo: add the option values as consts
type ListRepositoriesOptions struct {
	// Name (optional) if specified, this will limit the resulting repository
	// list to ones whose name matches this parameter's value. The match will
	// be done case-insensitive and any leading and/or trailing whitespace
	// characters on the name parameter will be stripped.zxl
	Name string `url:"name,omitempty"`

	// ProjectName (optional) if specified, this will limit the resulting
	// repository list to ones whose project's name matches this parameter's value.
	// The match will be done case-insensitive and any leading and/or trailing
	// whitespace characters on the `projectname` parameter will be stripped.
	ProjectName string `url:"projectname,omitempty"`

	// Permission (optional) if specified, it must be a valid repository permission
	// level name and will limit the resulting repository list to ones that the
	// requesting user has the specified permission level to. If not specified,
	// the default implicit 'read' permission level will be assumed. The currently
	// supported explicit permission values are REPO_READ, REPO_WRITE and REPO_ADMIN.
	Permission string `url:"permission,omitempty"`

	// State (optional) if specified, it must be a valid repository state name
	// and will limit the resulting repository list to ones that are in the
	// specified state. The currently supported explicit state values are AVAILABLE,
	// INITIALISING and INITIALISATION_FAILED.
	State string `url:"state,omitempty"`

	// Visibility (optional) if specified, this will limit the resulting repository
	// list based on the repositories visibility. Valid values are public or private.
	Visibility string `url:"visibility,omitempty"`

	ListOptions
}

// Create repository in a project. To create a personal repository the projectKey
// should be ~ then user slug (e.g., ~suhaib).
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp168
func (s *RepositoriesService) Create(ctx context.Context, projectKey, repoName string) (*Repository, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos", projectKey)
	b := map[string]string{"name": repoName}

	req, err := s.client.NewRequest(ctx, "POST", u, b)
	if err != nil {
		return nil, nil, err
	}

	repo := new(Repository)
	resp, err := s.client.Do(req, &repo)
	if err != nil {
		return nil, resp, err
	}

	return repo, resp, nil
}

// Fork repository in a project. To fork a personal repository the projectKey
// should be ~ then user slug (e.g., ~suhaib).
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp174
func (s *RepositoriesService) Fork(ctx context.Context, projectKey, repoSlug, repoName string) (*Repository, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s", projectKey, repoSlug)
	b := map[string]interface{}{
		"name": repoName,
		"project": map[string]interface{}{"key": projectKey},
	}

	req, err := s.client.NewRequest(ctx, "POST", u, b)
	if err != nil {
		return nil, nil, err
	}

	repo := new(Repository)
	resp, err := s.client.Do(req, &repo)
	if err != nil {
		return nil, resp, err
	}

	return repo, resp, nil
}

// List retrieves a page of repositories based on the options that control the search.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp393
func (s *RepositoriesService) List(ctx context.Context, opts *ListRepositoriesOptions) ([]*Repository, *Response, error) {
	u, err := addOptions("repos", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	page := &pagedResponse{
		Values: &repos,
	}
	resp, err := s.client.Do(req, page)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// ListByProject the repositories for a project. To list personal repositories, projectKey
// should be ~ then user slug (e.g., ~suhaib).
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp169
func (s *RepositoriesService) ListByProject(ctx context.Context, projectKey string, opts *ListOptions) ([]*Repository, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos", projectKey)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	page := &pagedResponse{
		Values: &repos,
	}
	resp, err := s.client.Do(req, page)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// Get fetches a repository.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp172
func (s *RepositoriesService) Get(ctx context.Context, projectKey, repositorySlug string) (*Repository, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s", projectKey, repositorySlug)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	repo := new(Repository)
	resp, err := s.client.Do(req, repo)
	if err != nil {
		return nil, resp, err
	}

	return repo, resp, nil
}

type RecentReposOptions struct {
	// Permission (optional) if specified, it must be a valid repository permission
	// level name and will limit the resulting repository list to ones that the
	// requesting user has the specified permission level to. If not specified,
	// the default REPO_READ permission level will be assumed.
	Permission string `url:"permission,omitempty"`

	ListOptions
}

// ListRecent retrieves a page of recently accessed repositories for the currently authenticated user.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp140
func (s *RepositoriesService) ListRecent(ctx context.Context, opts *RecentReposOptions) ([]*Repository, *Response, error) {
	u, err := addOptions("profile/recent/repos", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var repos []*Repository
	page := &pagedResponse{
		Values: &repos,
	}
	resp, err := s.client.Do(req, page)
	if err != nil {
		return nil, resp, err
	}

	return repos, resp, nil
}

// GetDefaultBranch returns the default branch of the repository.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp204
func (s *RepositoriesService) GetDefaultBranch(ctx context.Context, projectKey, repositorySlug string) (*Branch, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s/branches/default", projectKey, repositorySlug)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	b := new(Branch)
	resp, err := s.client.Do(req, b)
	if err != nil {
		return nil, resp, err
	}

	return b, resp, nil
}

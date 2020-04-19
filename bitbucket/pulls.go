package bitbucket

import (
	"context"
	"fmt"
)

// PullRequestsService handles communication with the pull request related
// methods of the Bitbucket Server API.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp280
type PullRequestsService service

// PullRequest represents a Bitbucket Server pull request on a repository.
type PullRequest struct {
	ID           int                `json:"id,omitempty"`
	Version      int                `json:"version,omitempty"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	State        string             `json:"state,omitempty"`
	Open         bool               `json:"open,omitempty"`
	Closed       bool               `json:"closed,omitempty"`
	CreatedDate  Time               `json:"createdDate,omitempty"`
	UpdatedDate  Time               `json:"updatedDate,omitempty"`
	FromRef      *PullRequestRef    `json:"from_ref,omitempty"`
	ToRef        *PullRequestRef    `json:"toRef,omitempty"`
	Locked       bool               `json:"locked,omitempty"`
	Author       *PullRequestUser   `json:"author,omitempty"`
	Reviewers    []*PullRequestUser `json:"reviewers,omitempty"`
	Participants []*PullRequestUser `json:"participants,omitempty"`
	Links        *SelfLinks         `json:"links,omitempty"`
}

type PullRequestRef struct {
	ID           string      `json:"id,omitempty"`
	DisplayId    string      `json:"displayId,omitempty"`
	LatestCommit string      `json:"latestCommit,omitempty"`
	Repository   *Repository `json:"repository,omitempty"`
}

type PullRequestUser struct {
	User               *User  `json:"user,omitempty"`
	LastReviewedCommit string `json:"lastReviewedCommit,omitempty"` // this populated only for pull request reviewers
	Role               string `json:"role,omitempty"`
	Approved           bool   `json:"approved,omitempty"`
	Status             string `json:"status,omitempty"`
}

// PullRequestListOptions specifies the optional parameters to the
// PullRequestsService.List method.
type PullRequestListOptions struct {
	// Direction (optional, defaults to INCOMING) the direction relative to the specified repository.
	// Either INCOMING or OUTGOING.
	Direction string `url:"direction,omitempty"`

	// At (optional) a fully-qualified branch ID to find pull requests to or from,
	// such as {@code refs/heads/master}
	At string `url:"at,omitempty"`

	// State (optional, defaults to OPEN). Supply ALL to return pull request in any state.
	// If a state is supplied only pull requests in the specified state will be returned.
	// Either OPEN, DECLINED or MERGED.
	State string `url:"state,omitempty"`

	// Order (optional, defaults to NEWEST) the order to return pull requests in,
	// either OLDEST (as in: "oldest first") or NEWEST.
	Order string `url:"order,omitempty"`

	// WithAttributes (optional, defaults to true) whether to return additional pull request attributes
	WithAttributes bool `url:"withAttributes,omitempty"`

	// WithProperties (optional, defaults to true) whether to return additional pull request properties
	WithProperties bool `url:"withProperties,omitempty"`

	ListOptions
}

// List retrieves a page of pull requests to or from the specified repository.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp281
func (s *PullRequestsService) List(ctx context.Context, projectKey, repo string, opts *PullRequestListOptions) ([]*PullRequest, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s/pull-requests", projectKey, repo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var pulls []*PullRequest
	page := &pagedResponse{
		Values: &pulls,
	}
	resp, err := s.client.Do(req, page)
	if err != nil {
		return nil, resp, err
	}

	return pulls, resp, nil
}

// Get retrieves a single pull request.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp284
func (s *PullRequestsService) Get(ctx context.Context, projectKey, repo string, id int) (*PullRequest, *Response, error) {
	u := fmt.Sprintf("projects/%s/repos/%s/pull-requests/%v", projectKey, repo, id)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	pull := new(PullRequest)
	resp, err := s.client.Do(req, pull)
	if err != nil {
		return nil, resp, err
	}

	return pull, resp, nil
}

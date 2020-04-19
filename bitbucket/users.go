package bitbucket

import (
	"bytes"
	"context"
	"fmt"
	"github.com/google/go-querystring/query"
	"net/url"
)

// UsersService handles communication with the user related
// methods of the  Bitbucket Server API.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp400
type UsersService service

type User struct {
	Name         string     `json:"name,omitempty"`
	EmailAddress string     `json:"emailAddress,omitempty"`
	Id           int        `json:"id,omitempty"`
	DisplayName  string     `json:"displayName,omitempty"`
	Active       bool       `json:"active,omitempty"`
	Slug         string     `json:"slug,omitempty"`
	Type         string     `json:"type,omitempty"`
	Links        *SelfLinks `json:"links,omitempty"`
}

// WhoAmI use the `whoami` endpoint to retrieve the current
// authenticated user slug as string.
func (s *UsersService) WhoAmI(ctx context.Context) (string, *Response, error) {
	// we use slash at the beginning in this case to avoid having the URL relative to the suffix `rest/api/1.0/`
	req, err := s.client.NewRequest(ctx, "GET", "/plugins/servlet/applinks/whoami", nil)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(req, buf)
	if err != nil {
		return "", resp, err
	}

	return buf.String(), resp, nil
}

// Myself performs two requests, the first one is through WhoAmI method to know the
// current user and the second one through Get method to retrieve the user info.
func (s *UsersService) Myself(ctx context.Context) (*User, *Response, error) {
	slug, resp, err := s.client.Users.WhoAmI(ctx)
	if err != nil {
		return nil, resp, err
	}

	return s.client.Users.Get(ctx, slug)
}

// Get retrieves the user matching the supplied userSlug.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp406
func (s *UsersService) Get(ctx context.Context, slug string) (*User, *Response, error) {
	u := fmt.Sprintf("users/%s", slug)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, nil
}

type ListUsersPermissions []ListUsersPermission

func (l ListUsersPermissions) EncodeValues(key string, v *url.Values) error {
	if len(l) == 0 {
		return nil
	}

	qu, err := query.Values(l[0])
	if err != nil {
		return err
	}
	v.Add(key, qu.Get(key))
	qu.Del(key)
	appendOptions(key, v, qu)

	for i := 1; i < len(l); i++ {
		prefix := fmt.Sprintf("%s.%d", key, i)
		qu, err := query.Values(l[i])
		if err != nil {
			return err
		}
		v.Add(prefix, qu.Get(key))
		qu.Del(key)
		appendOptions(prefix, v, qu)
	}

	return nil
}

func appendOptions(prefix string, v *url.Values, new url.Values) {
	for name := range new {
		v.Add(fmt.Sprintf("%s.%s", prefix, name), new.Get(name))
	}
}

type ListUsersOptions struct {
	// Filter return only users, whose username, name or email address contain this value.
	Filter string `url:"filter,omitempty"`

	// Group return only users who are members of the given group.
	Group string `url:"group,omitempty"`

	// Permission return users who have the specified permissions.
	Permission ListUsersPermissions `url:"permission,omitempty"`

	ListOptions
}

type ListUsersPermission struct {
	Permission     string `url:"permission"`
	ProjectKey     string `url:"projectKey,omitempty"`
	RepositorySlug string `url:"repositorySlug,omitempty"`
	RepositoryId   string `url:"repositoryId,omitempty"`
}

// List retrieves a page of users, optionally run through provided filters.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp401
func (s *UsersService) List(ctx context.Context, opts *ListUsersOptions) ([]*User, *Response, error) {
	u, err := addOptions("users", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	page := &pagedResponse{
		Values: &users,
	}
	resp, err := s.client.Do(req, page)
	if err != nil {
		return nil, resp, err
	}

	return users, resp, nil
}

package bitbucket

import (
	"bytes"
	"context"
	"fmt"
)

// UsersService handles communication with the user related
// methods of the  Bitbucket Server API.
//
// Bitbucket Server API doc: https://docs.atlassian.com/bitbucket-server/rest/7.0.1/bitbucket-rest.html#idp400
type UsersService service

type User struct {
	Name         string `json:"name,omitempty"`
	EmailAddress string `json:"emailAddress"`
	Id           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

// WhoAmI use the `whoami` endpoint to retrieve the current
// authenticated user slug as string.
func (s *UsersService) WhoAmI(ctx context.Context) (string, *Response, error) {
	// we use slash at the beginning in this case to avoid having the URL relative to the suffix `rest/api/1.0/`
	req, err := s.client.NewRequest("GET", "/plugins/servlet/applinks/whoami", nil)
	if err != nil {
		return "", nil, err
	}

	buf := new(bytes.Buffer)
	resp, err := s.client.Do(ctx, req, buf)
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

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	uResp := new(User)
	resp, err := s.client.Do(ctx, req, uResp)
	if err != nil {
		return nil, resp, err
	}

	return uResp, resp, nil
}

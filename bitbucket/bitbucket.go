package bitbucket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

const (
	userAgent = "go-bitbucket-server"
)

// A Client manages communication with the Bitbucket Server API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// User agent used when communicating with the Bitbucket Server API.
	UserAgent string

	common service

	// Base URL for API requests.
	baseURL *url.URL

	// Services used for talking to different parts of the Bitbucket Server API.
	Users        *UsersService
	PullRequests *PullRequestsService
}

func (c *Client) BaseURL() url.URL {
	return *c.baseURL
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// Start parameter indicates which item should be used as the first item in the page of results.
	Start int `url:"start,omitempty"`

	// Limit parameter indicates how many results to return per page.
	Limit int `url:"limit,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewServerClient returns a new Bitbucket Server API client with provided base URL.
// If either URL does not have the suffix "/rest/api/1.0/", it will be added automatically.
// If a nil httpClient is provided, a new http.Client will be used.
func NewServerClient(baseURL string, httpClient *http.Client) (*Client, error) {
	baseEndpoint, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}
	if !strings.HasSuffix(baseEndpoint.Path, "/rest/api/1.0/") {
		baseEndpoint.Path += "rest/api/1.0/"
	}

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	c := &Client{client: httpClient, baseURL: baseEndpoint, UserAgent: userAgent}
	c.common.client = c
	c.Users = (*UsersService)(&c.common)
	c.PullRequests = (*PullRequestsService)(&c.common)
	return c, nil
}

type service struct {
	client *Client
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash, otherwise
// the URL will be relative root of the base URL (ignoring the API suffix i.e., `/rest/api/1.0/`).
// If specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp, v)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
func CheckResponse(resp *http.Response) error {
	if c := resp.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return fmt.Errorf("Request failed. Please analyze the request body for more details. Status code: %d", resp.StatusCode)
}

// Response represents Bitbucket Server API response. It wraps http.Response returned from
// API and provides information about paging.
type Response struct {
	*http.Response

	*pagedResponse
}

type pagedResponse struct {
	Size          int         `json:"size"`
	Limit         int         `json:"limit"`
	IsLastPage    bool        `json:"isLastPage"`
	Values        interface{} `json:"values"`
	Start         int         `json:"start"`
	Filter        string      `json:"filter,omitempty"`
	NextPageStart int         `json:"nextPageStart,omitempty"`
}

func newResponse(r *http.Response, v interface{}) *Response {
	resp := &Response{Response: r}
	resp.populatePageValues(v)
	return resp
}

var notPagedResponse = pagedResponse{
	Size:       1,
	Limit:      1,
	IsLastPage: true,
}

// Sets paging values if response json was parsed to pagedResponse type
func (r *Response) populatePageValues(v interface{}) {
	switch value := v.(type) {
	case *pagedResponse:
		r.pagedResponse = value
	default:
		r.pagedResponse = &notPagedResponse
	}
}

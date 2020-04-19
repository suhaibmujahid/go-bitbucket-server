# go-bitbucket-server
[![Go](https://github.com/suhaibmujahid/go-bitbucket-server/workflows/Go/badge.svg)](https://github.com/suhaibmujahid/go-bitbucket-server/actions?query=workflow%3AGo)

go-bitbucket-server is a [Go](https://golang.orgs) client library for accessing the 
[Bitbucket](https://www.atlassian.com/software/bitbucket) Server API v1.

> This client is under development , I started with endpoints that is important to me.
> The goal eventually is to cover the all endpoints.
> Pull requests are welcomed to implement missed endpoints or fix bugs.  

The interface design of this package was inspired by [google/go-github](https://github.com/google/go-github) 
and [andygrunwald/go-jira](https://github.com/andygrunwald/go-jira).


## Usage

```go
import "github.com/suhaibmujahid/go-bitbucket-server/bitbucket"

client := bitbucket.NewServerClient("http://localhost:7990/rest/api/1.0/", http.DefaultClient)

// retrieve the authenticated user
user, _, err := client.Users.Myself(context.Background())
```

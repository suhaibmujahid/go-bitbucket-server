# go-bitbucket-server
go-bitbucket-server is a Go client library for accessing the Bitbucket Server API v1.

This client is under development , I started with endpoints that is important to me. The goal eventually is to cover the all endpoints. Pull requests are welcomed to implement missed endpoints or fix bugs.  

 ## Usage
 
```go
import "github.com/suhaibmujahid/go-bitbucket-server/bitbucket"

client := bitbucket.NewServerClient("http://localhost:7990/rest/api/1.0/", http.DefaultClient)

// retrieve the authenticated user
user, _, err := client.Users.Myself(context.Background())
```
package crawler

import (
	"net/http"
    "fmt"
	"github.com/juju/ratelimit"
    "strings"
)
/*
This is a simple framework for making ratelimited API calls.
*/

/* Parses http response to extract entities */
type Entity interface {
    ParseResponse(resp *http.Response) error
}

/* Any API call that is ratelimited */
type RateLimitedAPI struct {
    apiUrl string
    limiter *ratelimit.Bucket
}

/* Call a ratelimited API */
func (a *RateLimitedAPI) CallAPI(e Entity) error {
    fmt.Printf("get api\n")
    a.limiter.Wait(1)
    resp, err := http.Get(a.apiUrl)
    defer resp.Body.Close()
    if err != nil {
        return err
    }
    return e.ParseResponse(resp)
}

/* factory method for RateLimitedAPI */
func NewRateLimitedAPI(apiUrl string) *RateLimitedAPI {
    return &RateLimitedAPI{apiUrl, ratelimit.NewBucketWithRate(30, 30)}
}

/* helper function to extract next page link */
func ExtractNextLink(link string) string {
	next := strings.TrimLeft(link, "<")
	next = next[0:strings.LastIndex(next, ">")]
	return next
}
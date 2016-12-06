package main
import (
  "net/http"
  "fmt"
  "strings"
  "encoding/json"
  "github.com/juju/ratelimit"
)

func main() {
  // resp, err := http.Get("https://api.meetup.com/Women-Who-Code-SF/events?order=created&desc=1&status=past&page=5")
  StoreEvents("blah")
}

func Events(page int) string {
  //apikey should go here
  const apikey string = "blah"

  apiUrl := fmt.Sprintf("https://api.meetup.com/Women-Who-Code-SF/events?order=created&desc=1&status=past&page=%d)", page)
  resp, err := http.Get(apiUrl)
  if err != nil {
	// handle error
  }
  defer resp.Body.Close()

  // sanitize trim the angle brackets from left and right
  next := extractLink(resp.Header.Get("Link"))
  // next := strings.TrimLeft(resp.Header.Get("Link"), "<")
  // next = next[0: strings.LastIndex(next, ">")]
  fmt.Printf("======NEXT URL = %s", next)

  // decode json
  var dat interface{}
  json.NewDecoder(resp.Body).Decode(&dat)
  // array of events
  events := dat.([]interface{})

  for _, event := range events {
    //for each event
    tEvent := event.(map[string]interface{})
    link := tEvent["link"].(string)
    linkTrimmed := strings.TrimSuffix(link, "/")
    //parse the link to get the id as the last part of the url
    id :=  linkTrimmed[strings.LastIndex(linkTrimmed, "/") + 1:]
    fmt.Println(link)
    fmt.Println(id)
    fmt.Println(tEvent["name"].(string))
  }
  return next
}

func extractLink(link string) string {
  next := strings.TrimLeft(link, "<")
  next = next[0: strings.LastIndex(next, ">")]
  return next
}

func StoreEvents(lastWatermark string) {
  // for next := Events(page); next != ""; {

  // }
  const page int = 10
  // rate limit to 30 requests per second
  limiter := ratelimit.NewBucketWithRate(30, 30)
  for i := 0; i < 2; i++ {
    next := Events(page)
    fmt.Println(i)
    fmt.Println(next)
    limiter.Wait(1)
  }
}

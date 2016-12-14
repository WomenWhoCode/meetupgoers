package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/juju/ratelimit"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/WomenWhoCode/meetupgoers/mongodb"
)

const (
	DefaultWaterMark = "0"
)

type Event struct {
	id   string
	name string
	url  string
}

type WaterMark struct {
	Name      string
	ID        string
	timeStamp time.Time
}

var firstID string = "0"
var AuthDatabase string = os.Getenv("DBName")

func main() {
	StartTheEngine()
}

func StartTheEngine() string {
	// resp, err := http.Get("https://api.meetup.com/Women-Who-Code-SF/events?order=created&desc=1&status=past&page=5")
	waterMark := findWaterMark()
	StoreEvents(waterMark)
	return "success"
}

func findWaterMark() string {
	session := mongodb.ConnDB()
	var result WaterMark
	err := session.DB(AuthDatabase).C("waterMark").Find(bson.M{}).One(&result)
	if err != nil {
		session.DB(AuthDatabase).C("waterMark").Insert(&WaterMark{"event", DefaultWaterMark, time.Now()})
		return DefaultWaterMark
	}
	return result.ID
}

func updateWaterMark(ID string) {
	session := mongodb.ConnDB()
	change := bson.M{"$set": bson.M{"id": ID, "timestamp": time.Now()}}
	err := session.DB(AuthDatabase).C("waterMark").Update(bson.M{"name": "event"}, change)
	if err != nil {
		panic(err)
	}
}

func Events(apiUrl string, waterMark string) string {
	//apikey should go here
	const apikey string = "blah"

	resp, err := http.Get(apiUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	// sanitize trim the angle brackets from left and right
	next := extractLink(resp.Header.Get("Link"))
	// next := strings.TrimLeft(resp.Header.Get("Link"), "<")
	// next = next[0: strings.LastIndex(next, ">")]

	// decode json
	var dat interface{}
	json.NewDecoder(resp.Body).Decode(&dat)

	// array of events
	events := dat.([]interface{})
	session := mongodb.ConnDB()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(AuthDatabase).C("eventCollection")
	for _, event := range events {
		//for each event
		tEvent := event.(map[string]interface{})
		link := tEvent["link"].(string)
		linkTrimmed := strings.TrimSuffix(link, "/")
		//parse the link to get the id as the last part of the url
		id := linkTrimmed[strings.LastIndex(linkTrimmed, "/")+1:]
		if id == waterMark {
			// reach checkpoint
			return ""
		}
		if firstID == "0" {
			firstID = id
		}
		err = c.Insert(&Event{tEvent["name"].(string), id, link})
		if err != nil {
			log.Fatal(err)
		}
	}

	return next
}

func extractLink(link string) string {
	next := strings.TrimLeft(link, "<")
	next = next[0:strings.LastIndex(next, ">")]
	return next
}

func StoreEvents(lastWatermark string) {
	// for next := Events(page); next != ""; {

	// }
	const page int = 10
	apiUrl := fmt.Sprintf("https://api.meetup.com/Women-Who-Code-SF/events?order=created&desc=1&status=past&page=%d)", page)
	// rate limit to 30 requests per second
	limiter := ratelimit.NewBucketWithRate(30, 30)
	for next := Events(apiUrl, lastWatermark); next != ""; {
		fmt.Printf("======NEXT URL : %s ======\n", next)
		next = Events(next, lastWatermark)
		if strings.Contains(next, "before") {
			break
		}
		limiter.Wait(1)
	}
	if firstID != "0" {
		updateWaterMark(firstID)
	}
}

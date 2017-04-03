package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/WomenWhoCode/meetupgoers/mongodb"
)

const (
	DefaultWaterMark = "0"
	AuthDatabase     = "heroku_vb4zpgmk"
)

type Event struct {
	ID   string
	NAME string
	URL  string
}

type EventResponse struct {
	eventlist    []interface{}
	nextPageUrl string
}

type WaterMark struct {
	Name      string
	ID        string
	TIMESTAMP time.Time
	GROUPNAME string
}

var firstID string = "0"

func StartTheEngine() string {
	// resp, err := http.Get("https://api.meetup.com/Women-Who-Code-SF/events?order=created&desc=1&status=past&page=5")
	fmt.Printf("Start crawling\n")
	StoreEvents()
	return "success"
}

func findWaterMark(groupName string) string {
	session := mongodb.ConnDB()
	collectionName := fmt.Sprintf("waterMark")
	var result WaterMark
	err := session.DB(AuthDatabase).C(collectionName).Find(bson.M{"groupname": groupName}).One(&result)
	if err != nil {
		session.DB(AuthDatabase).C(collectionName).Insert(&WaterMark{
			"event", DefaultWaterMark,
			time.Now(), groupName})
		return DefaultWaterMark
	}
	print(result.ID)
	return result.ID
}

func updateWaterMark(ID string, groupName string) {
	session := mongodb.ConnDB()
	change := bson.M{"$set": bson.M{"id": ID, "timestamp": time.Now(), "groupname": groupName}}
	collectionName := fmt.Sprintf("waterMark")
	err := session.DB(AuthDatabase).C(collectionName).Update(bson.M{"name": "event", "groupname": groupName}, change)
	if err != nil {
		panic(err)
	}
}

func (eventResp *EventResponse) ParseResponse(resp *http.Response) error {
	var dat interface{}
	err := json.NewDecoder(resp.Body).Decode(&dat)
	eventResp.eventlist = dat.([]interface{})
	link := resp.Header.Get("Link")
	if strings.Contains(link, "prev") || len(link) == 0 {
		eventResp.nextPageUrl = ""
	} else {
		eventResp.nextPageUrl = extractLink(link)
	}
	return err
}

func Events(apiUrl string, waterMark string) string {
	api := NewRateLimitedAPI(apiUrl)
	eventResponse := EventResponse{}
	api.CallAPI(&eventResponse)

	session := mongodb.ConnDB()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(AuthDatabase).C("eventCollection")
	for _, event := range eventResponse.eventlist {
		//for each event
		tEvent := event.(map[string]interface{})
		link := tEvent["link"].(string)
		linkTrimmed := strings.TrimSuffix(link, "/")
		//parse the link to get the id as the last part of the url
		id := linkTrimmed[strings.LastIndex(linkTrimmed, "/")+1:]
		if firstID == "0" {
			firstID = id
		}
		if id == waterMark {
			// reach checkpoint
			fmt.Printf("Encouter the same event: %s\n", waterMark)
			return ""
		}
		err := c.Insert(&Event{ID: id, NAME: tEvent["name"].(string), URL: link})
		if err != nil {
			log.Fatal(err)
		}
	}

	return eventResponse.nextPageUrl
}

func extractLink(link string) string {
	next := strings.TrimLeft(link, "<")
	next = next[0:strings.LastIndex(next, ">")]
	return next
}

func StoreEventsByGroup(lastWatermark string, groupName string) {
	const page int = 10
	apiUrl := fmt.Sprintf("https://api.meetup.com/%s/events?order=created&desc=1&status=past&page=%d)", groupName, page)
	for next := Events(apiUrl, lastWatermark); next != ""; {
		fmt.Printf("======NEXT URL : %s ======\n", next)
		next = Events(next, lastWatermark)
		if strings.Contains(next, "before") {
			break
		}
	}
	if firstID != "0" {
		updateWaterMark(firstID, groupName)
	}
}

func StoreEvents() {
	session := mongodb.ConnDB()
	session.SetMode(mgo.Monotonic, true)

	var groups []Group
	err := session.DB(AuthDatabase).C("groupCollection").Find(nil).All(&groups)
	if err != nil {
		panic(err)
	}
	// Iterate each group to get events
	for _, group := range groups {
		waterMark := findWaterMark(group.URLNAME)
		fmt.Printf("WaterMark found: %s for %s\n", waterMark, group.URLNAME)
		StoreEventsByGroup(waterMark, group.URLNAME)
	}

}

package crawler

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/WomenWhoCode/meetupgoers/mongodb"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"os"
)

type Group struct {
	AVG_AGE float64
	CITY string
	COUNTRY string
	DESCRIPTION string
	FOUNDED_DATE float64
	GEN_MALE float64
	GEN_FEMAILE float64
	GEN_OTHER float64
	GEN_UNKNOWN float64
	ID float64
	LAST_EVENT string
	LAT float64
	LON float64
	MEMBER_COUNT float64
	NAME string
	NEXT_EVENT string
	PAST_EVENTS float64
	PAST_RSVPS float64
	PRO_JOIN_DATE string
	REPEAT_RSVPERS float64
	RSVPS_PER_EVENT float64
	STATE string
	URLNAME string
}

type GroupResponse struct {
	grouplist    []interface{}
	nextPageUrl string
}

func StartTheGroupEngine() string {
	fmt.Printf("Start crawling group info\n")
	// Get API key from env
	key := os.Getenv("MEETUP_API_KEY")
	refreshGroups(key)
	return "success"
}

func (groupResp *GroupResponse) ParseResponse(resp *http.Response) error {
	var dat interface{}
	err := json.NewDecoder(resp.Body).Decode(&dat)
	groupResp.grouplist = dat.([]interface{})
	link := resp.Header.Get("Link")
	next := extractLink(link)
	if strings.Contains(link, "prev") {
		next = ""
	}
	groupResp.nextPageUrl = next
	return err
}

func groups(apiUrl string) string {

	api := NewRateLimitedAPI(apiUrl)

	groupResponse := GroupResponse{}
	api.CallAPI(&groupResponse)

	session := mongodb.ConnDB()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB(AuthDatabase).C("groupCollectionTMP")

	for _, group := range groupResponse.grouplist {
		//for each event
		tGroup := group.(map[string]interface{})
		group := populateGroup(tGroup)
		err := c.Insert(&group)
		if err != nil {
			log.Fatal(err)
		}
	}
	return groupResponse.nextPageUrl
}

func populateGroup(group map[string]interface{}) Group {
	return Group{
		AVG_AGE: group["average_age"].(float64),
		CITY: group["city"].(string),
		COUNTRY: group["country"].(string),
		DESCRIPTION: group["description"].(string),
		FOUNDED_DATE: group["founded_date"].(float64),
		GEN_MALE: group["gender_male"].(float64),
		GEN_FEMAILE: group["gender_female"].(float64),
		GEN_OTHER: group["gender_other"].(float64),
		GEN_UNKNOWN: group["gender_unknown"].(float64),
		ID: group["id"].(float64),
		LAT: group["lat"].(float64),
		LON: group["lon"].(float64),
		MEMBER_COUNT: group["member_count"].(float64),
		NAME: group["name"].(string),
		PAST_EVENTS: group["past_events"].(float64),
		PAST_RSVPS: group["past_rsvps"].(float64),
		REPEAT_RSVPERS: group["repeat_rsvpers"].(float64),
		RSVPS_PER_EVENT: group["rsvps_per_event"].(float64),
		URLNAME: group["urlname"].(string),
	}
}

func refreshGroups(key string) {
	const page int = 10
	apiUrl := fmt.Sprintf("https://api.meetup.com/pro/womenwhocode/groups?sign=true&key=%s&page=%d", key, page)
	for next := groups(apiUrl); next != ""; {
		next = fmt.Sprintf(next + "&key=%s", key)
		next = groups(next)
	}
	session := mongodb.ConnDB()
	// Rename collections for updating
	mongodb.RenameCollection("groupCollectionTMP", "groupCollection", session)
}


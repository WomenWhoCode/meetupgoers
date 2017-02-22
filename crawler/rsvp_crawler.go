package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/WomenWhoCode/meetupgoers/mongodb"
	mgo "gopkg.in/mgo.v2"
)

type Answer struct {
	ID          string
	EVENT_ID    string
	QUESTION_ID string
	ANSWER      string
	QUESTION    string
	MEMBER_ID   string
	GROUP_ID    string
}

type RsvpResponse struct {
	rsvplist    []interface{}
	nextPageUrl string
}

/* Parses response to RSVP api and Satisfies the Entity interface */
func (rsvpResp *RsvpResponse) ParseResponse(resp *http.Response) error {
	var dat interface{}
	d := json.NewDecoder(resp.Body)
	d.UseNumber()
	err := d.Decode(&dat)
	rsvpResp.rsvplist = dat.([]interface{})
	return err
}

var apikey string = os.Getenv("MEETUP_API_KEY")

func GetAnswers(eventId string) string {
	fmt.Printf("get rsvp\n")
	//generate api url
	apiUrl := fmt.Sprintf("https://api.meetup.com/women-who-code-sf/events/%s/rsvps?fields=answers&key=%s", eventId, apikey)
	fmt.Printf(apiUrl)
	//make a new api object
	api := NewRateLimitedAPI(apiUrl)
	rsvpResponse := RsvpResponse{}
	//Call the api. If this api was paginated, use a for loop
	//for rsvpResponse.nextPageUrl != ''
	api.CallAPI(&rsvpResponse)

	session := mongodb.ConnDB()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(AuthDatabase).C("answersCollection")

	for _, rsvp := range rsvpResponse.rsvplist {
		tRsvp := rsvp.(map[string]interface{})
		groupId := string(tRsvp["group"].(map[string]interface{})["id"].(json.Number))
		memberId := string(tRsvp["member"].(map[string]interface{})["id"].(json.Number))
		fmt.Printf("%s\n", memberId)
		if tRsvp["answers"] == nil {
			fmt.Printf("No answer in rsvp")
			continue
		}
		tAnswers := tRsvp["answers"].([]interface{})
		for _, ans := range tAnswers {
			tAns := ans.(map[string]interface{})
			ansText := tAns["answer"].(string)
			question := tAns["question"].(string)
			questionId := string(tAns["question_id"].(json.Number))
            //TODO sanity check for not null
			generatedAnsId := eventId + "|" + memberId + "|" + questionId
			a := Answer{
				ID:          generatedAnsId,
				EVENT_ID:    eventId,
				QUESTION_ID: questionId,
				QUESTION:    question,
				ANSWER:      ansText,
				MEMBER_ID:   memberId,
				GROUP_ID:    groupId,
			}
			err := c.Insert(&a)
			if err != nil {
				log.Fatal(err)
			}

		}
	}
	return "success"
}

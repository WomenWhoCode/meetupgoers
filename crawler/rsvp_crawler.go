package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Answer struct {
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

func GetAnswers(event_id string) string {
	fmt.Printf("get rsvp\n")
	//generate api url
	apiUrl := fmt.Sprintf("https://api.meetup.com/women-who-code-sf/events/%s/rsvps?key=%s", event_id, apikey)
	fmt.Printf(apiUrl)
	//make a new api object
	api := NewRateLimitedAPI(apiUrl)
	rsvpResponse := RsvpResponse{}
	//Call the api. If this api was paginated, use a for loop
	//for rsvpResponse.nextPageUrl != ''
	api.CallAPI(&rsvpResponse)
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
			tAns := ans.(map[string]string)
			fmt.Printf("Got an answer")
			ansText := tAns["answer"]
			question := tAns["question"]
			question_id := tAns["question_id"]
			a := Answer{
				EVENT_ID:    "blah",
				QUESTION_ID: question_id,
				QUESTION:    question,
				ANSWER:      ansText,
				MEMBER_ID:   memberId,
				GROUP_ID:    groupId,
			}
			fmt.Printf("%s", a)
			//TODO store this in mongo
		}
	}
	return "success"
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/WomenWhoCode/meetupgoers/crawler"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "scheduler" {
		crawler.StartTheEngine()
		crawler.StartTheGroupEngine()
	} else if len(os.Args) >= 2 && os.Args[1] == "rsvp_test" {
		event_id := "slmrkmywdbcc"
		crawler.GetAnswers(event_id)
	} else {
		http.HandleFunc("/", root)
		http.ListenAndServe(GetPort(), nil)
	}
}

// GetPort from the environment so we can run on Heroku
func GetPort() string {
	var port = os.Getenv("PORT")
	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, crawler.StartTheEngine())
}

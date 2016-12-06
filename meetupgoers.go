package nomain

import (
  "fmt"
  "net/http"
  "os"
)

func main() {
  http.HandleFunc("/", root)
  http.ListenAndServe(GetPort(), nil)
}

// Get the Port from the environment so we can run on Heroku
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
    fmt.Fprint(w, "Meetup Go'ers")
}
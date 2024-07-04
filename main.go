package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    ld "github.com/launchdarkly/go-server-sdk/v7"
    "github.com/launchdarkly/go-server-sdk/v7/ldcomponents"
    "github.com/launchdarkly/go-sdk-common/v3/ldlog"
    "github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

var ldClient *ld.LDClient

func main() {
    // Set custom config with a longer timeout and info logging level
    config := ld.Config{
        DataSource: ldcomponents.StreamingDataSource(),
        Events:     ldcomponents.SendEvents(),
        Logging:    ldcomponents.Logging().MinLevel(ldlog.Info),
    }

    // Initialize LaunchDarkly client with custom config and 10 second timeout
    var err error
    ldClient, err = ld.MakeCustomClient("sdk-00f6231c-4043-43e2-93bb-e9eab88a6d6b", config, 10*time.Second)
    if err != nil {
        log.Fatalf("Error initializing LaunchDarkly client: %s", err)
    }
    defer ldClient.Close()

    // HTTP handler for the home page
    http.HandleFunc("/", homeHandler)
    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Home page handler
func homeHandler(w http.ResponseWriter, r *http.Request) {
    user := ldcontext.New("example-user-key")

    // Check feature flag
    flagValue, err := ldClient.BoolVariation("new-homepage", user, false)
    if err != nil {
        log.Printf("Error reading feature flag: %s", err)
        http.Error(w, "Error reading feature flag", http.StatusInternalServerError)
        return
    }

    if flagValue {
        fmt.Fprintln(w, "Welcome to the new homepage of Hugo's Pet Shop!")
    } else {
        fmt.Fprintln(w, "Welcome to Hugo's Pet Shop!")
    }
}
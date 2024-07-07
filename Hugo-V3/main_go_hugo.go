package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    ld "github.com/launchdarkly/go-server-sdk/v7"
    "github.com/launchdarkly/go-server-sdk/v7/ldcomponents"
    "github.com/launchdarkly/go-sdk-common/v3/ldcontext"
    "github.com/launchdarkly/go-sdk-common/v3/ldvalue"
)

    var ldClient *ld.LDClient
    // FeatureFlags holds the feature flag states with a mutex for concurrent access
    type FeatureFlags struct {
    instantRollback bool
    newShopFeature  bool
    v3Feature       bool // Add this line
    mu              sync.RWMutex
    }
    var featureFlags = &FeatureFlags{}
    // sseClients manages the set of SSE clients
    var sseClients = struct {
    mu      sync.Mutex
    clients map[chan bool]struct{}
    }{clients: make(map[chan bool]struct{})}

    var jokes = []string{
    "Why did the Golden Retriever sit in the shade? He didn't want to be a hot dog!",
    "What do you get when you cross a Golden Retriever and a telephone? A golden receiver!",
    "What do you call a frozen dog? A pupsicle!",
}

    // LaunchDarkly configuration
    func main() {
    config := ld.Config{
        Events: ldcomponents.SendEvents(),
    }
    var err error
    ldClient, err = ld.MakeCustomClient("sdk-00f6231c-4043-43e2-93bb-e9eab88a6d6b", config, 10*time.Second)
    if err != nil {
        log.Fatalf("Error initializing LaunchDarkly client: %s", err)
    }
    defer ldClient.Close()

      // Create new user context named Hugo with additional attributes
      // Replace "Hugo" with your own user key and define attributes as needed
    user := ldcontext.NewBuilder("Hugo").
        SetString("firstName", "Hugo").
        SetString("lastName", "Henderson-Mills").
        SetInt("age", 3).
        SetString("location", "Worcester Park").
        SetString("breed", "Golden Retriever").
        Build()

    featureFlags.mu.Lock()
    featureFlags.instantRollback, _ = ldClient.BoolVariation("instant-rollback", user, false)
    featureFlags.newShopFeature, _ = ldClient.BoolVariation("new-shop-feature", user, false)
    featureFlags.v3Feature, _ = ldClient.BoolVariation("v3-feature", user, false)
    featureFlags.mu.Unlock()

    setupFlagListeners(user)

    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/events", sseHandler)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

    func setupFlagListeners(user ldcontext.Context) {
    updateCh1 := ldClient.GetFlagTracker().AddFlagValueChangeListener("instant-rollback", user, ldvalue.Bool(false))
    go func() {
        for event := range updateCh1 {
            if event.Key == "instant-rollback" {
                featureFlags.mu.Lock()
                featureFlags.instantRollback = event.NewValue.BoolValue()
                featureFlags.mu.Unlock()
                notifyClients()
                log.Printf("Flag %q for context %q has changed from %s to %s", event.Key, user.Key(), event.OldValue, event.NewValue)
            }
        }
    }()

    updateCh2 := ldClient.GetFlagTracker().AddFlagValueChangeListener("new-shop-feature", user, ldvalue.Bool(false))
    go func() {
        for event := range updateCh2 {
            if event.Key == "new-shop-feature" {
                featureFlags.mu.Lock()
                featureFlags.newShopFeature = event.NewValue.BoolValue()
                featureFlags.mu.Unlock()
                notifyClients()
                log.Printf("Flag %q for context %q has changed from %s to %s", event.Key, user.Key(), event.OldValue, event.NewValue)
            }
        }
    }()
    updateV3 := ldClient.GetFlagTracker().AddFlagValueChangeListener("v3-feature", user, ldvalue.Bool(false))
    go func() {
        for event := range updateV3 {
            if event.Key == "v3-feature" {
                featureFlags.mu.Lock()
                featureFlags.v3Feature = event.NewValue.BoolValue()
                featureFlags.mu.Unlock()
                notifyClients()
                log.Printf("Flag %q for context %q has changed from %s to %s", event.Key, user.Key(), event.OldValue, event.NewValue)
            }
        }
    }()
}

func notifyClients() {
    sseClients.mu.Lock()
    defer sseClients.mu.Unlock()
    for client := range sseClients.clients {
        select {
        case client <- true:
        default:
        }
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    featureFlags.mu.RLock()
    instantRollback := featureFlags.instantRollback
    newShopFeature := featureFlags.newShopFeature
    v3Feature := featureFlags.v3Feature
    featureFlags.mu.RUnlock()

    fmt.Fprintln(w, `<html><head><title>Hugo the Golden</title><style>
        body { font-family: Arial, sans-serif; background-color: #f8f9fa; margin: 0; padding: 0; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { background-color: #343a40; color: white; padding: 10px 0; text-align: center; }
        .header h1 { margin: 0; font-size: 2.5em; }
        .main { display: flex; justify-content: space-around; flex-wrap: wrap; }
        .item { background-color: white; border: 1px solid #ddd; border-radius: 4px; margin: 20px; padding: 20px; text-align: center; width: calc(33% - 40px); box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
        .item img { max-width: 100%; border-bottom: 1px solid #ddd; margin-bottom: 15px; }
        .item h2 { font-size: 1.5em; margin: 0 0 10px; }
        .item p { font-size: 1em; color: #666; margin: 0 0 15px; }
        .footer { background-color: #343a40; color: white; text-align: center; padding: 10px 0; margin-top: 20px; }
        .map { background-color: white; border: 1px solid #ddd; border-radius: 4px; margin: 20px; padding: 20px; text-align: center; width: calc(50% - 40px); box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
        .map iframe { width: 100%; border: none; border-bottom: 1px solid #ddd; margin-bottom: 15px; }
        .joke-container { display: flex; justify-content: center; align-items: center; flex-direction: column; margin: 20px auto; padding: 20px; background-color: #f8f9fa; border: 1px solid #ddd; border-radius: 4px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); width: 50%; text-align: center; }
        .joke-container h2 { font-size: 1.5em; margin-bottom: 10px; }
        .joke-container p { font-size: 1em; color: #666; margin-bottom: 15px; }
        .joke-container button { padding: 10px 20px; font-size: 1em; color: #fff; background-color: #343a40; border: none; border-radius: 4px; cursor: pointer; }
    </style></head><body>`)
    fmt.Fprintln(w, `<div class="header"><h1>Hugo the Golden</h1></div><div class="container"><div class="main">`)
    if instantRollback {
        fmt.Fprintln(w, `<div class="item"><img src="https://www.akc.org/wp-content/uploads/2017/11/Bernese-Mountain-Dog_Puppy_Bone.jpg" alt="Bernese Mountain Dog"><h2>Bernese Mountain Dog</h2><p>Welcome to Hugo's Pet Shop! Oh no, you're still on the old version with the incorrect dog (although still cute)</p></div>`)
    } else {
        fmt.Fprintln(w, `<div class="item"><img src="https://www.pdsa.org.uk/media/7657/golden-retriever-gallery-1.jpg?" alt="Golden Retriever"><h2>Golden Retriever</h2><p>Welcome to the new homepage of Hugo's Pet Shop! This is the new version of the app with the correct dog - congratulations!</p></div>`)
    }
    if newShopFeature {
        fmt.Fprintln(w, `<div class="map"><iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2522.314574568677!2d-0.6160482840138047!3d51.76069337967569!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0xbfb195792f22448f!2sRoss%20%26%20Friends%20Dog%20Experience!5e0!3m2!1sen!2suk!4v1622549094231!5m2!1sen!2suk" style="height: 450px;" allowfullscreen="" loading="lazy"></iframe><h2>Best Dog Park in the UK!</h2></div>`)
    }
    fmt.Fprintln(w, `</div></div>`)
    if v3Feature {
        fmt.Fprintln(w, `<div class="joke-container">
            <h2>Golden Retriever Jokes</h2>
            <p id="joke">` + jokes[0] + `</p>
            <button id="joke-button">Click for more jokes</button>
        </div>`)
    }
    fmt.Fprintln(w, `<div class="footer"><p>&copy; 2024 Hugo the Golden. All rights reserved.</p></div>
    <script>
    if (!window.evtSource) {
        window.evtSource = new EventSource("/events"); 
        evtSource.onmessage = function(event) { location.reload(); };
    }

    const jokes = [
        "Why did the Golden Retriever sit in the shade? He didn't want to be a hot dog!",
        "What do you get when you cross a Golden Retriever and a telephone? A golden receiver!",
        "What do you call a frozen dog? A pupsicle!"
    ];
    let jokeIndex = 0;
    document.getElementById('joke-button').addEventListener('click', function() {
        jokeIndex = (jokeIndex + 1) % jokes.length;
        document.getElementById('joke').innerText = jokes[jokeIndex];
    });
    </script></body></html>`)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }
    notify := make(chan bool)
    sseClients.mu.Lock()
    sseClients.clients[notify] = struct{}{}
    sseClients.mu.Unlock()

    defer func() {
        sseClients.mu.Lock()
        delete(sseClients.clients, notify)
        sseClients.mu.Unlock()
        close(notify)
    }()

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    for {
        select {
        case <-notify:
            fmt.Fprintf(w, "data: update\n\n")
            flusher.Flush()
        case <-r.Context().Done():
            return
        }
    }
}
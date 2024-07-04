package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    ld "github.com/launchdarkly/go-server-sdk/v7"
    "github.com/launchdarkly/go-server-sdk/v7/ldcomponents"
    "github.com/launchdarkly/go-sdk-common/v3/ldcontext"
)

var ldClient *ld.LDClient

func main() {
    // Set up LaunchDarkly client configuration
    config := ld.Config{
        Events: ldcomponents.SendEvents(),
    }

    // Initialize LaunchDarkly client with SDK key and 10 second timeout
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

    // Check feature flags
    showNewHomepage, err := ldClient.BoolVariation("new-homepage", user, false)
    if err != nil {
        http.Error(w, "Error reading feature flag", http.StatusInternalServerError)
        return
    }

    newShopFeature, err := ldClient.BoolVariation("new-shop-feature", user, false)
    if err != nil {
        http.Error(w, "Error reading feature flag", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, `<html><head><title>Hugo's Pet Shop</title>`)
    fmt.Fprintln(w, `<style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f8f9fa;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #343a40;
            color: white;
            padding: 10px 0;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 2.5em;
        }
        .main {
            display: flex;
            justify-content: space-around;
            flex-wrap: wrap;
        }
        .product {
            background-color: white;
            border: 1px solid #ddd;
            border-radius: 4px;
            margin: 20px;
            padding: 20px;
            text-align: center;
            width: calc(33% - 40px);
            box-shadow: 0 4px 8px rgba(0,0,0,0.1);
        }
        .product img {
            max-width: 100%;
            border-bottom: 1px solid #ddd;
            margin-bottom: 15px;
        }
        .product h2 {
            font-size: 1.5em;
            margin: 0 0 10px;
        }
        .product p {
            font-size: 1em;
            color: #666;
            margin: 0 0 15px;
        }
        .footer {
            background-color: #343a40;
            color: white;
            text-align: center;
            padding: 10px 0;
            margin-top: 20px;
        }
        .map {
            margin-top: 20px;
            text-align: center;
        }
    </style></head><body>`)

    fmt.Fprintln(w, `<div class="header"><h1>Hugo's Pet Shop</h1></div>`)
    fmt.Fprintln(w, `<div class="container"><div class="main">`)

    if showNewHomepage {
        fmt.Fprintln(w, `<div class="product">
            <img src="https://www.pdsa.org.uk/media/7657/golden-retriever-gallery-1.jpg?" alt="Golden Retriever">
            <h2>Golden Retriever</h2>
            <p>Welcome to the new homepage of Hugo's Pet Shop!</p>
        </div>`)
    } else {
        fmt.Fprintln(w, `<div class="product">
            <img src="https://www.akc.org/wp-content/uploads/2017/11/Bernese-Mountain-Dog_Puppy_Bone.jpg" alt="Bernese Mountain Dog">
            <h2>Bernese Mountain Dog</h2>
            <p>Welcome to Hugo's Pet Shop!</p>
        </div>`)
    }

    if newShopFeature {
        fmt.Fprintln(w, `<div class="map">
            <h2>Our Location</h2>
            <iframe src="https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d2522.314574568677!2d-0.6160482840138047!3d51.76069337967569!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0xbfb195792f22448f!2sRoss%20%26%20Friends%20Dog%20Experience!5e0!3m2!1sen!2suk!4v1622549094231!5m2!1sen!2suk" width="600" height="450" style="border:0;" allowfullscreen="" loading="lazy"></iframe>
        </div>`)
    }

    fmt.Fprintln(w, `</div></div>`)
    fmt.Fprintln(w, `<div class="footer"><p>&copy; 2024 Hugo's Pet Shop. All rights reserved.</p></div>`)
    fmt.Fprintln(w, `</body></html>`)
}
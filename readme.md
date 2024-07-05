# Hugo's Pet Shop Feature Flag Application

This is a Go-based web application for Hugo's Pet Shop that leverages LaunchDarkly feature flags to dynamically control the content of the website. The application demonstrates how to implement feature flags, handle real-time updates with Server-Sent Events (SSE), and manage concurrent access to feature flag states.

## Prerequisites

- Go 1.16 or later
- LaunchDarkly account and SDK key

## Getting Started

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/hugos-pet-shop.git
    cd hugos-pet-shop
    ```

2. Install the required dependencies:
    ```sh
    go mod tidy
    ```

### Configuration

1. Replace the LaunchDarkly SDK key in `main.go`:
    ```go
    // Replace the SDK key with your own LaunchDarkly SDK key
    ldClient, err = ld.MakeCustomClient("YOUR_SDK_KEY", config, 10*time.Second)
    ```

2. Replace the user key in `main.go`:
    ```go
    // Replace "example-user-key" with your own user key
    user := ldcontext.New("YOUR_USER_KEY")
    ```

### Running the Application

1. Build and run the application:
    ```sh
    go run main.go
    ```

2. Open your web browser and navigate to `http://localhost:8080`.

## Application Structure

- `main.go`: The main application file that initializes the LaunchDarkly client, sets up feature flag listeners, and defines HTTP handlers.
- `FeatureFlags`: A struct that holds the state of feature flags with concurrent access control.
- `sseClients`: A struct that manages the set of SSE clients.

## Feature Flags

The application uses the following feature flags:

1. `instant-rollback`: Controls the display of the old or new version of the homepage.
2. `new-shop-feature`: Controls the display of a Google Maps iframe showing the best dog park in the UK.
3. `v3-feature`: Controls the display of a section with Golden Retriever jokes.

## HTTP Handlers

- `homeHandler`: Renders the homepage with content controlled by feature flags.
- `sseHandler`: Handles SSE connections to provide real-time updates to clients.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Contact

For questions or suggestions, please contact [James](mailto:james@hendersonmills.co.uk).
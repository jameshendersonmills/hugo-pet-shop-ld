Hugo's Pet Store with Go & LaunchDarkly

Pre-Reqs:
- Make sure you have a LaunchDarkly SDK Key
- Ensure that Go is installed on your local machine - goland.org 
- go mod tidy 
- Ensure you have some flags created in the LD Dashboard - https://app.launchdarkly.com/projects/default/flags


1. Navigate to the clone repo folder Hugos Pet Store and run: 'go mod init hugos-pet-shop'
2. Install the LaunchDarkly Go SDK by running 'go get gopkg.in/launchdarkly/go-server-sdk.v5'
3. The 'main.go' file contains all of the code to run the application locally
4. Replace "YOUR_SDK_KEY" with your actual LaunchDarkly SDK key 
5. flags.json allows you to define the feature flags in Hugo's Pet Store
6. Use 'go run main.go' to run the application locally 



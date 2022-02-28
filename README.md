# MQTT

Lightweight MQTT client broker module written in Golang.
Based on Eclipse Paho MQTT Go client and Mochi MQTT broker libraries

## How to use

To use this module navigate to your project's root folder through a terminal and run the following:
1. Run `go mod init`
2. Run `go get -u github.com/DeAntoLei/go-mqtt`
3. Run `go mod tidy`

Then open all your source files that use the package and add the following to the imports section
1. Add `import github.com/DeAntoLei/go-mqtt`

## Demo

Download the demo folder. This folder contains a client and a broker application.

### Steps to run the apps
#### Broker app
1. Open a terminal
2. Navigate to demo folder
3. Navigate to server folder: `cd broker`
4. Run broker application: `go run main.go`

#### Client app
1. Open a terminal
2. Navigate to demo folder
3. Navigate to client folder: `cd client`
4. Run broker application: `go run main.go`
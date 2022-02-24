package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/DeAntoLei/go-mqtt/broker"
)

var topics = []string{"demo/events", "demo/events2"}

func main() {
	// Signal channel to terminate server application
	close := make(chan os.Signal)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	var mqttconn broker.MqttConnector
	mqttconn.Topics = topics
	mqttconn.AdCount = map[string]int{}
	mqttconn.NewAd = map[string]chan bool{}
	for _, topic := range topics {
		mqttconn.AdCount[topic] = 0
		mqttconn.NewAd[topic] = make(chan bool)
	}

	serverId := "go_mqtt_example_server"

	if err := mqttconn.ConfigureServer(serverId); err != nil {
		fmt.Println("Error while configuring server: ", err)
		os.Exit(1)
	}

	// Start the server
	go mqttconn.StartServer()

	// Publish received advertisements
	go mqttconn.Publish()

	fmt.Println("  Started!  ")

	<-close
	fmt.Println("  Caught Signal  ")

	mqttconn.StopServer()
	fmt.Println("  Finished  ")
}

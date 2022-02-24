package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DeAntoLei/go-mqtt/client"
)

var topics = []string{"demo/events", "demo/events2"}

/*
	Responsible for publishing to all given topics
*/
func publish(mqttconn client.MqttConnector) {
	// Declare a wait group
	var wg sync.WaitGroup
	// Repeat for each one of the client's PubTopics
	for _, topic := range mqttconn.PubTopics {
		wg.Add(1)
		// Run publish subroutine
		go mqttconn.PublishRepeatedly(topic, &wg)
	}
	wg.Wait()
}

func main() {
	// Signal channel to terminate client application
	close := make(chan os.Signal)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	// Connect to broker
	// test mosquitto brokers
	// broker := "tcp://test.mosquitto.org:1883"
	// broker := "ws://test.mosquitto.org:8080"
	var mqttconn client.MqttConnector
	mqttconn.PubTopics = topics
	mqttconn.SubTopics = topics
	mqttconn.TerminatePub = map[string]bool{"demo/events": false, "demo/events2": false}

	// local broker
	brokerUrl := "ws://localhost:8080"
	clientId := "go_mqtt_example"
	mqttconn.ConfigureClient(brokerUrl, clientId)

	// Connect to broker
	go mqttconn.Connect()

	// Disconnect after last routine returns
	defer mqttconn.Disconnect(100)

	// Publish to all given topics
	go publish(mqttconn)

	// Publish to topics and wait until client receives terminate message for all topics
	// go mqttconn.Publish()

	// Closing client application
	<-close
	fmt.Println("  Caught Signal  ")

	fmt.Println("  Finished  ")
}

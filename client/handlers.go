/*
	Based on Eclipse Paho MQTT Go client
	Github: https://github.com/eclipse/paho.mqtt.golang
	Go pkg repo: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang
*/
package client

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
	Simulates a broker response. Holds an incremented
	sequence number, a message type and a timestamp
*/
type BrokerMessage struct {
	SeqNum    int    `json:"seqNum"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

const defaultQoS = 0

/*
	Handles the successfull connection of the client
	to a broker and subscribes the client to its
	SubTopics
*/
func (c *MqttConnector) onConnected(client mqtt.Client) {
	if len(c.SubTopics) > 0 {
		// Subscribe client to its SubTopics, if any
		fmt.Println("> onConnected(): subscribing to all client's sub topics")

		topicFilters := make(map[string]byte)
		for _, topic := range c.SubTopics {
			fmt.Printf("> onConnected(): subscribing to topic %s\n", topic)
			topicFilters[topic] = defaultQoS
		}
		client.SubscribeMultiple(topicFilters, c.messageHandler)
	} else {
		fmt.Println("> onConnected(): no client's sub topics")
	}
}

/*
	Handles the connection loss between the client and
	a broker. Basically acts as a logger
*/
func (c *MqttConnector) onConnectionLost(client mqtt.Client, err error) {
	reader := client.OptionsReader()
	fmt.Printf("> onConnectionLost(): client %s lost connection to the broker: %v\n", reader.ClientID(), err.Error())
}

/*
	Handles the reception
*/
func (c *MqttConnector) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// Decode reply received
	var payload BrokerMessage
	err := json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		fmt.Println("> messageHandler(): error decoding message")
		return
	}
	payloadStr := fmt.Sprintf("seqNum: %d, messageType: %s, timestamp: %s", payload.SeqNum, payload.Type, payload.Timestamp)

	reader := client.OptionsReader()
	fmt.Printf("> messageHandler(): %s received message %q on topic %s\n", reader.ClientID(), payloadStr, msg.Topic())

	// Send terminate signal if respective message type received
	if payload.Type == "terminate" {
		c.TerminatePub[msg.Topic()] = true
	}

	// Alternative without json decoding
	// payloadStr := msg.Payload()
	// fmt.Printf("Client %s received message %s on topic %s\n", reader.ClientID(), payloadStr, msg.Topic())
}

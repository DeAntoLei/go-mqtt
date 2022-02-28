/*
	Based on Eclipse Paho MQTT Go client
	Github: https://github.com/eclipse/paho.mqtt.golang
	Go pkg repo: https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang
*/
package client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/*
	Holds client configuration, broker url, topics
	the clients subscribes and publishes to as well
	as a helper map that shows if the publishing
	proccess for each topic should terminate or not
*/
type MqttConnector struct {
	Client    mqtt.Client
	BrokerUrl string
	PubTopics []string
	SubTopics []string

	TerminatePub map[string]bool
}

/*
	Simulates an event captured by a device. Holds an
	incremented sequence number, a message type, the
	device id and a timestamp
*/
type ClientMessage struct {
	SeqNum    int    `json:"seqNum"`
	Type      string `json:"type"`
	DeviceId  int    `json:"deviceId"`
	Timestamp string `json:"timestamp"`
}

const types int = 4

// Message types sent from client to server
var msgTypes = []string{"onBody", "offBody", "brfChange", "enuEvent"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
	Configures the connector struct and most importantly
	its client field. It receives the client id of
	the new client and the url of the broker where
	the client publishes and subscribes to topics
*/
func (c *MqttConnector) ConfigureClient(brokerUrl, clientId string) {
	// Set broker url
	c.BrokerUrl = brokerUrl

	// Set client options and handlers
	clientOptions := mqtt.NewClientOptions()
	clientOptions.AddBroker(brokerUrl)
	clientOptions.SetClientID(clientId)

	// Set client handlers
	clientOptions.SetDefaultPublishHandler(c.messageHandler)
	clientOptions.OnConnect = c.onConnected
	clientOptions.OnConnectionLost = c.onConnectionLost

	// Set new client
	c.Client = mqtt.NewClient(clientOptions)
}

/*
	Repeatedly attempts to connect the client to
	the specified broker
*/
func (c *MqttConnector) Connect() {
	for {
		fmt.Printf("> connect(): connecting to the broker %v\n", c.BrokerUrl)
		// Break if the client is already connected to the client
		if c.Client.IsConnected() {
			break
		}

		// Attempt to connect
		token := c.Client.Connect()
		if token.Wait() && token.Error() == nil {
			break
		}
		fmt.Printf("> connect(): failed to connect: %v\n", token.Error())
	}

	// Client connected to broker successfully
	fmt.Printf("> connect(): connected to the broker %v\n", c.BrokerUrl)
}

/*
	Disconnect client from broker after provided amount
	of milliseconds
*/
func (c *MqttConnector) Disconnect(milliseconds uint) {
	c.Client.Disconnect(milliseconds)
}

/*
	Repeatedly publishes advertisement to given topic
*/
func (c *MqttConnector) PublishRepeatedly(topic string, wg *sync.WaitGroup) {
	payload := generateAdvertisement()
	// Repeat until a terminate signal has arrived
	for !c.TerminatePub[topic] {
		if c.Client.IsConnected() {
			// Publish advertisement
			token := c.Publish(topic, 0, false, <-payload)
			token.Wait()
			// Wait for one second
			time.Sleep(time.Second)
		}
	}
	wg.Done()
}

/*
	Published single advertisement to given topic
*/
func (c *MqttConnector) Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	return c.Client.Publish(topic, qos, retained, payload)
}

/*
	Generates random advertisement payload
*/
func generateAdvertisement() <-chan []byte {
	ch := make(chan []byte)

	go func() {
		seqNum := 0
		for {
			// Generate JSON payload
			seqNum++
			message := ClientMessage{seqNum, msgTypes[rand.Intn(types)], rand.Intn(1000), time.Now().Format("02/01/2006 15:04:05")}
			payload, err := json.Marshal(message)
			if err != nil {
				fmt.Printf("> publish(): error while preparing message: %s\n", err)
				break
			}

			// Sending payload through the channel
			ch <- payload
		}
	}()

	return ch
}

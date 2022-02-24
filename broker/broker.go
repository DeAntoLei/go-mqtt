/*
	Based on Mochi MQTT server
	Github: https://github.com/mochi-co/mqtt
	Go pkg repo: https://pkg.go.dev/github.com/mochi-co/mqtt/server
*/
package broker

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/mochi-co/mqtt/server"
	"github.com/mochi-co/mqtt/server/listeners"
)

/*
	Holds server configuration, topics that will be
	served as well as helper maps that show whether
	a new advertisement was received and keep count
	of the total number of advertisements received
	on a certain topic, respectively
*/
type MqttConnector struct {
	Server *mqtt.Server
	Topics []string

	NewAd   map[string]chan bool
	AdCount map[string]int
}

/*
	Simulates a broker response. Holds an incremented
	sequence number, a message type and a timestamp
*/
type BrokerMessage struct {
	SeqNum    int    `json:"seqNum"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

const types int = 3

// Message types sent from server to client
var msgTypes = []string{"erc", "drc", "res"}

func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
	Configures the connector struct and most importantly
	its server field. It receives the new server id
*/
func (c *MqttConnector) ConfigureServer(serverId string) error {
	fmt.Println("> configureServer(): setting up mqtt server... WS")

	// Set up new server
	c.Server = mqtt.New()
	ws := listeners.NewWebsocket(serverId, ":8080")
	if err := c.Server.AddListener(ws, nil); err != nil {
		return err
	}

	// Set server handlers
	c.Server.Events.OnConnect = c.onConnected
	c.Server.Events.OnDisconnect = c.onConnectionLost
	c.Server.Events.OnMessage = c.messageHandler

	return nil
}

/*
	Starts server application
*/
func (c *MqttConnector) StartServer() {
	err := c.Server.Serve()
	if err != nil {
		fmt.Println("> startServer(): error while starting server =>", err)
	}
}

/*
	Stops server application
*/
func (c *MqttConnector) StopServer() {
	c.Server.Close()
}

/*
	Concurrently publishes random event messages
	to all topics so long as a new advertisement
	has been received
*/
func (c *MqttConnector) Publish() {
	// Repeat for each one of the topics
	for _, topic := range c.Topics {
		go func(topic string) {
			for {
				<-c.NewAd[topic]
				var msgType string
				if c.AdCount[topic] < 15 {
					msgType = msgTypes[rand.Intn(types)]
				} else {
					msgType = "terminate"
				}
				// responseTxt := fmt.Sprintf(`{ "seqNum": "%d", "type": "%s", "timestamp": "%s" }`, msgCount[pkx.TopicName], msgType, time.Now().Format("02/01/2006 15:04:05"))
				response := BrokerMessage{c.AdCount[topic], msgType, time.Now().Format("02/01/2006 15:04:05")}
				responseTxt, err := json.Marshal(response)
				if err != nil {
					fmt.Printf("> publish(): error while preparing message => %s\n", err)
					break
				}

				err = c.Server.Publish(topic, []byte(responseTxt), false)
				if err != nil {
					fmt.Println(err)
					break
				}
				fmt.Printf("> publish(): issued message %s to %s\n", responseTxt, topic)
				c.NewAd[topic] = make(chan bool)
			}
		}(topic)
	}
}

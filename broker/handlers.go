/*
	Based on Mochi MQTT server
	Github: https://github.com/mochi-co/mqtt
	Go pkg repo: https://pkg.go.dev/github.com/mochi-co/mqtt/server
*/
package broker

import (
	"fmt"

	"github.com/mochi-co/mqtt/server/events"
)

func (c *MqttConnector) onConnected(client events.Client, pk events.Packet) {
	fmt.Printf("> onConnected(): client %s connected => %+v\n", client.ID, pk)
}

func (c *MqttConnector) onConnectionLost(client events.Client, err error) {
	fmt.Printf("> onConnectionLost(): client %s disconnected => %v\n", client.ID, err)
}

func (c *MqttConnector) messageHandler(client events.Client, pk events.Packet) (events.Packet, error) {
	packet := pk
	fmt.Printf("> messageHandler(): received message %s from client %s\n", string(packet.Payload), client.ID)

	c.AdCount[packet.TopicName]++
	c.NewAd[packet.TopicName] <- true

	return packet, nil
}

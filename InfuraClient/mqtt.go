//Original Author: bdjukic

package main

import (
	"fmt"
	// "io/ioutil"
	// "net/http"
	"os"
	"os/signal"
	// "strings"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

func main() {
	// Set up channel on which to send signal notifications.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "m13.cloudmqtt.com:12423",
		ClientID: []byte("Infura-Client"),
		UserName: []byte("skupntaq"),
		Password: []byte("yRS1mpvJp8su"),
	})
	if err != nil {
		panic(err)
	}

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("p2penergy/photon/events"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(("New Transaction received"))
					fmt.Println(string(message));
				},
			},
			&client.SubReq{
				TopicFilter: []byte("power1"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(("New power1 received"))
					fmt.Println(string(message));
				},
			},
			&client.SubReq{
				TopicFilter: []byte("power2"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(("New power2 received"))
					fmt.Println(string(message));
				},
			},
			&client.SubReq{
				TopicFilter: []byte("power3"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					fmt.Println(("New power3 received"))
					fmt.Println(string(message));
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// Wait for receiving a signal.
	<-sigc

	// Disconnect the Network Connection.
	if err := cli.Disconnect(); err != nil {
		panic(err)
	}
}

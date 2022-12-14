package client

import (
	"fmt"
	"time"
)

// Heartbeats section

func startHeartbeat(client Client, interval int32) {
	// create a ticker that ticks every 5 milliseconds
	client.heartbeatTicker = time.NewTicker(5000 * time.Millisecond)

	sendIdentify(client)


	for {
		if !client.isConnected {
			// stop ticker
			stopHeartbeat(client)
			break
		}

		select {
		case <-client.heartbeatTicker.C:
			// send heartbeat
			fmt.Println("SENDING HEARTBEAT")
			client.writeData(map[string]interface{}{"op": 1, "d": client.sequence})
		}
	}
}

func stopHeartbeat(client Client) {
	client.heartbeatTicker.Stop()
}

// Identify section

func sendIdentify(client Client) {
	client.writeData(map[string]interface{}{"op": 2, "d": map[string]interface{}{"token": client.Token, "properties": map[string]interface{}{"$os": "windows", "$browser": "gophercord", "$device": "gophercord"}, "compress": false, "large_threshold": 250, "shard": []int{0, 1}}})
}

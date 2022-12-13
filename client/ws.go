package client

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/davipatricio/gophercord/payloads"
	"nhooyr.io/websocket"
)

// create a context with a timeout of 1 minute
var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)

func (c *Client) Connect() (err error) {
	// TODO: handle error
	c.ws, _, err = websocket.Dial(ctx, "wss://gateway.discord.gg/?v=11&encoding=json&compress=zlib-stream", nil)
	if err != nil {
		fmt.Println("failed to connect to discord websocket:", err)
		return err
	}

	c.isConnected = true

	go c.listenToMessages(ctx, cancel)
	return
}

func (c *Client) Disconnect() error {
	return c.ws.Close(1000, "")
}

func (c *Client) writeData(data map[string]interface{}) error {
	// convert the data map to json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// write the data to the websocket
	err = c.ws.Write(ctx, websocket.MessageText, jsonData)
	if err != nil {
		return err
	}

	return nil
}

// start listening to messages from discord
// blocks the current thread
func (c *Client) listenToMessages(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		messageType, message, err := c.ws.Read(ctx)
		if err != nil {
			// check if the connection was closed. if the status code is not -1, it means that the connection was closed
			if websocket.CloseStatus(err) != -1 {
				fmt.Println("connection closed")
				c.isConnected = false
				break
			}

			continue
		}

		go c.handleMessage(messageType, message)
	}
}

func (c *Client) handleMessage(messageType websocket.MessageType, data []byte) (err error) {
	var jsonData payloads.BasePayload
	var dataReader *bytes.Buffer

	// check if discord sent a compressed message
	if messageType == websocket.MessageBinary {
		if dataReader, err = decompressBytes(data); err != nil {
			return err
		}
	} else {
		dataReader = bytes.NewBuffer(data)
	}

	decoder := json.NewDecoder(dataReader)
	err = decoder.Decode(&jsonData)
	if err != nil {
		return err
	}

	fmt.Println(jsonData)

	switch jsonData.Op {
	case 10:
		data := payloads.CreateHelloPayload(jsonData)
		// start sending heartbeats
		go startHeartbeat(*c, data.HeartbeatInterval)
	}

	return nil
}

func decompressBytes(data []byte) (out *bytes.Buffer, err error) {
	if z, err := zlib.NewReader(bytes.NewReader(data)); err == nil {
		defer z.Close()

		out := new(bytes.Buffer)
		out.ReadFrom(z)
	}

	return
}

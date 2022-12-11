package client

import (
	"bytes"
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/davipatricio/gophercord/payloads"
	"nhooyr.io/websocket"
)

// TODO: fix compression
var zlibReader, _ = zlib.NewReader(bytes.NewReader(make([]byte, 0)))

func (c *Client) Connect() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	// TODO: handle error
	c.ws, _, err = websocket.Dial(ctx, "wss://gateway.discord.gg/?v=11&encoding=json&compress=zlib-stream", nil)

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
	err = c.ws.Write(context.Background(), websocket.MessageText, jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listenToMessages(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	for {
		messageType, message, err := c.ws.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				fmt.Println("connection closed")
				return
			}

			fmt.Println(err)
			continue
		}

		go c.handleMessage(messageType, message)
	}
}

func (c *Client) handleMessage(messageType websocket.MessageType, data []byte) error {
	var jsonData payloads.BasePayload

	// check the message type
	if messageType == websocket.MessageBinary {
		fmt.Println("compressed message detected!")
		// use zlib to decompress the data

		// read the decompressed data
		decompressedData := bytes.NewBuffer(make([]byte, 0))
		_, err := io.Copy(decompressedData, zlibReader)
		if err != nil {
			fmt.Println(err)
			// check if err is unexpected EOF
			if err != io.ErrUnexpectedEOF {
				return err
			}
		}

		// json decode the decompressed data
		err = json.Unmarshal(decompressedData.Bytes(), &jsonData)
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		// parse the data array into a json object
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}
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

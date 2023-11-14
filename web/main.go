package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := "wss://ws.blockchain.info/inv"
	log.Printf("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	// Start a goroutine for receiving messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("Received message: %s\n", message)
		}
	}()

	// Subscribe to new blocks
	subscribeMsg := `{"op":"blocks_sub"}`
	err = c.WriteMessage(websocket.TextMessage, []byte(subscribeMsg))
	if err != nil {
		log.Fatal("write:", err)
	}

	// Start a goroutine for the heartbeat
	go heartbeat(c, done)

	// Main loop to handle interrupt signal
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

// heartbeat sends a ping message to the server periodically
func heartbeat(c *websocket.Conn, done chan struct{}) {
	heartbeatInterval := 30 * time.Second
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			err := c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Println("write ping:", err)
				close(done)
				return
			}
			fmt.Println("Heartbeat sent")
		}
	}
}

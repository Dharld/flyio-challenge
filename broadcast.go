package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Node struct {
	*maelstrom.Node
	messages []int // Store for broadcast messages
	mutex sync.Mutex
}

func NewNode() *Node {
	return &Node{
		Node: maelstrom.NewNode(),
		messages: make([]int, 0),
	}
}

func (n *Node) HandleBroadcast(msg maelstrom.Message) error {
	// Decode the message
	var body map[string]any

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Add the message to the list
	num := body["message"].(int)

	// Protect from concurrent write
	n.mutex.Lock()
	n.messages = append(n.messages, num)
	n.mutex.Unlock()

	// Send acknowledgement
	response := map[string]any{
		"type": "broadcast_ok",
	}

	return n.Reply(msg, response)
}

func broadcast() {
	n := NewNode()

	// Listen to the broadcast event
	n.Handle("broadcast", n.HandleBroadcast)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
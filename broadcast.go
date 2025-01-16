package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Node struct {
	*maelstrom.Node
	messages []int // Store for broadcast messages
	mutex sync.Mutex
	neighboursID []string // Store the ids of neighbours
}

func NewNode() *Node {
	return &Node{
		Node: maelstrom.NewNode(),
		messages: make([]int, 0),
		neighboursID: make([]string, 0),
	}
}

func (n *Node) HandleBroadcast(msg maelstrom.Message) error {
    // Decode the message
    var body map[string]any
    if err := json.Unmarshal(msg.Body, &body); err != nil {
        return err
    }

    // Add the message to the list
    // Convert float64 to int
    messageFloat, ok := body["message"].(float64)
    if !ok {
        return fmt.Errorf("message is not a number")
    }
    num := int(messageFloat)  // Convert float64 to int

    // Protect from concurrent write
    n.mutex.Lock()
    defer n.mutex.Unlock()

    n.messages = append(n.messages, num)

    // Send acknowledgement
    response := map[string]any{
        "type": "broadcast_ok",
    }

    return n.Reply(msg, response)
}

func (n *Node) HandleRead(msg maelstrom.Message) error {

	n.mutex.Lock()
	defer n.mutex.Unlock()

	messages := make([]int, len(n.messages))
	copy(messages, n.messages)

	// Send the list of messages
	response := map[string]any{
		"type": "read_ok",
		"messages": messages,
	}

	return n.Reply(msg, response)
}

func (n *Node) HandleTopology(msg maelstrom.Message) error {
    var body map[string]any
    if err := json.Unmarshal(msg.Body, &body); err != nil {
        return err
    }

    // Get topology from body
    topology, ok := body["topology"].(map[string]any)
    if !ok {
        return fmt.Errorf("invalid topology format")
    }

    // Get our node's neighbors
    if neighbours, ok := topology[n.ID()].([]any); ok {
        // Convert []any to []string
        n.neighboursID = make([]string, len(neighbours))
        for i, v := range neighbours {
            if s, ok := v.(string); ok {
                n.neighboursID[i] = s
            }
        }
    }

    // Send acknowledgement
    response := map[string]any{
        "type": "topology_ok",
    }

    return n.Reply(msg, response)
}


func Broadcast() {
	n := NewNode()

	// Listen to the broadcast event
	n.Handle("broadcast", n.HandleBroadcast)
	n.Handle("read", n.HandleRead)
	n.Handle("topology", n.HandleTopology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
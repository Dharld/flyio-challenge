package main

import (
	"encoding/json"
	"log"
	"maelstrom-echo/internal/service/id"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)


func UniqueIds() {
	n := maelstrom.NewNode()

	// Listen to the generate event
	n.Handle("generate", func (msg maelstrom.Message) error {
		println("Generating a new ID")
		// Generate a new ID
		generator := id.NewGenerator()
		newID := generator.Generate()

		// Unmarshall the message
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Update the message type to return back
		body["type"] = "generate_ok"
		body["id"] = newID

		// Reply to the server with the newID
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
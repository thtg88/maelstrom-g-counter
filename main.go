package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	AddType   = "add"
	AddOkType = "add_ok"

	ReadType   = "read"
	ReadOkType = "read_ok"

	ReplicateType   = "replicate"
	ReplicateOkType = "replicate_ok"
)

// RPC: `add`
// Your node should accept add requests and increment the value of a single global counter.
// Your node will receive a request message body that looks like this:
// {
//   "type": "add",
//   "delta": 123
// }
// and it will need to return an "add_ok" acknowledgement message:
// {
//   "type": "add_ok"
// }

type AddBody struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type AddOkBody struct {
	Type string `json:"type"`
}

// RPC: `read`
// Your node should accept read requests and return the current value of the global counter.
// Remember that the counter service is only sequentially consistent.
// Your node will receive a request message body that looks like this:
// {
//   "type": "read"
// }
// and it will need to return a "read_ok" message with the current value:
// {
//   "type": "read_ok",
//   "value": 1234
// }

type ReadBody struct {
	Type string `json:"type"`
}

type ReadOkBody struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

// We add a `replicate` message so that we can propagate our value to other nodes
// We don't have a topology given to us at the beginning of the test from Maelstrom
// So we assume we replicate to all nodes in the topology 🥵
// {
// 	"type": "replicate",
// 	"value": 12345
// }
// and it will need to return a "replicate_ok" acknowledgment message:
// {
// 	"type": "replicate_ok"
// }

type ReplicateBody struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type ReplicateOkBody struct {
	Type string `json:"type"`
}

func main() {
	var mutex sync.RWMutex

	counters := make(map[string]int)

	node := maelstrom.NewNode()

	// Continuosly replicate value to other nodes
	go func() {
		for {
			mutex.RLock()
			value := counters[node.ID()]
			mutex.RUnlock()

			for _, nodeID := range node.NodeIDs() {
				// Don't send to ourselves lol
				if nodeID == node.ID() {
					continue
				}

				// No need to wait for a reply from the node, as it will _eventually_ send
				node.Send(nodeID, ReplicateBody{
					Type:  ReplicateType,
					Value: value,
				})
			}

			time.Sleep(1 * time.Second)
		}
	}()

	node.Handle(AddType, func(msg maelstrom.Message) error {
		var body AddBody

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		mutex.Lock()
		defer mutex.Unlock()

		// Init node's value if first time accessing it
		if _, ok := counters[node.ID()]; !ok {
			counters[node.ID()] = 0
		}

		counters[node.ID()] = counters[node.ID()] + body.Delta

		return node.Reply(msg, AddOkBody{Type: AddOkType})
	})

	node.Handle(ReadType, func(msg maelstrom.Message) error {
		total := 0

		// Value of the counter is the sum of all the nodes values
		mutex.RLock()
		for _, value := range counters {
			total = total + value
		}
		mutex.RUnlock()

		return node.Reply(msg, ReadOkBody{
			Type:  ReadOkType,
			Value: total,
		})
	})

	node.Handle(ReplicateType, func(msg maelstrom.Message) error {
		var body ReplicateBody

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// We may receive an old counter value (e.g. if we were partitioned)
		// So make sure to store the max of the 2
		mutex.Lock()
		counters[msg.Src] = max(counters[msg.Src], body.Value)
		mutex.Unlock()

		return node.Reply(msg, ReplicateOkBody{Type: ReplicateOkType})
	})

	node.Handle(ReplicateOkType, func(msg maelstrom.Message) error {
		return nil
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}

// max returns the larger of x or y.
func max(x int, y int) int {
	if x < y {
		return y
	}

	return x
}

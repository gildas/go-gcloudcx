package purecloud

import (
	"github.com/pkg/errors"
)

// Quieue defines a PureCloud Queue
type Queue struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (client *Client) FindQueueByName(name string) (*Queue, error) {
	response := struct {
		Entities []*Queue `json:"entities"`
	}{}
	// /routing/queues?name=$name
	err := client.Get("/routing/queues?pageSize=200", &response)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, queue := range response.Entities {
		if queue.Name == name {
			return queue, nil
		}
	}
	return nil, errors.Errorf("Queue not found: %s", name)
}

// GetID gets the identifier of this
//   implements Identifiable
func (queue Queue) GetID() string {
	return queue.ID
}

func (queue Queue) String() string {
	if len(queue.Name) > 0 {
		return queue.Name
	}
	return queue.ID
}

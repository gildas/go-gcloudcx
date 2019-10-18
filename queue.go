package purecloud

// Quieue defines a PureCloud Queue
type Queue struct {
	ID string   `json:"id"`
	Name string `json:"name"`
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
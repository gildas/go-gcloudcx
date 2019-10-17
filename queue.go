package purecloud

// Quieue defines a PureCloud Queue
type Queue struct {
	ID string `json:"id"`
}

// GetID gets the identifier of this
//   implements Identifiable
func (queue Queue) GetID() string {
	return queue.ID
}
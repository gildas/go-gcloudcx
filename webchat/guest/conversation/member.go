package conversation

import (
	"fmt"
)

// GetMembers fetches the Members of this Conversation
func (conversation *Conversation) GetMembers() ([]*Member, error) {
	var members []*Member

	if err := conversation.Client.Get(fmt.Sprintf("webchat/guest/conversations/%s/members", conversation.ID), nil, &members); err != nil {
		return nil, err
	}
	return members, nil
}
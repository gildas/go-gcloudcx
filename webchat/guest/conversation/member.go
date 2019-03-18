package conversation

import (
	"fmt"
)

// GetMembers fetches the Members of this Conversation
func (conv *Conversation) GetMembers() ([]*Member, error) {
	var members []*Member

	if err := conv.Client.Get(fmt.Sprintf("webchat/guest/conversations/%s/members", conv.ID), nil, &members); err != nil {
		return nil, err
	}
	return members, nil
}
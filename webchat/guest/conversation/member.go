package conversation

import (
	"fmt"

	"github.com/gildas/go-core"
)

type getMembersResponse struct {
	Members     []Member `json:"entities"`
	PageSize    uint32   `json:"pageSize,omitempty"`
	PageNumber  uint32   `json:"pageNumber,omitempty"`
	PageCount   uint32   `json:"pageCount,omitempty"`
	Total       uint32   `json:"total,omitempty"`
	FirstURI    string   `json:"firstUri,omitempty"`
	SelfURI     string   `json:"selfUri,omitempty"`
	NextURI     string   `json:"nextUri,omitempty"`
	PreviousURI string   `json:"previousUri,omitempty"`
	LastURI     string   `json:"lastUri,omitempty"`
}
// GetMembers fetches the Members of this Conversation
func (conversation *Conversation) GetMembers() ([]Member, error) {
	response := &getMembersResponse{}
	err := conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members", conversation.ID),
		&core.RequestOptions{
			Authorization: "bearer " + conversation.JWT,
		},
		&response,
	)
	if err != nil {
		return nil, err
	}
	conversation.Logger.Record("scope", "getmembers").Debugf("Response: %+v", response)
	return response.Members, nil
}

// GetMember fetches the given member of this Conversation (caches the member)
func (conversation *Conversation) GetMember(id string) (*Member, error) {
	if member, ok := conversation.Members[id]; ok {
		return member, nil
	}
	member := &Member{}
	err := conversation.Client.SendRequest(
		fmt.Sprintf("/webchat/guest/conversations/%s/members/%s", conversation.ID, id),
		&core.RequestOptions{
			Authorization: "bearer " + conversation.JWT,
		},
		&member,
	)
	if err != nil {
		return nil, err
	}
	conversation.Logger.Record("scope", "getmembers").Debugf("Response: %+v", member)
	conversation.Members[member.ID] = member
	return member, nil

}
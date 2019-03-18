package conversation

import (
	"github.com/gildas/go-purecloud"
	"fmt"
)

type getMembersResponse struct {
	Members     []Member `json:"entities"`
	PageSize    uint32   `json:"pageSize,omitifempty"`
	PageNumber  uint32   `json:"pageNumber,omitifempty"`
	PageCount   uint32   `json:"pageCount,omitifempty"`
	Total       uint32   `json:"total,omitifempty"`
	FirstURI    string   `json:"firstUri,omitifempty"`
	SelfURI     string   `json:"selfUri,omitifempty"`
	NextURI     string   `json:"nextUri,omitifempty"`
	PreviousURI string   `json:"previousUri,omitifempty"`
	LastURI     string   `json:"lastUri,omitifempty"`
}
// GetMembers fetches the Members of this Conversation
func (conversation *Conversation) GetMembers() ([]Member, error) {
	response := &getMembersResponse{}
	if err := conversation.Client.Get(fmt.Sprintf("webchat/guest/conversations/%s/members", conversation.ID), nil, &response, purecloud.RequestOptions{ Authorization: "bearer " + conversation.JWT}); err != nil {
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
	if err := conversation.Client.Get(fmt.Sprintf("webchat/guest/conversations/%s/members/%s", conversation.ID, id), nil, &member, purecloud.RequestOptions{ Authorization: "bearer " + conversation.JWT}); err != nil {
		return nil, err
	}
	conversation.Logger.Record("scope", "getmembers").Debugf("Response: %+v", member)
	conversation.Members[member.ID] = member
	return member, nil

}
package gcloudcx

import (
	"context"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ResponseManagementResponse is the interface for the Response Management Response
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementResponse struct {
	ID            uuid.UUID                        `json:"id"`
	Name          string                           `json:"name"`
	DateCreated   time.Time                        `json:"dateCreated,omitempty"`
	CreatedBy     *DomainEntityRef                 `json:"createdBy,omitempty"`
	Libraries     []DomainEntityRef                `json:"libraries,omitempty"`
	Texts         []ResponseManagementContent      `json:"texts,omitempty"`
	Substitutions []ResponseManagementSubstitution `json:"substitutions,omitempty"`
	Version       int                              `json:"version"`
	SelfURI       URI                              `json:"selfUri,omitempty"`
}

// ResponseManagementContent is the interface for the Response Management Content
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementContent struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

type ResponseManagementQuery struct {
	Query   string                          `json:"queryPhrase,omitempty"`
	PageSize int                            `json:"pageSize,omitempty"`
	Filters []ResponseManagementQueryFilter `json:"filters,omitempty"`
}

type ResponseManagementQueryFilter struct {
	Name     string   `json:"name"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// ResponseManagementSubstitution is the interface for the Response Management Substitution
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementSubstitution struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Default     string    `json:"defaultValue"`
}

// Fetch fetches this from the given Client
//
//  implements Fetchable
func (response *ResponseManagementResponse) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	// TODO: Allow to filter per Library ID
	id, name, uri, _ := client.ParseParameters(ctx, response, parameters...)
	if len(uri) > 0 {
		return client.Get(ctx, uri, response)
	}
	if id != uuid.Nil {
		return client.Get(ctx, NewURI("/responsemanagement/responses/%s", id), response)
	}
	if len(name) > 0 {
		nameLowercase := strings.ToLower(name)
		results := struct {
			Results   struct {
				Responses []ResponseManagementResponse `json:"entities"`
				paginatedEntities
			} `json:"results"`
		}{}
		err := client.Post(
			ctx,
			NewURI("/responsemanagement/responses/query"),
			ResponseManagementQuery{
				Filters: []ResponseManagementQueryFilter{
					{Name: "name", Operator: "EQUALS", Values: []string{name}},
				},
			},
			&results,
		)
		if err != nil {
			return err
		}
		for _, entity := range results.Results.Responses {
			if strings.Compare(strings.ToLower(entity.Name), nameLowercase) == 0 {
				*response = entity
				return nil
			}
		}
		return errors.NotFound.With("ResponseManagementResponse", name)
	}
	return errors.NotFound.With("ResponseManagementResponse")
}

// GetID gets the identifier of this
//
//   implements Identifiable
func (response *ResponseManagementResponse) GetID() uuid.UUID {
	return response.ID
}

// String gets a string version
//
//   implements the fmt.Stringer interface
func (response ResponseManagementResponse) String() string {
	if len(response.Name) > 0 {
		return response.Name
	}
	return response.ID.String()
}

package gcloudcx

import (
	"context"
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
	Query    string                          `json:"queryPhrase,omitempty"`
	PageSize int                             `json:"pageSize,omitempty"`
	Filters  []ResponseManagementQueryFilter `json:"filters,omitempty"`
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
	ID          string `json:"id"`
	Description string `json:"description"`
	Default     string `json:"defaultValue"`
}

// Initialize initializes the object
//
// implements Initializable
func (response *ResponseManagementResponse) Initialize(parameters ...interface{}) {
}

// GetID gets the identifier of this
//
// implements Identifiable
func (response ResponseManagementResponse) GetID() uuid.UUID {
	return response.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (response ResponseManagementResponse) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/responsemanagement/responses/%s", ids[0])
	}
	if response.ID != uuid.Nil {
		return NewURI("/api/v2/responsemanagement/responses/%s", response.ID)
	}
	return URI("/api/v2/responsemanagement/responses/")
}

// String gets a string version
//
// implements the fmt.Stringer interface
func (response ResponseManagementResponse) String() string {
	if len(response.Name) > 0 {
		return response.Name
	}
	return response.ID.String()
}

func (response ResponseManagementResponse) FetchByFilters(context context.Context, client *Client, filters ...ResponseManagementQueryFilter) (*ResponseManagementResponse, error) {
	results := struct {
		Results struct {
			Responses []ResponseManagementResponse `json:"entities"`
			paginatedEntities
		} `json:"results"`
	}{}
	err := client.Post(
		context,
		NewURI("/responsemanagement/responses/query"),
		ResponseManagementQuery{
			Filters: filters,
		},
		&results,
	)
	if err != nil {
		return nil, err
	}
	if len(results.Results.Responses) == 0 {
		return nil, errors.NotFound.With("ResponseManagementResponse")
	}
	/*
		for _, entity := range results.Results.Responses {
			if strings.Compare(strings.ToLower(entity.Name), nameLowercase) == 0 {
				*response = entity
				return nil
			}
		}
	*/
	return &results.Results.Responses[0], nil
}

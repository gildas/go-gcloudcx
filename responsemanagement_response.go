package gcloudcx

import (
	"context"
	"strings"
	"text/template"
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

// ResponseManagementContent represent a Response Management Content
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

// ApplySubstitutions applies the substitutions to the response text that matches the given content type
func (response ResponseManagementResponse) ApplySubstitutions(contentType string, substitutions map[string]string) (string, error) {
	var text string
	for _, content := range response.Texts {
		if strings.Compare(strings.ToLower(content.ContentType), strings.ToLower(contentType)) == 0 {
			text = content.Content
			break
		}
	}
	if len(text) == 0 {
		return "", errors.NotFound.With("text of type ", contentType)
	}

	if len(substitutions) == 0 {
		return text, nil
	}

	// Apply defaults from the response to the given substitutions
	for _, substitution := range response.Substitutions {
		if s, ok := substitutions[substitution.ID]; ok && len(s) == 0 {
			substitutions[substitution.ID] = substitution.Default
		}
	}

	// change the placeholders to Go Template
	for id, _ := range substitutions {
		text = strings.ReplaceAll(text, "{{"+id+"}}", "{{."+id+"}}")
	}

	// apply the substitutions
	tpl, err := template.New("response").Parse(text)
	if err != nil {
		return "", errors.WrapErrors(errors.ArgumentInvalid.With("text", "..."), err)
	}
	var expanded strings.Builder
	err = tpl.Execute(&expanded, substitutions)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return expanded.String(), nil
}

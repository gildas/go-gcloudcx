package gcloudcx

import (
	"context"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// ResponseManagementResponse is the interface for the Response Management Response
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementResponse struct {
	ID                uuid.UUID                        `json:"id"`
	Name              string                           `json:"name"`
	Type              string                           `json:"ResponseType"`
	DateCreated       time.Time                        `json:"dateCreated,omitempty"`
	CreatedBy         *DomainEntityRef                 `json:"createdBy,omitempty"`
	Libraries         []DomainEntityRef                `json:"libraries,omitempty"`
	Texts             []ResponseManagementContent      `json:"texts,omitempty"`
	Substitutions     []ResponseManagementSubstitution `json:"substitutions,omitempty"`
	TemplateType      string                           `json:"-"`
	TemplateName      string                           `json:"-"`
	TemplateNamespace string                           `json:"-"`
	TemplateLanguage  string                           `json:"-"`
	Version           int                              `json:"version"`
	SelfURI           URI                              `json:"selfUri,omitempty"`
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
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			response.ID = parameter
		}
	}
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

// GetType gets the type of this
//
// implements core.TypeCarrier
func (response ResponseManagementResponse) GetType() string {
	return response.Type
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
func (response ResponseManagementResponse) ApplySubstitutions(context context.Context, contentType string, substitutions map[string]string) (string, error) {
	log := logger.Must(logger.FromContext(context, logger.Create("gcloudcx", "nil"))).Child("response", "applysubstitutions", "response", response.ID)
	// Logging is done at the TRACE level since the response and/or the text could contain sensitive information

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

	if len(substitutions) == 0 && len(response.Substitutions) == 0 {
		return text, nil
	}
	if substitutions == nil {
		substitutions = make(map[string]string)
	}

	// Apply defaults from the response to the given substitutions
	log.Record("substitutions", substitutions).Tracef("Substitutions")
	for _, substitution := range response.Substitutions {
		if s, ok := substitutions[substitution.ID]; ok {
			if len(s) == 0 {
				substitutions[substitution.ID] = substitution.Default
			}
		} else {
			substitutions[substitution.ID] = substitution.Default
		}
	}
	log.Record("substitutions", substitutions).Tracef("Substitutions with defaults")

	// change the Genesys Cloud placeholders to Go Template placeholders
	for id := range substitutions {
		text = strings.ReplaceAll(text, "{{"+id+"}}", "{{."+id+"}}")
	}
	log.Tracef("Replaced gcloud placeholders with Go Template placeholders: %s", text)

	// apply the substitutions
	tpl, err := template.New("response").Funcs(sprig.TxtFuncMap()).Parse(text)
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

// MarshalJSON marshals into JSON
//
// implements json.Marshaler
func (response ResponseManagementResponse) MarshalJSON() ([]byte, error) {
	type surrogate ResponseManagementResponse
	type MessagingTemplate struct {
		WhatsApp struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
			Language  string `json:"language"`
		} `json:"whatsApp"`
	}

	data, err := json.Marshal(struct {
		surrogate
		MessagingTemplate *MessagingTemplate `json:"messagingTemplate,omitempty"`
	}{
		surrogate: surrogate(response),
		MessagingTemplate: func() *MessagingTemplate {
			switch response.TemplateType {
			case "whatsApp":
				template := MessagingTemplate{}
				template.WhatsApp.Name = response.TemplateName
				template.WhatsApp.Namespace = response.TemplateNamespace
				template.WhatsApp.Language = response.TemplateLanguage
				return &template
			default:
				return nil
			}
		}(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals JSON
//
// implements json.Unmarshaler
func (response *ResponseManagementResponse) UnmarshalJSON(payload []byte) (err error) {
	type surrogate ResponseManagementResponse
	var inner struct {
		surrogate
		MessagingTemplate *struct {
			WhatsApp *struct {
				Name      string `json:"name"`
				Namespace string `json:"namespace"`
				Language  string `json:"language"`
			} `json:"whatsApp"`
		}
	}

	if err := json.Unmarshal(payload, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*response = ResponseManagementResponse(inner.surrogate)
	if inner.MessagingTemplate != nil {
		if inner.MessagingTemplate.WhatsApp != nil {
			response.TemplateType = "whatsApp"
			response.TemplateName = inner.MessagingTemplate.WhatsApp.Name
			response.TemplateNamespace = inner.MessagingTemplate.WhatsApp.Namespace
			response.TemplateLanguage = inner.MessagingTemplate.WhatsApp.Language
		}
	}
	return
}

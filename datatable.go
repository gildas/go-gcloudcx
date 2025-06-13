package gcloudcx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-request"
	"github.com/google/uuid"
)

// DataTable describes Data tables used by Architect
//
// See: https://developer.genesys.cloud/routing/architect/flows
type DataTable struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Division    *Division        `json:"division,omitempty"`
	Schema      *DataTableSchema `json:"schema,omitempty"`
	client      *Client          `json:"-"`
	logger      *logger.Logger   `json:"-"`
}

// DataTableRow describes a row in a Data Table
type DataTableRow map[string]interface{}

type DataTableSchema struct {
	Schema               string                             `json:"$schema"`
	Type                 string                             `json:"type"`
	Title                string                             `json:"title"`
	Description          string                             `json:"description,omitempty"`
	DataTableID          uuid.UUID                          `json:"datatableId"`
	Required             []string                           `json:"required"`
	Properties           map[string]DataTableSchemaProperty `json:"properties"`
	AdditionalProperties interface{}                        `json:"additionalProperties"`
}

type DataTableSchemaProperty struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	ID        string `json:"$id"`
	Order     string `json:"displayOrder"`
	MinLength int    `json:"minLength"`
	MaxLength int    `json:"maxLength"`
}

// Initialize initializes the object
//
// implements Initializable
func (table *DataTable) Initialize(parameters ...interface{}) {
	for _, raw := range parameters {
		switch parameter := raw.(type) {
		case uuid.UUID:
			table.ID = parameter
		case *Client:
			table.client = parameter
		case *logger.Logger:
			table.logger = parameter.Child("user", "table", "id", table.ID)
		}
	}
	if table.logger == nil {
		table.logger = logger.Create("gcloudcx", &logger.NilStream{})
	}
}

// GetID gets the identifier of this
//
//	implements Identifiable
func (table DataTable) GetID() uuid.UUID {
	return table.ID
}

// GetURI gets the URI of this
//
// implements Addressable
func (table DataTable) GetURI(ids ...uuid.UUID) URI {
	if len(ids) > 0 {
		return NewURI("/api/v2/flows/datatables/%s", ids[0])
	}
	if table.ID != uuid.Nil {
		return NewURI("/api/v2/flows/datatables/%s", table.ID)
	}
	return URI("/api/v2/flows/datatables/")
}

// AddRow adds a row to this table
func (table DataTable) AddRow(context context.Context, key string, row DataTableRow) (correlationID string, err error) {
	if row == nil {
		return "", errors.ArgumentMissing.With("row")
	}
	row["key"] = key
	return table.client.SendRequest(
		context,
		NewURI("%s/rows", table.GetURI()),
		&request.Options{
			Method:      http.MethodPost,
			PayloadType: "application/json",
			Payload:     row,
		},
		nil,
	)
}

// UpdateRow updates a row to this table
func (table DataTable) UpdateRow(context context.Context, key string, row DataTableRow) (correlationID string, err error) {
	if row == nil {
		return "", errors.ArgumentMissing.With("row")
	}
	row["key"] = key
	return table.client.SendRequest(
		context,
		NewURI("%s/rows/%s", table.GetURI(), key),
		&request.Options{
			Method:      http.MethodPut,
			PayloadType: "application/json",
			Payload:     row,
		},
		nil,
	)
}

// DeleteRow deletes a row from this table
func (table DataTable) DeleteRow(context context.Context, key string) (correlationID string, err error) {
	return table.client.Delete(context, NewURI("%s/rows/%s", table.GetURI(), key), nil)
}

// GetRows gets the rows of this table
func (table DataTable) GetRows(context context.Context) (rows []DataTableRow, correlationID string, err error) {
	rows = make([]DataTableRow, 0)

	entities := Entities{}
	page := uint64(1)

	for {
		uri := NewURI("%s/rows", table.GetURI()).WithQuery(Query{"pageNumber": page}).WithQuery(Query{"showbrief": false})
		if correlationID, err = table.client.Get(context, uri, &entities); err != nil {
			return []DataTableRow{}, correlationID, err
		}
		for _, entity := range entities.Entities {
			var row DataTableRow
			if err = json.Unmarshal(entity, &row); err != nil {
				return []DataTableRow{}, correlationID, err
			}
			rows = append(rows, row)
		}
		if page++; page > entities.PageCount {
			break
		}
	}

	return rows, correlationID, nil
}

// GetRow gets a row from this table
func (table DataTable) GetRow(context context.Context, key string) (row DataTableRow, correlationID string, err error) {
	uri := NewURI("%s/rows/%s", table.GetURI(), key).WithQuery(Query{"showbrief": false})

	if correlationID, err = table.client.Get(context, uri, &row); err != nil {
		return nil, correlationID, errors.WrapErrors(NotFoundError.With("row.key", key), err)
	}
	return row, correlationID, nil
}

// String gets a string representation of this
//
// implements fmt.Stringer
func (table DataTable) String() string {
	return table.Name
}

// MarshalJSON marshals this into JSON
//
// implements json.Marshaler
func (table DataTable) MarshalJSON() ([]byte, error) {
	type surrogate DataTable
	data, err := json.Marshal(&struct {
		surrogate
		SelfURI URI `json:"selfUri"`
	}{
		surrogate: surrogate(table),
		SelfURI:   table.GetURI(),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

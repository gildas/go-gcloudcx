package gcloudcx

import (
	"context"
	"encoding/json"
)

type Entities struct {
	Entities   [][]byte `json:"-"`
	PageSize   int64    `json:"pageSize"`
	PageNumber int64    `json:"pageNumber"`
	PageCount  uint64   `json:"pageCount"`
	PageTotal  uint64   `json:"total"`
	FirstURI   string   `json:"firstUri"`
	SelfURI    string   `json:"selfUri"`
	LastURI    string   `json:"lastUri"`
}

func (client *Client) FetchEntities(context context.Context, uri URI) ([][]byte, error) {
	entities := Entities{}
	values := [][]byte{}

	page := uint64(1)
	for {
		if err := client.Get(context, uri.WithQuery(Query{"pageNumber": page}), &entities); err != nil {
			return nil, err
		}
		values = append(values, entities.Entities...)
		if page++; page > entities.PageCount {
			break
		}
	}
	return values, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (entities *Entities) UnmarshalJSON(data []byte) error {
	type surrogate Entities
	var inner struct {
		surrogate
		Entities []json.RawMessage `json:"entities"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}
	*entities = Entities(inner.surrogate)
	entities.Entities = make([][]byte, 0, len(inner.Entities))
	for _, entity := range inner.Entities {
		entities.Entities = append(entities.Entities, entity)
	}
	return nil
}

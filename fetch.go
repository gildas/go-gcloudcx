package gcloudcx

import (
	"context"
	"encoding/json"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

// Fetch fetches a resource from the Genesys Cloud API
//
// # The object must implement the Fetchable interface
//
// Resources can be fetched by their ID:
//
//	user, err := Fetch[gcloudcx.User](context, client, uuid.UUID)
//
//	user, err := Fetch[gcloudcx.User](context, client, gcloudcx.User{ID: uuid.UUID})
//
// or by their URI:
//
//	user, err := Fetch[gcloudcx.User](context, client, gcloudcx.User{}.GetURI(uuid.UUID))
func Fetch[T Fetchable, PT interface {
	Initializable
	*T
}](context context.Context, client *Client, parameters ...any) (*T, error) {
	id, query, selfURI, log := parseFetchParameters(context, client, parameters...)

	if len(selfURI) > 0 {
		var object T
		if err := client.Get(context, selfURI.WithQuery(query), &object); err != nil {
			return nil, err
		}
		PT(&object).Initialize(client, log)
		return &object, nil
	}
	if id != uuid.Nil {
		var object T
		if err := client.Get(context, object.GetURI(id).WithQuery(query), &object); err != nil {
			return nil, err
		}
		PT(&object).Initialize(client, log)
		return &object, nil
	}
	return nil, errors.NotFound.WithStack()
}

// FetchWithStringID fetches a resource from the Genesys Cloud API
//
// # The object must implement the Fetchable interface
//
// Resources can be fetched by their ID:
//
//	integrationType, err := Fetch[gcloudcx.IntegrationType](context, client, stringid)
//
// or by their URI:
//
//	integrationType, err := Fetch[gcloudcx.IntegrationType](context, client, gcloudcx.IntegrationType{}.GetURI(stringid))
func FetchWithStringID[T FetchableByStringID, PT interface {
	Initializable
	*T
}](context context.Context, client *Client, parameters ...any) (*T, error) {
	id, query, selfURI, log := parseFetchParametersWithNamedID(context, client, parameters...)

	if len(selfURI) > 0 {
		var object T
		if err := client.Get(context, selfURI.WithQuery(query), &object); err != nil {
			return nil, err
		}
		PT(&object).Initialize(client, log)
		return &object, nil
	}
	if id != "" {
		var object T
		if err := client.Get(context, object.GetURI(id).WithQuery(query), &object); err != nil {
			return nil, err
		}
		PT(&object).Initialize(client, log)
		return &object, nil
	}
	return nil, errors.NotFound.WithStack()
}

// FetchBy fetches a resource from the Genesys Cloud API by a match function
//
// The resource must implement the Fetchable interface
//
//	match := func(user gcloudcx.User) bool {
//	    return user.Name == "John Doe"
//	}
//	user, err := FetchBy(context, client, match)
//
// A gcloudcx.Query can be added to narrow the request:
//
//	user, err := FetchBy(context, client, match, gcloudcx.Query{Language: "en-US"})
func FetchBy[T Fetchable, PT interface {
	Initializable
	*T
}](context context.Context, client *Client, match func(T) bool, parameters ...interface{}) (*T, error) {
	if match == nil {
		return nil, errors.ArgumentMissing.With("match function")
	}
	_, query, _, log := parseFetchParameters(context, client, parameters...)
	entities := Entities{}
	page := uint64(1)
	var addressable T
	for {
		uri := addressable.GetURI().WithQuery(query).WithQuery(Query{"pageNumber": page})
		if err := client.Get(context, uri, &entities); err != nil {
			return nil, err
		}
		for _, entity := range entities.Entities {
			var object T
			if err := json.Unmarshal(entity, &object); err == nil && match(object) {
				PT(&object).Initialize(client, log)
				return &object, nil
			}
		}
		if page++; page > entities.PageCount {
			break
		}
	}
	return nil, errors.NotFound.WithStack()
}

// FetchAll fetches all objects from the Genesys Cloud API
//
// The objects must implement the Fetchable interface
//
//	users, err := FetchAll[gcloudcx.User](context, client)
//
// A gcloudcx.Query can be added to narrow the request:
//
//	users, err := FetchAll[gcloudcx.User](context, client, gcloudcx.Query{Language: "en-US"})
func FetchAll[T Fetchable, PT interface {
	Initializable
	*T
}](context context.Context, client *Client, parameters ...interface{}) ([]*T, error) {
	_, query, _, log := parseFetchParameters(context, client, parameters...)
	entities := Entities{}
	objects := []*T{}
	page := uint64(1)
	var addressable T
	for {
		uri := addressable.GetURI().WithQuery(query).WithQuery(Query{"pageNumber": page})
		if err := client.Get(context, uri, &entities); err != nil {
			return nil, err
		}
		for _, entity := range entities.Entities {
			var object T
			if err := json.Unmarshal(entity, &object); err == nil {
				PT(&object).Initialize(client, log)
				objects = append(objects, &object)
			}
		}
		if page++; page > entities.PageCount {
			break
		}
	}
	return objects, nil
}

/*
func (client *Client) FetchAll(context context.Context, object Addressable) ([]interface{}, error) {
	entities := struct {
		Entities   []json.RawMessage `json:"entities"`
		PageSize   int               `json:"pageSize"`
		PageNumber int               `json:"pageNumber"`
		PageCount  int               `json:"pageCount"`
		PageTotal  int               `json:"total"`
		FirstURI   string            `json:"firstUri"`
		SelfURI    string            `json:"selfUri"`
		LastURI    string            `json:"lastUri"`
	}{}
	page := 1
	for {
		uri := URI(addressable.GetURI().Base().String() + "?messengerType=" + messengerType + "&pageNumber=" + strconv.FormatUint(page, 10))
		if err := client.Get(context, uri, &entities); err != nil {
			return nil, err
		}
		log.Record("response", entities).Infof("Got a response")

	}
	return []interface{}{}, nil
}
*/

func parseFetchParameters(context context.Context, client *Client, parameters ...any) (uuid.UUID, Query, URI, *logger.Logger) {
	var id uuid.UUID
	var query Query
	var uri URI
	log, _ := logger.FromContext(context)

	if log == nil {
		log = client.Logger
	}
	for _, parameter := range parameters {
		switch parameter := parameter.(type) {
		case uuid.UUID:
			id = parameter
		case Query:
			query = parameter
		case URI:
			uri = parameter
		case *logger.Logger:
			log = parameter
		default:
			if identifiable, ok := parameter.(Identifiable); ok {
				id = identifiable.GetID()
			} else if addressable, ok := parameter.(Addressable); ok {
				uri = addressable.GetURI()
			}
		}
	}
	return id, query, uri, log
}

func parseFetchParametersWithNamedID(context context.Context, client *Client, parameters ...any) (string, Query, URI, *logger.Logger) {
	var id string
	var query Query
	var uri URI
	log, _ := logger.FromContext(context)

	if log == nil {
		log = client.Logger
	}
	for _, parameter := range parameters {
		switch parameter := parameter.(type) {
		case string:
			id = parameter
		case Query:
			query = parameter
		case URI:
			uri = parameter
		case *logger.Logger:
			log = parameter
		default:
			if identifiable, ok := parameter.(core.StringIdentifiable); ok {
				id = identifiable.GetID()
			} else if addressable, ok := parameter.(Addressable); ok {
				uri = addressable.GetURI()
			}
		}
	}
	return id, query, uri, log
}

package gcloudcx

import (
	"context"
	"strings"
	"time"

	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

// ResponseManagementLibrary is the interface for the Response Management Library
//
// See https://developer.genesys.cloud/api/rest/v2/responsemanagement
type ResponseManagementLibrary struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	DateCreated time.Time        `json:"dateCreated,omitempty"`
	CreatedBy   *DomainEntityRef `json:"createdBy,omitempty"`
	Version     int              `json:"version"`
	SelfURI     URI              `json:"selfUri,omitempty"`
}

// Fetch fetches this from the given Client
//
//  implements Fetchable
func (library *ResponseManagementLibrary) Fetch(ctx context.Context, client *Client, parameters ...interface{}) error {
	id, name, uri, log := client.ParseParameters(ctx, library, parameters...)
	if len(uri) > 0 {
		return client.Get(ctx, uri, library)
	}
	if id != uuid.Nil {
		return client.Get(ctx, NewURI("/responsemanagement/libraries/%s", id), library)
	}
	if len(name) > 0 {
		nameLowercase := strings.ToLower(name)
		entities := struct {
			Libraries []ResponseManagementLibrary `json:"entities"`
			paginatedEntities
		}{}
		for pageNumber := 1; ; pageNumber++ {
			log.Debugf("Fetching page %d", pageNumber)
			if err := client.Get(ctx, NewURI("/responsemanagement/libraries/?pageNumber=%d", pageNumber), &entities); err != nil {
				return err
			}
			for _, entity := range entities.Libraries {
				if strings.Compare(strings.ToLower(entity.Name), nameLowercase) == 0 {
					*library = entity
					return nil
				}
			}
			if pageNumber >= entities.PageCount {
				break
			}
		}
	}
	return errors.NotFound.With("ResponseManagementLibrary")
}

// GetID gets the identifier of this
//
//   implements Identifiable
func (library *ResponseManagementLibrary) GetID() uuid.UUID {
	return library.ID
}

// GetType gets the identifier of this
//
//   implements Identifiable
func (library *ResponseManagementLibrary) GetType() string {
	return "responsemanagement.library"
}

// String gets a string version
//
//   implements the fmt.Stringer interface
func (library ResponseManagementLibrary) String() string {
	if len(library.Name) > 0 {
		return library.Name
	}
	return library.ID.String()
}

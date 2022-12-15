# go-gcloudcx

![GoVersion](https://img.shields.io/github/go-mod/go-version/gildas/go-gcloudcx)
[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/gildas/go-gcloudcx) 
[![License](https://img.shields.io/github/license/gildas/go-gcloudcx)](https://github.com/gildas/go-gcloudcx/blob/master/LICENSE) 
[![Report](https://goreportcard.com/badge/github.com/gildas/go-gcloudcx)](https://goreportcard.com/report/github.com/gildas/go-gcloudcx)  

A Package to send requests to HTTP/REST services.

Genesys Cloud CX Client Library in GO

Have a look at the examples/ folder for complete examples on how to use this library.

## Usage

You first start by creating a `gcloudcx.Client` that will allow to send requests to Genesys Cloud:  
```go
log    := logger.Create("gcloudcx")
client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       log,
})
```

You can choose the authorization grant right away as well:  
```go
log    := logger.Create("gcloudcx")
client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       log,
}).SetAuthorizationGrant(&gcloudcx.AuthorizationCodeGrant{
	ClientID:    "hlkjshdgpiuy123387",
	Secret:      "879e8ugspojdgj",
	RedirectURL: "http://my.acme.com/token",
})
```

Or,  
```go
log    := logger.Create("gcloudcx")
client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       log,
}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
	ClientID: "jklsdufg89u9j234",
	Secret:   "sdfgjlskdfjglksdfjg",
})
```

As of today, *Authorization Code* and *Client Credentials* grants are implemented.

In the case of the Authorization Code, the best is to run a Webserver in your code and to handle the authentication requests in the router. The library provides two helpers to manage the authentication:

- `AuthorizeHandler()` that can be used to ensure a page has an authenticated client,
- `LoggedInHandler()` that can be used in the *RedirectURL* to process the results of the authentication.

They can be used like this (using the [gorilla/mux](https://github.com/gorilla/mux) router, for example):  
```go
router := mux.NewRouter()
// This route is used as the RedirectURL of the client
router.Methods("GET").Path("/token").Handler(client.LoggedInHandler()(myhandler()))

authorizedRouter := router.PathPrefix("/").Subrouter()

authorizedRouter.Use(client.AuthorizeHandler())

// This is the main route of this application, we want a fully functional gcloudcx.Client
authorizedRouter.Methods("GET").Path("/").HandlerFunc(mainRouteHandler)
```

In you *HttpHandler*, the client will be available from the request's context:  
```go
func mainRouteHandler(w http.ResponseWriter, r *http.Request) {
	client, err := gcloudcx.ClientFromContext(r.Context())
	if err != nil {
		core.RespondWithError(w, http.StatusServiceUnavailable, err)
		return
	}

	// Let's get my organization here, as an example...
	organization, err := client.GetMyOrganization(r.Context())
	if err != nil {
		core.RespondWithError(w, http.StatusServiceUnavailable, err)
		return
	}
	core.RespondWithJSON(w, http.StatusOK, struct {
		OrgName  string `json:"organization"`
	}{
		OrgName:  organization.String(),
	})
}
```

When using the Client Credential grant, you can configure the grant to tell when the token gets updated, allowing you to store it and re-use it in the future:  
```go
var TokenUpdateChan = make(chan gcloudcx.UpdatedAccessToken)

func main() {
	// ...
	defer close(TokenUpdateChan)

	go func() {
		for {
			data, opened := <-TokenUpdateChan

			if !opened {
				return
			}

			log.Printf("Received Token: %s\n", data.Token)

			myID, ok := data.CustomData.(string)
			if (!ok) {
				log.Printf("Shoot....\n")
			} else {
				log.Printf("myID: %s\n", myID)
			}
		}
	}()

	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		// ...
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: "1234",
		Secret:   "s3cr3t",
		TokenUpdated: TokenUpdateChan,
		CustomData:   "myspecialID",
		Token: savedToken,
	})

	// do stuff with the client, etc.
}
```

As you can see, you can even pass some custom data (`interface{}`, so anything really) to the grant and that data will be passed back to the `func` that handles the `chan`.

## Using Go's contexts

All functions that will end up calling the Genesys Cloud API must use a [context](https://pkg.go.dev/context) as their first argument.

This allow developers to, eventually, control timeouts, and other things.

A useful pattern is to add a [logger](https://github.com/gildas/go-logger) to the context with various records, when the library function executes, it will send its logs to that logger, thus using all the records that were set:
```go
log := logger.Create("MYAPP").Record("coolId", myID)

user := client.GetMyUser(log.ToContext(someContext))
```
In the logs, you will see the value of `coolId` in every line produced by client.GetMyUser.

If the context does not contain a logger, the client.Logger is used.

## Fetch resources

The library provides a `Fetch` function that will fetch a resource from the Genesys Cloud API.

```go
// Fetch a user by its ID
user, err := gcloudcx.Fetch[gcloud.User](context, client, userID) // userID is a uuid.UUID
```

You can use `FetchBy` to fetch a resource by a specific field:
```go
integration, err := gcloudcx.FetchBy(context, client, func (integration gcloudcx.OpenMessagingIntegration) bool {
	return integration.Name == "My Integration"
})
```

**Note:** This method can be rather slow as it fetches all the resources of the type and then filters them.

You can also add query criteria to the `FetchBy` function:
```go
recipient, err := gcloudcx.FetchBy(
	context,
	client,
	func (recipient gcloudcx.Recipient) bool {
		return recipient.Name == "My Recipient"
	},
	gcloudcx.Query{
		"messengerType": "open",
		"pageSize":      100,
	},
)
```

Finally, you can use `FetchAll` to fetch all the resources of a type:
```go
integrations, err := gcloudcx.FetchAll[gcloudcx.OpenMessagingIntegration](context, client)
```

Again, some query criteria can be added:
```go
integrations, err := gcloudcx.FetchAll[gcloudcx.OpenMessagingIntegration](context, client, gcloudcx.Query{
	"pageSize": 100,
})
```

## Notifications

The Genesys Cloud Notification API is accessible via the `NotificationChannel` and `NotificationTopic` types.

Here is a quick example:  
```go
user, err := client.GetMyUser(context.Background())
if err != nil {
	log.Errorf("Failed to retrieve my User", err)
	panic(err)
}

notificationChannel, err := client.CreateNotificationChannel(context.Background())
if err != nil {
	log.Errorf("Failed to create a notification channel", err)
	panic(err)
}

topics, err := config.NotificationChannel.Subscribe(
	context.Background(),
	gcloudcx.UserPresenceTopic{}.TopicFor(user),
)
if err != nil {
	log.Errorf("Failed to subscribe to topics", err)
	panic(err)
}
log.Infof("Subscribed to topics: [%s]", strings.Join(topics, ","))

// Call the Genesys Cloud Notification Topic loop
go func() {
	for {
		select {
		case receivedTopic := notificationChannel.TopicReceived:
			if receivedTopic == nil {
				log.Infof("Terminating Topic message loop")
				return
			}
			switch topic := receivedTopic.(type) {
			case *gcloudcx.UserPresenceTopic:
				log.Infof("User %s, Presence: %s", topic.User, topic.Presence)
			default:
				log.Warnf("Unknown topic: %s", topic)
			}
		case <-time.After(30 * time.Second):
			// This timer makes sure the for loop does not execute too quickly 
			//   when no topic is received for a while
			log.Debugf("Nothing in the last 30 seconds")
		}
	}
}()
```

## Response Management (Canned Responses)

Responses canbe fetched, like any other resource, via the `Fetch` function:

```go
response, err := gcloudcx.Fetch[gcloudcx.Response](context, client, responseID)
```

Or via the dedicated `FetchByFilter` function:

```go
response, err = gcloudcx.ResponseManagementResponse{}.FetchByFilters(
	context,
	client,
	gcloudcx.ResponseManagementQueryFilter{
		Name: "name", Operator: "EQUALS", Values: []string{response.Name},
	},
)
```
This calls the query API of Genesys Cloud, which is more efficient than fetching all the responses and then filtering them (`FetchBy` method).

Responses can apply substitutions to provide the final text that can be sent to the user. Custom substitutions are provided as a `map[string]string` to the func and default substitutions are provided by the resource itself.

```go
text, err := response.ApplySubstitutions(context, "text/plain", map[string]string{
	"firstName": "John",
	"lastName":  "Doe",
})
```
The second argument allows to specify the content type of the text to be returned, that content type has to be present in the resource, otherwise an error is returned. `text/plain` and `text/html` are supported by the Genesys Cloud API as of today.

`ApplySubstitutions` supports both the Genesys Cloud substitutions (`{{substitutionId}}`) and the Go Template system, which is far more powerful. See the [Go Template documentation](https://pkg.go.dev/text/template) for more information. On top of the basic features, `ApplySubstitutions` also provides the functions from the Sprig library, which is a collection of useful functions for templates. See the [Sprig documentation](https://masterminds.github.io/sprig/) for more information.

For example the following response will return the text `Hello John Doe`:
```go
response := gcloudcx.ResponseManagementResponse{ // This is a fake response, just to show the template
	Texts: []gcloudcx.ResponseManagementContent{
		{
			ContentType: "text/plain",
			Content: `Hello {{firstName}} {{lastName}}`, // Genesys Cloud substitutions
		}
	},
}
text, err := response.ApplySubstitutions(context, "text/plain", map[string]string{
	"firstName": "John",
	"lastName":  "Doe",
})
assert.Equal(t, "Hello John Doe", text)
```

To get the same result with Go Templates:
```go
response := gcloudcx.ResponseManagementResponse{ // This is a fake response, just to show the template
	Texts: []gcloudcx.ResponseManagementContent{
		{
			ContentType: "text/plain",
			Content: `Hello {{.firstName}} {{.lastName}}`, // Go Template substitutions
		}
	},
}
text, err := response.ApplySubstitutions(context, "text/plain", map[string]string{
	"firstName": "John",
	"lastName":  "Doe",
})
assert.Equal(t, "Hello John Doe", text)
```

You can mix both approaches:
```go
response := gcloudcx.ResponseManagementResponse{ // This is a fake response, just to show the template
	Texts: []gcloudcx.ResponseManagementContent{
		{
			ContentType: "text/plain",
			Content: `
Hello {{.firstName}} {{.lastName}}, you are on {{location}}!   # Genesys Cloud and Go substitutions
{{if eq .location "Earth"}}You are on the right place!{{end}}  # Go Template condition
And you are visiting {{ default "Atlantis" .country }}.        # Sprig function
`,
		}
	},
	Substitutions: []gcloudcx.ResponseManagementSubstitution{{
		ID: "location",
		Description: "The location of the person to greet",
		Default: "Earth",
	}},
}
text, err := response.ApplySubstitutions(context, "text/plain", map[string]string{
	"firstName": "John",
	"lastName":  "Doe",
})
```

## Agent Chat API

## Guest Chat API

## OpenMessaging API

# TODO

This library implements only a very small set of Genesys Cloud CX's API at the moment, but I keep adding stuff...
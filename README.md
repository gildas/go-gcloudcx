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

## Agent Chat API

## Guest Chat API

## OpenMessaging API

# TODO

This library implements only a very small set of Genesys Cloud CX's API at the moment, but I keep adding stuff...
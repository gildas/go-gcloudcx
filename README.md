# go-purecloud

![GoVersion](https://img.shields.io/github/go-mod/go-version/gildas/go-purecloud)
[![GoDoc](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/gildas/go-purecloud) 
[![License](https://img.shields.io/github/license/gildas/go-purecloud)](https://github.com/gildas/go-purecloud/blob/master/LICENSE) 
[![Report](https://goreportcard.com/badge/github.com/gildas/go-purecloud)](https://goreportcard.com/report/github.com/gildas/go-purecloud)  

A Package to send requests to HTTP/REST services.

PureCloud Client Library in GO

Have a look at the examples/ folder for complete examples on how to use this library.

## Usage


You first start by creating a `purecloud.Client` that will allow to send requests to PureCloud:  
```go
Log    := logger.Create("purecloud")
client := purecloud.NewClient(&purecloud.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       Log,
})
```

You can choose the authorization grant right away as well:  
```go
Log    := logger.Create("purecloud")
client := purecloud.NewClient(&purecloud.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       Log,
}).SetAuthorizationGrant(&purecloud.AuthorizationCodeGrant{
	ClientID:    "hlkjshdgpiuy123387",
	Secret:      "879e8ugspojdgj",
	RedirectURL: "http://my.acme.com/token",
})
```

Or,  
```go
Log    := logger.Create("purecloud")
client := purecloud.NewClient(&purecloud.ClientOptions{
	DeploymentID: "123abc0981234i0df8g0",
	Logger:       Log,
}).SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
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
// This is the main route of this application, we want a fully functional purecloud.Client
router.Methods("GET").Path("/").Handler(Client.AuthorizeHandler()(mainRouteHandler()))
// This route is used as the RedirectURL of the client
router.Methods("GET").Path("/token").Handler(Client.LoggedInHandler()(myhandler()))
```

In you *HttpHandler*, the client will be available from the request's context:  
```go
func mainRouteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Let's get my organization here, as an example...
		organization, err := client.GetMyOrganization()
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		core.RespondWithJSON(w, http.StatusOK, struct {
			OrgName  string `json:"organization"`
		}{
			OrgName:  organization.String(),
		})
	})
}
```

## Notifications

The PureCloud Notification API is accessible via the `NotificationChannel` and `NotificationTopic` types.

Here is a quick example:  
```go
user, err := client.GetMyUser()
if err != nil {
	log.Errorf("Failed to retrieve my User", err)
	panic(err)
}

notificationChannel, err := client.CreateNotificationChannel()
if err != nil {
	log.Errorf("Failed to create a notification channel", err)
	panic(err)
}

topics, err := config.NotificationChannel.Subscribe(
	purecloud.UserPresenceTopic{}.TopicFor(user),
)
if err != nil {
	log.Errorf("Failed to subscribe to topics", err)
	panic(err)
}
log.Infof("Subscribed to topics: [%s]", strings.Join(topics, ","))

// Call the PureCloud Notification Topic loop
go func() {
	for {
		select {
		case receivedTopic := notificationChannel.TopicReceived:
			if receivedTopic == nil {
				log.Infof("Terminating Topic message loop")
				return
			}
			switch topic := receivedTopic.(type) {
			case *purecloud.UserPresenceTopic:
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

# TODO

This library implements only a very small set of PureCloud's API at the moment, but I keep adding stuff...
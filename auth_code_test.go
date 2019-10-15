package purecloud_test

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
	"github.com/gorilla/mux"
)

// Log is the application Logger
var Log *logger.Logger

// Client is the PureCloud Client
var Client *purecloud.Client

func loggedInHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Scope("logged_in")

		// Do stuff here...
		// For Example go to the root path:
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}

func mainRouteHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, err := logger.FromContext(r.Context())
		if err != nil {
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		log = log.Scope("main")

		client, err := purecloud.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the PureCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Let's get my user and organization here, as an example...

		organization, err := client.GetMyOrganization()
		if err != nil {
			log.Errorf("Failed to retrieve my Organization", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
		}

		user, err := client.GetMyUser()
		if err != nil {
			log.Errorf("Failed to retrieve my User", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
		}
		core.RespondWithJSON(w, http.StatusOK, struct{
			UserName string `json:"user"`
			OrgName  string `json:"organization"`
		}{
			UserName: user.ID,
			OrgName: organization.Name,
		})
		return
	})
}

func ExampleAuthorizationCodeGrant() {
	var (
		region       = flag.String("region", core.GetEnvAsString("REGION", "mypurecloud.com"), "the PureCloud Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("CLIENTID", ""), "the PureCloud Client ID for authentication")
		secret       = flag.String("secret", core.GetEnvAsString("SECRET", ""), "the PureCloud Client Secret for authentication")
		deploymentID = flag.String("deploymentid", core.GetEnvAsString("DEPLOYMENTID", ""), "the PureCloud Application Deployment ID")
		root         = flag.String("root", core.GetEnvAsString("ROOT", ""), "The root uri to give to PureCloud as a Redirect URI")
		port         = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
	)
	Log := logger.Create("AuthCode_Example")

	if len(*root) == 0 {
		*root = fmt.Sprintf("http://localhost:%d", *port)
	}
	redirectURL, err := url.Parse(*root + "/token")
	if err != nil {
		Log.Fatalf("Invalid Redirect URL: %s/token", *root, err)
		os.Exit(-1)
	}

	Client = purecloud.New(purecloud.ClientOptions{
		Region:       *region,
		DeploymentID: *deploymentID,
		Logger:       Log,
	}).SetAuthorizationGrant(&purecloud.AuthorizationCodeGrant{
		ClientID:    *clientID,
		Secret:      *secret,
		RedirectURL: redirectURL,
	})

	// Create the HTTP Incoming Request Router
	router := mux.NewRouter().StrictSlash(true)
	// This route actually performs login the user using the grant of the purecloud.Client
	//   Upon success, your route httpHandler is called
	router.Methods("GET").Path("/token").Handler(Log.HttpHandler()(Client.LoginHandler()(loggedInHandler())))

	// This route performs your actions, but makes sure the client is authorized,
	//   if authorized, your route http.Handler is called
	//   otherwise, the AuthorizeHandler will redirect the user to the PureCloud Login page
	//   that will end up with the grant.RedirectURL defined ealier
	router.Methods("GET").Path("/").Handler(Log.HttpHandler()(Client.AuthorizeHandler()(mainRouteHandler())))

	WebServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", *port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
		//ErrorLog:     Log,
	}

	// Starting the server
	go func() {
		log := Log.Topic("webserver").Scope("run")

		log.Infof("Starting WEB server on port %d", *port)
		if err := WebServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Fatalf("Failed to start the WEB server on port: %d", *port, err)
			}
		}
	}()
	fmt.Println("WEB Server started")
	// Output:
	// WEB Server started

	// Accepting shutdowns from SIGINT (^C) and SIGTERM (docker, heroku)
	interruptChannel := make(chan os.Signal, 1)
	exitChannel      := make(chan struct{})

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// The go routine that wait for cleaning stuff when exiting
	go func() {
		sig := <-interruptChannel // Block until we have to stop

		context, cancel := context.WithTimeout(context.Background(), 15 * time.Second)
		defer cancel()

		Log.Infof("Application is stopping (%+v)", sig)

		// Stopping the WEB server
		Log.Debugf("WEB server is shutting down")
		WebServer.SetKeepAlivesEnabled(false)
		WebServer.Shutdown(context)
		Log.Infof("WEB server is stopped")
		close(exitChannel)
	}()

	<- exitChannel
	os.Exit(0)
}
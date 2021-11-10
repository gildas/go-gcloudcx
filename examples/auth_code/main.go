package main

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
	"github.com/gildas/go-gcloudcx"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Log is the application Logger
var Log *logger.Logger

// Client is the GCloud Client
var Client *gcloudcx.Client

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
		log.Infof("Redirecting to /")
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

		client, err := gcloudcx.ClientFromContext(r.Context())
		if err != nil {
			log.Errorf("Failed to retrieve the GCloud Client", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		// Let's get my user and organization here, as an example...

		organization, err := client.GetMyOrganization(context.Background())
		if err != nil {
			log.Errorf("Failed to retrieve my Organization", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}

		user, err := client.GetMyUser(r.Context(), "presence")
		if err != nil {
			log.Errorf("Failed to retrieve my User", err)
			core.RespondWithError(w, http.StatusServiceUnavailable, err)
			return
		}
		core.RespondWithJSON(w, http.StatusOK, struct {
			UserName string `json:"user"`
			Presence string `json:"presence"`
			OrgName  string `json:"organization"`
		}{
			UserName: user.String(),
			Presence: user.Presence.String(),
			OrgName:  organization.String(),
		})
	})
}

func main() {
	var (
		region       = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the GCloud CX Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the GCloud CX Client ID for authentication")
		secret       = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the GCloud CX Client Secret for authentication")
		deploymentID = flag.String("deploymentid", core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""), "the GCloud CX Application Deployment ID")
		redirectRoot = flag.String("redirecturi", core.GetEnvAsString("PURECLOUD_REDIRECTURI", ""), "The root uri to give to GCloud CX as a Redirect URI")
		port         = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
	)
	flag.Parse()

	Log = logger.Create("AuthCode_Example")

	Log.Infof("redirect root: %s", *redirectRoot)
	if len(*redirectRoot) == 0 {
		*redirectRoot = fmt.Sprintf("http://localhost:%d", *port)
	}
	redirectURL, err := url.Parse(*redirectRoot + "/token")
	if err != nil {
		Log.Fatalf("Invalid Redirect URL: %s/token", *redirectRoot, err)
		os.Exit(-1)
	}
	Log.Infof("Make sure your GCloud OAUTH accepts redirects to: %s", redirectURL.String())

	Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       *region,
		DeploymentID: uuid.MustParse(*deploymentID),
		Logger:       Log,
	}).SetAuthorizationGrant(&gcloudcx.AuthorizationCodeGrant{
		ClientID:    uuid.MustParse(*clientID),
		Secret:      *secret,
		RedirectURL: redirectURL,
	})

	// Create the HTTP Incoming Request Router
	router := mux.NewRouter().StrictSlash(true)

	router.Use(Log.HttpHandler())

	// This route actually performs login the user using the grant of the gcloudcx.Client
	//   Upon success, your route httpHandler is called
	router.Methods("GET").Path("/token").Handler(Client.LoggedInHandler()(loggedInHandler()))

	// This route performs your actions, but makes sure the client is authorized,
	//   if authorized, your route http.Handler is called
	//   otherwise, the AuthorizeHandler will redirect the user to the GCloud Login page
	//   that will end up with the grant.RedirectURL defined earlier
	router.Methods("GET").Path("/").Handler(Client.AuthorizeHandler()(mainRouteHandler()))

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

	// Accepting shutdowns from SIGINT (^C) and SIGTERM (docker, heroku)
	interruptChannel := make(chan os.Signal, 1)
	exitChannel := make(chan struct{})

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)

	// The go routine that wait for cleaning stuff when exiting
	go func() {
		// Block until we have to stop
		for sig := range interruptChannel {
			switch sig {
			case syscall.SIGUSR1:
				Log.Flush()
			default:
				context, cancel := context.WithTimeout(context.Background(), 15*time.Second)

				Log.Infof("Application is stopping (%+v)", sig)

				// Stopping the WEB server
				Log.Debugf("WEB server is shutting down")
				WebServer.SetKeepAlivesEnabled(false)
				err := WebServer.Shutdown(context)
				if err != nil {
					Log.Errorf("Failed to stop the WEB server", err)
				} else {
					Log.Infof("WEB server is stopped")
				}
				Log.Flush()
				close(exitChannel)
				cancel()
			}
		}
	}()

	<-exitChannel
	os.Exit(0)
}

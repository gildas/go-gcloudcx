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
	"github.com/gildas/go-purecloud"
	"github.com/gorilla/mux"
)

// Log is the application Logger
var Log *logger.Logger

// Client is the PureCloud Client
var Client *purecloud.Client

// The Organization this session belongs to
// The Queue to transfer to
var AgentQueue *purecloud.Queue

var WebRootPath string

func main() {
	var (
		region       = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the PureCloud Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the PureCloud Client ID for authentication")
		secret       = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the PureCloud Client Secret for authentication")
		deploymentID = flag.String("deploymentid", core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""), "the PureCloud Application Deployment ID")
		redirectRoot = flag.String("redirecturi", core.GetEnvAsString("PURECLOUD_REDIRECTURI", ""), "The root uri to give to PureCloud as a Redirect URI")
		queueName    = flag.String("queue", core.GetEnvAsString("PURECLOUD_QUEUE", ""), "The queue to transfer to")
		webrootpath  = flag.String("webrootpath", core.GetEnvAsString("WEBROOT_PATH", ""), "The path to use before each endpoint (useful for nginx config)")
		port         = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
	)
	flag.Parse()

	Log = logger.Create("ChatBot_Example")

	if len(*redirectRoot) == 0 {
		*redirectRoot = fmt.Sprintf("http://localhost:%d", *port)
	}
	redirectURL, err := url.Parse(*redirectRoot + "/token")
	if err != nil {
		Log.Fatalf("Invalid Redirect URL: %s/token", *redirectRoot, err)
		os.Exit(-1)
	}
	Log.Infof("Make sure your PureCloud OAUTH accepts redirects to: %s", redirectURL.String())

	WebRootPath = *webrootpath
	Log.Infof("Web root path: %s", WebRootPath)

	Client = purecloud.NewClient(purecloud.ClientOptions{
		Region:       *region,
		DeploymentID: *deploymentID,
		Logger:       Log,
	}).SetAuthorizationGrant(&purecloud.AuthorizationCodeGrant{
		ClientID:    *clientID,
		Secret:      *secret,
		RedirectURL: redirectURL,
	})

	// TODO: Make this better... Too Simple for now
	AgentQueue = &purecloud.Queue{Name: *queueName}

	// Create the HTTP Incoming Request Router
	router := mux.NewRouter().StrictSlash(true)
	// This route actually performs login the user using the grant of the purecloud.Client
	//   Upon success, your route httpHandler is called
	router.Methods("GET").Path("/token").Handler(Log.HttpHandler()(Client.LoginHandler()(LoggedInHandler())))
	// This route performs the login process makes sure the client is authorized,
	//   if authorized, the LoginHandler is called to setup some variables
	//   otherwise, the purecloud.AuthorizeHandler will redirect the user to the PureCloud Login page
	//   that will end up with the grant.RedirectURL defined ealier
	router.Methods("POST").Path("/login").Handler(Log.HttpHandler()(Client.AuthorizeHandler()(LoginHandler())))

	// This route shows the main page, with login infor and a Chat Widget
	//  See: https://developer.mypurecloud.com/api/webchat/widget-version2.html
	router.Methods("GET").Path("/").Handler(Log.HttpHandler()(Client.HttpHandler()(MainHandler())))

	// This route gives the PureCloud Widget Javascript config to use
	//  See: https://developer.mypurecloud.com/api/webchat/widget-version2.html
	router.Methods("GET").Path("/widget").Handler(Log.HttpHandler()(Client.HttpHandler()(WidgetHandler())))

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

	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// The go routine that wait for cleaning stuff when exiting
	go func() {
		sig := <-interruptChannel // Block until we have to stop

		context, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

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
		close(exitChannel)
	}()

	<-exitChannel
	os.Exit(0)
}

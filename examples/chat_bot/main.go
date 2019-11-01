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

// MyAppConfig holds the configuration of this Application
var MyAppConfig *AppConfig

// NotFoundHandler is called when all other routes did not match
func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log, _ := logger.FromContext(r.Context())
		log = log.Topic("route").Scope("notfound")
		log.Errorf("Route not Found %s", r.URL.String())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Path Not Found"))
	})
}

func main() {
	var (
		region         = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the PureCloud Region. \nDefault: mypurecloud.com")
		clientID       = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the PureCloud Client ID for authentication")
		secret         = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the PureCloud Client Secret for authentication")
		deploymentID   = flag.String("deploymentid", core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", ""), "the PureCloud Application Deployment ID")
		redirectRoot   = flag.String("redirecturi", core.GetEnvAsString("PURECLOUD_REDIRECTURI", ""), "The root uri to give to PureCloud as a Redirect URI")
		agentQueueName = flag.String("agentqueue", core.GetEnvAsString("PURECLOUD_AGENTQUEUE", ""), "The queue to transfer to agents")
		botURL         = flag.String("boturl", core.GetEnvAsString("PURECLOUD_BOTURL", ""), "The Bot URL to query for interpretation")
		botQueueName   = flag.String("botqueue", core.GetEnvAsString("PURECLOUD_BOTQUEUE", ""), "The queue to send customers to initially")
		queueName      = flag.String("queue", core.GetEnvAsString("PURECLOUD_QUEUE", ""), "(legacy) the queue to send to")
		webrootpath    = flag.String("webrootpath", core.GetEnvAsString("WEBROOT_PATH", ""), "The path to use before each endpoint (useful for nginx config)")
		port           = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
		err error
	)
	flag.Parse()

	Log = logger.Create("ChatBot_Example")

	if len(*agentQueueName) == 0 && len(*botQueueName) == 0 && len(*queueName) == 0 {
		Log.Fatalf("Agent and Bot Queues are empty")
		os.Stderr.WriteString("Agent and Bot Queues are empty, please use --agentqueue and/or --botqueue\n")
		os.Exit(-2)
	} else if len(*agentQueueName) == 0 && len(*botQueueName) == 0 {
		Log.Warnf("Queue %s will be used for both the Agents and the Bot (legacy), you should use --agentqueue and --botqueue", *queueName)
		fmt.Fprintf(os.Stderr, "Queue %s will be used for both the Agents and the Bot (legacy), you should use --agentqueue and --botqueue\n", *queueName)
		agentQueueName = queueName
		botQueueName   = queueName
	} else if len(*agentQueueName) == 0 {
		Log.Warnf("Bot Queue %s will be used also for the Agent Queue", *botQueueName)
		agentQueueName = botQueueName
	} else if len(*botQueueName) == 0 {
		Log.Warnf("Agent Queue %s will be used also for the Bot Queue", *agentQueueName)
		botQueueName = agentQueueName
	}

	MyAppConfig = &AppConfig{
		AgentQueue:  &purecloud.Queue{Name: *agentQueueName},
		BotQueue:    &purecloud.Queue{Name: *botQueueName},
		WebRootPath: *webrootpath,
		Logger:      Log.Topic("config"),
	}
	if MyAppConfig.BotURL, err = url.Parse(*botURL); err != nil {
		Log.Fatalf("The Chat BOT URL %s is invalid", *botURL, err)
		fmt.Fprintf(os.Stderr, "The Chat BOT URL %s is invalid. Error: %s", *botURL, err)
		os.Exit(-2)
	}

	if len(MyAppConfig.WebRootPath) > 0 {
		Log.Infof("Web root path: %s, Please use this path in your NGINX", MyAppConfig.WebRootPath)
	}

	if len(*redirectRoot) == 0 {
		*redirectRoot = fmt.Sprintf("http://localhost:%d", *port)
	}
	redirectURL, err := url.Parse(*redirectRoot + "/token")
	if err != nil {
		Log.Fatalf("Invalid Redirect URL: %s/token", *redirectRoot, err)
		os.Exit(-1)
	}
	Log.Infof("Make sure your PureCloud OAUTH accepts redirects to: %s", redirectURL.String())

	Client = purecloud.NewClient(purecloud.ClientOptions{
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
	// All routes will use the Logger and the AppConfig
	router.Use(Log.HttpHandler())
	router.Use(MyAppConfig.HttpHandler())

	// This route actually performs login the user using the grant of the purecloud.Client
	//   Upon success, your route httpHandler is called
	router.Methods("GET").Path("/token").Handler(Client.LoggedInHandler()(LoggedInHandler()))

	// This route performs the login process makes sure the client is authorized,
	//   if authorized, the LoggedInHandler is called to setup some variables
	//   otherwise, the purecloud.AuthorizeHandler will redirect the user to the PureCloud Login page
	//   that will end up with the grant.RedirectURL defined ealier
	router.Methods("POST").Path("/login").Handler(Client.AuthorizeHandler()(LoggedInHandler()))

	// This route performs the logout process
	router.Methods("POST").Path("/logout").Handler(Client.LogoutHandler()(LoggedOutHandler()))

	// This route shows the main page, with login infor and a Chat Widget
	//  See: https://developer.mypurecloud.com/api/webchat/widget-version2.html
	router.Methods("GET").Path("/").Handler(Client.HttpHandler()(MainHandler()))

	// This route gives the PureCloud Widget Javascript config to use
	//  See: https://developer.mypurecloud.com/api/webchat/widget-version2.html
	router.Methods("GET").Path("/widget").Handler(Client.HttpHandler()(WidgetHandler()))

	// This route catches all other routes
	router.PathPrefix("/").Handler(NotFoundHandler())

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

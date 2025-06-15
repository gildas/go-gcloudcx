package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func UpdateEnvFile(config *Config) {
	config.Client.Logger.Infof("Updating the .env file")
	_ = godotenv.Write(map[string]string{
		"PURECLOUD_REGION":       config.Client.Region,
		"PURECLOUD_CLIENTID":     config.Client.Grant.(*gcloudcx.ClientCredentialsGrant).ClientID.String(),
		"PURECLOUD_CLIENTSECRET": config.Client.Grant.(*gcloudcx.ClientCredentialsGrant).Secret,
		"PURECLOUD_CLIENTTOKEN":  config.Client.Grant.AccessToken().Token,
		"PURECLOUD_DEPLOYMENTID": config.Client.DeploymentID.String(),
		"INTEGRATION_NAME":       config.IntegrationName,
		"INTEGRATION_WEBHOOK":    config.IntegrationWebhookURL.String(),
		"INTEGRATION_TOKEN":      config.IntegrationWebhookToken,
	}, ".env")
}

func main() {
	_ = godotenv.Load()
	var (
		region       = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the GENESYS Cloud Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the GENESYS Cloud Client ID for authentication")
		clientSecret = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the GENESYS Cloud Client Secret for authentication")
		clientToken  = flag.String("token", core.GetEnvAsString("PURECLOUD_CLIENTTOKEN", ""), "the GENESYS Cloud Client Token if any. If expired, it will be replaced")

		integrationName  = flag.String("integration", core.GetEnvAsString("INTEGRATION_NAME", ""), "the Integration Name")
		integrationHook  = flag.String("webhook", core.GetEnvAsString("INTEGRATION_WEBHOOK", ""), "the Integration Webhook URL")
		integrationToken = flag.String("webhook-token", core.GetEnvAsString("INTEGRATION_TOKEN", ""), "the Integration Webhook Token")

		port          = flag.Int("port", core.GetEnvAsInt("PORT", 3000), "the port to listen to")
		wait          = flag.Duration("graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish")
		err           error
		correlationID string
	)
	flag.Parse()

	log := logger.Create("OpenMessaging_Example", logger.TRACE)
	defer log.Flush()
	log.Infof("%s", strings.Repeat("-", 80))
	log.Infof("log Destination: %s", log)
	log.Infof("Webserver Port=%d", *port)

	if *port == 0 {
		log.Fatalf("Missing Webserver port, stopping...")
		os.Exit(-1)
	}

	// Initializing the Config
	config := &Config{
		IntegrationName:         *integrationName,
		IntegrationWebhookURL:   core.Must(url.Parse(*integrationHook)),
		IntegrationWebhookToken: *integrationToken,
		Client: gcloudcx.NewClient(&gcloudcx.ClientOptions{
			Region: *region,
			Logger: log,
		}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
			ClientID: uuid.MustParse(*clientID),
			Secret:   *clientSecret,
			Token: gcloudcx.AccessToken{
				Type:  "bearer",
				Token: *clientToken,
			},
		}),
	}
	defer UpdateEnvFile(config)

	// Initializing the OpenMessaging Integration
	config.Client.Logger.Infof("Fetching OpenMessaging Integration %s", *integrationName)
	match := func(integration gcloudcx.OpenMessagingIntegration) bool {
		return integration.Name == *integrationName
	}
	config.Integration, correlationID, err = gcloudcx.FetchBy(context.Background(), config.Client, match)
	log = log.Record("genesys-correlation", correlationID)

	if errors.Is(err, errors.NotFound) {
		log.Infof("Creating a new OpenMessaging Integration for %s", *integrationName)
		config.Integration, correlationID, err = config.Client.CreateOpenMessagingIntegration(context.Background(), config.IntegrationName, config.IntegrationWebhookURL, config.IntegrationWebhookToken, nil)
		log = log.Record("genesys-correlation", correlationID)
		if err != nil {
			log.Fatalf("Failed creating integration", err)
			os.Exit(1)
		}
		log.Record("integration", config.Integration).Infof("Created new integration")
	} else if err != nil {
		log.Fatalf("Failed to retrieve OpenMessaging Integration", err)
		os.Exit(1)
	}

	if strings.Compare(config.Integration.WebhookURL.String(), config.IntegrationWebhookURL.String()) != 0 || strings.Compare(config.Integration.WebhookToken, config.IntegrationWebhookToken) != 0 {
		log.Warnf("OpenMessaging Integration has changed, we need to update it in GENESYS Cloud")
		if correlationID, err := config.Integration.Update(context.Background(), config.IntegrationName, config.IntegrationWebhookURL, config.IntegrationWebhookToken); err != nil {
			log = log.Record("genesys-correlation", correlationID)
			log.Fatalf("Failed to update the OpenMessaging Integration")
			os.Exit(1)
		}
		log.Record("integration", config.Integration).Infof("Updated integration")
	}

	// Setting up web routes
	router := mux.NewRouter().StrictSlash(true)
	router.Use(log.HttpHandler())
	router.Use(config.HttpHandler())
	router.Methods("POST").Path("/hook").HandlerFunc(mainRouteHandler)

	// Routes for the internal Chat Server (used by the chat web client)
	ChatRoutes(router)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public/"))))
	router.PathPrefix("/").Handler(NotFoundHandler())

	// Initializing the web server
	webServer := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", *port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	// Start the Chat Server
	log.Infof("Starting Chat server")
	config.ChatServer = NewChatServer(router, log)
	go config.ChatServer.Start(config)

	// Starting the server
	go func() {
		log := log.Child("webserver", "run")

		log.Infof("Starting WEB server on port %d", *port)
		log.Infof("Serving routes:")
		_ = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			message := strings.Builder{}
			args := []interface{}{}

			if methods, err := route.GetMethods(); err == nil {
				message.WriteString("%s ")
				args = append(args, strings.Join(methods, ","))
			} else {
				return nil
			}
			if path, err := route.GetPathTemplate(); err == nil {
				message.WriteString("%s ")
				args = append(args, path)
			}
			if path, err := route.GetPathRegexp(); err == nil {
				message.WriteString("%s ")
				args = append(args, path)
			}
			log.Infof(message.String(), args...)
			return nil
		})
		if err := webServer.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				log.Fatalf("Failed to start the WEB server on port: %d", *port, err)
			}
		}
	}()

	// Accepting shutdowns from SIGINT (^C) and SIGTERM (docker, heroku)
	interruptChannel := make(chan os.Signal, 1)
	exitChannel := make(chan struct{})
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	// Waiting to clean stuff up when exiting
	go func() {
		sig := <-interruptChannel // Block until we have to stop
		context, cancel := context.WithTimeout(context.Background(), *wait)
		defer cancel()

		log.Infof("Application is stopping (%+v)", sig)

		// Stopping the WEB server
		log.Debugf("WEB server is shutting down")
		webServer.SetKeepAlivesEnabled(false)
		if err = webServer.Shutdown(context); err != nil {
			log.Errorf("Failed to stop WEB server", err)
		} else {
			log.Infof("WEB server is stopped")
		}

		// Stopping the application
		close(exitChannel)
	}()

	<-exitChannel
	os.Exit(0)
}

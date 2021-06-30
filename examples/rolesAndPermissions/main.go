package main

import (
	"flag"
	"strings"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gildas/go-purecloud"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Config struct {
	Client *purecloud.Client
}

func UpdateEnvFile(config *Config) {
	config.Client.Logger.Infof("Updating the .env file")
	_ = godotenv.Write(map[string]string{
		"PURECLOUD_REGION":       config.Client.Region,
		"PURECLOUD_CLIENTID":     config.Client.Grant.(*purecloud.ClientCredentialsGrant).ClientID.String(),
		"PURECLOUD_CLIENTSECRET": config.Client.Grant.(*purecloud.ClientCredentialsGrant).Secret,
		"PURECLOUD_CLIENTTOKEN":  config.Client.Grant.AccessToken().Token,
	}, ".env")
}

func main() {
	_ = godotenv.Load()
	var (
		region       = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the GENESYS Cloud Region. \nDefault: mypurecloud.com")
		clientID     = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the GENESYS Cloud Client ID for authentication")
		clientSecret = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the GENESYS Cloud Client Secret for authentication")
		clientToken  = flag.String("token", core.GetEnvAsString("PURECLOUD_CLIENTTOKEN", ""), "the GENESYS Cloud Client Token if any. If expired, it will be replaced")
	)
	flag.Parse()

	log := logger.Create("RolesAndPermissions_Example", &logger.StdoutStream{FilterLevel: logger.TRACE, Unbuffered: true})
	defer log.Flush()
	log.Infof(strings.Repeat("-", 80))
	log.Infof("Log Destination: %s", log)

	// Initializing the Config
	config := &Config{
		Client: purecloud.NewClient(&purecloud.ClientOptions{
			Region:       *region,
			Logger:       log,
		}).SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
			ClientID: uuid.MustParse(*clientID),
			Secret:   *clientSecret,
			Token:    purecloud.AccessToken{
				Type:  "bearer",
				Token: *clientToken,
			},
		}),
	}
	defer UpdateEnvFile(config)

	log.Infof("Permissions: %d", len(flag.Args()))
	if permitted, missing := config.Client.CheckPermissions(flag.Args()...); len(missing) == 0 {
		log.Infof("You can do %s", strings.Join(permitted, ", "))
	} else {
		log.Errorf("You are missing %s", strings.Join(missing, ", "))
	}
}

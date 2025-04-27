package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/gildas/go-core"
	"github.com/gildas/go-gcloudcx"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
)

func main() {
	var (
		region         = flag.String("region", core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com"), "the GCloud CX Region. \nDefault: mypurecloud.com")
		clientID       = flag.String("clientid", core.GetEnvAsString("PURECLOUD_CLIENTID", ""), "the GCloud CX Client ID for authentication")
		secret         = flag.String("secret", core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""), "the GCloud CX Client Secret for authentication")
		conversationID = flag.String("conversationid", core.GetEnvAsString("PURECLOUD_CONVERSATIONID", ""), "the GCloud CX Conversation ID")
		err            error
	)
	flag.Parse()

	log := logger.Create("Disconnect_Example")
	mainctx := log.ToContext(context.Background())

	if len(*conversationID) == 0 {
		log.Errorf("You must provide both a Conversation ID and a Participant ID")
		os.Exit(1)
	}

	client := gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region: *region,
		Logger: log,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: uuid.MustParse(*clientID),
		Secret:   *secret,
	})

	// Let's make sure the conversation ecists
	conversation, err := gcloudcx.Fetch[gcloudcx.Conversation](mainctx, client, uuid.MustParse(*conversationID))
	if err != nil {
		log.Errorf("Failed to get conversation: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to get conversation: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Conversation: %s\n", conversation)

	var agent *gcloudcx.Participant

	for _, participant := range conversation.Participants {
		fmt.Printf("Participant: %s (%s), %s\n", participant.Name, participant.ID, participant.Purpose)
		if participant.Purpose == "agent" {
			agent = &participant
		}
	}
	if agent == nil {
		log.Errorf("No agent found in conversation")
		fmt.Fprintf(os.Stderr, "No agent found in conversation\n")
		os.Exit(1)
	}

	err = conversation.Disconnect(mainctx, agent)
	if err != nil {
		log.Errorf("Failed to disconnect from conversation: %v", err)
		fmt.Fprintf(os.Stderr, "Failed to disconnect from conversation: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Disconnected from conversation\n")
}

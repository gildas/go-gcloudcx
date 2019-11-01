package purecloud_test

import (
	"testing"

	"github.com/gildas/go-logger"
	"github.com/stretchr/testify/require"

	purecloud "github.com/gildas/go-purecloud"
)

func TestCanExtractClientAndLogger(t *testing.T) {
	logger := CreateLogger("test-helpers")
	client := purecloud.NewClient(purecloud.ClientOptions{
		DeploymentID: "12345676890",
		Logger:       logger,
	})
	_, _, err := purecloud.ExtractClientAndLogger(client)
	if err != nil {
		logger.Errorf("Failed", err)
	}
	require.Nil(t, err, "Failed to fetch stuff")
}

func TestCanRunInitializable(t *testing.T) {
	logger := CreateLogger("test-helpers")
	client := purecloud.NewClient(purecloud.ClientOptions{
		DeploymentID: "12345676890",
		Logger:       logger,
	})

	stuff := &Stuff{}
	err := stuff.Initialize(client)
	if err != nil {
		logger.Errorf("Failed", err)
	}
	require.Nil(t, err, "Failed to fetch stuff")
}

func TestCanInitializeWithFetch(t *testing.T) {
	logger := CreateLogger("test-helpers")
	client := purecloud.NewClient(purecloud.ClientOptions{
		DeploymentID: "12345676890",
		Logger:       logger,
	})

	stuff := &Stuff{}
	err := client.Fetch(stuff)
	if err != nil {
		logger.Errorf("Failed", err)
	}
	require.Nil(t, err, "Failed to fetch stuff")
}

type Stuff struct {
	ID     string            `json:"id"`
	Client *purecloud.Client `json:"-"`
	Logger *logger.Logger    `json:"-"`
}

func (stuff *Stuff) Initialize(parameters ...interface{}) error {
	client, logger, err := purecloud.ExtractClientAndLogger(parameters...)
	if err != nil {
		return err
	}
	stuff.Client = client
	stuff.Logger = logger
	return nil
}

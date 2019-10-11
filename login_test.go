package purecloud_test

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-logger"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	purecloud "github.com/gildas/go-purecloud"
)

type LoginSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	Client *purecloud.Client
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (suite *LoginSuite) TestCanLogin() {
	err := suite.Client.SetAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID: core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		Secret:   core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	}).Login()
	suite.Assert().Nil(err, "Failed to login")
}

func (suite *LoginSuite) TestFailsLoginWithInvalidGrant() {
	err := suite.Client.LoginWithAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID:  "DEADID",
		Secret:    "WRONGSECRET",
	})
	suite.Assert().NotNil(err, "Should have failed login in")

	apierr, ok := errors.Cause(err).(purecloud.APIError)
	suite.Require().True(ok, "Error is not a purecloud.APIError")
	suite.Logger.Record("apierr", apierr).Errorf("API Error", err)
	suite.Assert().Equal(400, apierr.Status)
	suite.Assert().Equal("bad.credentials", apierr.Code)
}

func (suite *LoginSuite) TestCanLoginWithClientCredentialsGrant() {
	err := suite.Client.LoginWithAuthorizationGrant(&purecloud.ClientCredentialsGrant{
		ClientID:  core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		Secret:    core.GetEnvAsString("PURECLOUD_CLIENTSECRET", ""),
	})
	suite.Assert().Nil(err, "Failed to login")
}

func (suite *LoginSuite) TestCanLoginWithImplicitGrant() {
	redirectURL, _ := url.Parse("http://localhost:35000/oauth2/callback")
	grant := &purecloud.ImplicitGrant{
		ClientID:    core.GetEnvAsString("PURECLOUD_CLIENTID", ""),
		RedirectURL: redirectURL,

	}

	endTest := make(chan struct{})

	// Setting up web routes
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("GET").Path("/").Handler(suite.Logger.HttpHandler()(grant.HttpHandler()(func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log, err := logger.FromContext(r.Context())
			if err != nil {
				log := suite.Logger
			}
			log.Infof("Handling Authentication Callback")
			grant, err := purecloud.AuthorizationGrantFromContext(r.Context())
			suite.Assert().Nil(err, "Failed to retrieve Authorization Grant from Request Context")
			if err != nil {
				core.RespondWithError(w, http.StatusServiceUnavailable, errors.New("Failed to retrieve Grant"))
				return
			}
			core.RespondWithJSON(w, http.StatusOK, struct{}{})
			return
		})
	})))
	router.Methods("GET").Path("/token/{token}").Handler(suite.Logger.HttpHandler()(grant.HttpHandler()(func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log, err := logger.FromContext(r.Context())
			if err != nil {
				log := suite.Logger
			}
			log.Infof("Handling Incoming Token")
			grant, err := purecloud.AuthorizationGrantFromContext(r.Context())
			suite.Assert().Nil(err, "Failed to retrieve Authorization Grant from Request Context")
			if err != nil {
				log.Errorf("Failed to retrieve Grant from request")
				core.RespondWithError(w, http.StatusServiceUnavailable, errors.New("Failed to retrieve Grant"))
				return
			}

			params := mux.Vars(r)
			tokenString, ok := params["token"]
			if !ok {
				log.Errorf("Parameter token was missing")
				core.RespondWithError(w, http.StatusBadRequest, errors.New("Parameter token was missing"))
				return
			}
			core.RespondWithJSON(w, http.StatusOK, struct{}{})
			return
		})
	})))

	// Setting up the server
	WebServer := &http.Server{
		Addr:         "0.0.0.0:35000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
		//ErrorLog:     Log,
	}

	suite.Client.LoginWithAuthorizationGrant(grant)

	<- endTest
	context, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	WebServer.SetKeepAlivesEnabled(false)
	WebServer.Shutdown(context)

	token := grant.AccessToken()
	suite.Require().NotNil(token)
	suite.Assert().NotEmpty(token.Token)
}

// Suite Tools

func (suite *LoginSuite) SetupSuite() {
	suite.Name = strings.TrimSuffix(reflect.TypeOf(*suite).Name(), "Suite")
	logFolder := filepath.Join(".", "log")
	os.MkdirAll(logFolder, os.ModePerm)
	suite.Logger = logger.CreateWithDestination("test", fmt.Sprintf("file://%s/test-%s.log", logFolder, strings.ToLower(suite.Name)))
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	var (
		region       = core.GetEnvAsString("PURECLOUD_REGION", "")
		clientID     = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
		secret       = core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
		deploymentID = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	)
	suite.Client = purecloud.New(purecloud.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	})
	suite.Require().NotNil(suite.Client, "PureCloud Client is nil")

}

func (suite *LoginSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
}

func (suite *LoginSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))

	suite.Start = time.Now()
}

func (suite *LoginSuite) AfterTest(suiteName, testName string) {
	duration := time.Now().Sub(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

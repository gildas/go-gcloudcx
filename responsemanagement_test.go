package gcloudcx_test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
	"github.com/gildas/go-logger"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"

	"github.com/gildas/go-gcloudcx"
)

type ResponseManagementSuite struct {
	suite.Suite
	Name   string
	Logger *logger.Logger
	Start  time.Time

	LibraryName  string
	LibraryID    uuid.UUID
	ResponseName string
	ResponseID   uuid.UUID
	Client       *gcloudcx.Client
}

func TestResponseManagementSuite(t *testing.T) {
	suite.Run(t, new(ResponseManagementSuite))
}

// *****************************************************************************
// #region: Suite Tools {{{
func (suite *ResponseManagementSuite) SetupSuite() {
	var err error
	var value string

	_ = godotenv.Load()
	suite.Name = strings.TrimSuffix(reflect.TypeOf(suite).Elem().Name(), "Suite")
	suite.Logger = logger.Create("test",
		&logger.FileStream{
			Path:         fmt.Sprintf("./log/test-%s.log", strings.ToLower(suite.Name)),
			Unbuffered:   true,
			SourceInfo:   true,
			FilterLevels: logger.NewLevelSet(logger.TRACE),
		},
	).Child("test", "test")
	suite.Logger.Infof("Suite Start: %s %s", suite.Name, strings.Repeat("=", 80-14-len(suite.Name)))

	region := core.GetEnvAsString("PURECLOUD_REGION", "mypurecloud.com")

	value = core.GetEnvAsString("PURECLOUD_CLIENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_CLIENTID is not set")

	clientID, err := uuid.Parse(value)
	suite.Require().NoError(err, "PURECLOUD_CLIENTID is not a valid UUID")

	secret := core.GetEnvAsString("PURECLOUD_CLIENTSECRET", "")
	suite.Require().NotEmpty(secret, "PURECLOUD_CLIENTSECRET is not set")

	value = core.GetEnvAsString("PURECLOUD_DEPLOYMENTID", "")
	suite.Require().NotEmpty(value, "PURECLOUD_DEPLOYMENTID is not set")

	deploymentID, err := uuid.Parse(value)
	suite.Require().NoError(err, "PURECLOUD_DEPLOYMENTID is not a valid UUID")

	value = core.GetEnvAsString("RESPONSE_MANAGEMENT_LIBRARY_ID", "")
	suite.Require().NotEmpty(value, "RESPONSE_MANAGEMENT_LIBRARY_ID is not set in your environment")

	suite.LibraryID, err = uuid.Parse(value)
	suite.Require().NoError(err, "RESPONSE_MANAGEMENT_LIBRARY_ID is not a valid UUID")

	suite.LibraryName = core.GetEnvAsString("RESPONSE_MANAGEMENT_LIBRARY_NAME", "")
	suite.Require().NotEmpty(suite.LibraryName, "RESPONSE_MANAGEMENT_LIBRARY_NAME is not set in your environment")

	value = core.GetEnvAsString("RESPONSE_MANAGEMENT_RESPONSE_ID", "")
	suite.Require().NotEmpty(value, "RESPONSE_MANAGEMENT_RESPONSE_ID is not set in your environment")

	suite.ResponseID, err = uuid.Parse(value)
	suite.Require().NoError(err, "RESPONSE_MANAGEMENT_RESPONSE_ID is not a valid UUID")

	suite.ResponseName = core.GetEnvAsString("RESPONSE_MANAGEMENT_RESPONSE_NAME", "")
	suite.Require().NotEmpty(suite.ResponseName, "RESPONSE_MANAGEMENT_RESPONSE_NAME is not set in your environment")

	suite.Client = gcloudcx.NewClient(&gcloudcx.ClientOptions{
		Region:       region,
		DeploymentID: deploymentID,
		Logger:       suite.Logger,
	}).SetAuthorizationGrant(&gcloudcx.ClientCredentialsGrant{
		ClientID: clientID,
		Secret:   secret,
	})
	suite.Require().NotNil(suite.Client, "GCloudCX Client is nil")
}

func (suite *ResponseManagementSuite) TearDownSuite() {
	if suite.T().Failed() {
		suite.Logger.Warnf("At least one test failed, we are not cleaning")
		suite.T().Log("At least one test failed, we are not cleaning")
	} else {
		suite.Logger.Infof("All tests succeeded, we are cleaning")
	}
	suite.Logger.Infof("Suite End: %s %s", suite.Name, strings.Repeat("=", 80-12-len(suite.Name)))
	suite.Logger.Close()
}

func (suite *ResponseManagementSuite) BeforeTest(suiteName, testName string) {
	suite.Logger.Infof("Test Start: %s %s", testName, strings.Repeat("-", 80-13-len(testName)))
	suite.Start = time.Now()

	// Reuse tokens as much as we can
	if !suite.Client.IsAuthorized() {
		suite.Logger.Infof("Client is not logged in...")
		err := suite.Client.Login(context.Background())
		suite.Require().NoError(err, "Failed to login")
		suite.Logger.Infof("Client is now logged in...")
	} else {
		suite.Logger.Infof("Client is already logged in...")
	}
}

func (suite *ResponseManagementSuite) AfterTest(suiteName, testName string) {
	duration := time.Since(suite.Start)
	suite.Logger.Record("duration", duration.String()).Infof("Test End: %s %s", testName, strings.Repeat("-", 80-11-len(testName)))
}

// #endregion: Suite Tools }}}
// *****************************************************************************

func (suite *ResponseManagementSuite) TestCanFetchLibraryByID() {
	library, err := gcloudcx.Fetch[gcloudcx.ResponseManagementLibrary](context.Background(), suite.Client, suite.LibraryID)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoErrorf(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.LibraryID, library.GetID(), "Library ID is not the same")
	suite.Assert().Equal(suite.LibraryName, library.String(), "Library Name is not the same")
	suite.Logger.Record("library", library).Infof("Library Details")
}

func (suite *ResponseManagementSuite) TestCanFetchLibraryByName() {
	match := func(library gcloudcx.ResponseManagementLibrary) bool {
		return library.Name == suite.LibraryName
	}
	library, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoErrorf(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.LibraryID, library.GetID(), "Library ID is not the same")
	suite.Assert().Equal(suite.LibraryName, library.String(), "Library Name is not the same")
	suite.Logger.Record("library", library).Infof("Library Details")
}

func (suite *ResponseManagementSuite) TestCanFetchResponseByID() {
	response, err := gcloudcx.Fetch[gcloudcx.ResponseManagementResponse](context.Background(), suite.Client, suite.ResponseID)
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoErrorf(err, "Failed to fetch Response Management Library, Error: %s", err)
	suite.Assert().Equal(suite.ResponseID, response.GetID(), "Client's Organization ID is not the same")
	suite.Assert().Equal(suite.ResponseName, response.String(), "Client's Organization Name is not the same")
	suite.Logger.Record("response", response).Infof("Response Details")
}

func (suite *ResponseManagementSuite) TestCanFetchResponseByFilters() {
	response, err := gcloudcx.ResponseManagementResponse{}.FetchByFilters(context.Background(), suite.Client, gcloudcx.ResponseManagementQueryFilter{
		Name: "name", Operator: "EQUALS", Values: []string{suite.ResponseName},
	})
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoErrorf(err, "Failed to fetch Response Management Response, Error: %s", err)
	suite.Assert().Equal(suite.ResponseID, response.GetID(), "Response ID is not the same")
	suite.Assert().Equal(suite.ResponseName, response.String(), "Response Name is not the same")
	suite.Logger.Record("response", response).Infof("Response Details")
}

func (suite *ResponseManagementSuite) TestCanFetchResponseByName() {
	match := func(response gcloudcx.ResponseManagementResponse) bool {
		return response.Name == suite.ResponseName
	}
	response, err := gcloudcx.FetchBy(context.Background(), suite.Client, match, gcloudcx.Query{"libraryId": suite.LibraryID})
	if err != nil {
		suite.Logger.Errorf("Failed", err)
	}
	suite.Require().NoErrorf(err, "Failed to fetch Response Management Response, Error: %s", err)
	suite.Assert().Equal(suite.ResponseID, response.GetID(), "Response ID is not the same")
	suite.Assert().Equal(suite.ResponseName, response.String(), "Response Name is not the same")
	suite.Logger.Record("response", response).Infof("Response Details")
}

func (suite *ResponseManagementSuite) TestShouldFailFetchingLibraryWithUnknownName() {
	match := func(library gcloudcx.ResponseManagementLibrary) bool {
		return library.Name == "unknown library"
	}
	_, err := gcloudcx.FetchBy(context.Background(), suite.Client, match)
	suite.Require().Error(err, "Should have failed to fetch Response Management Library")
	suite.Logger.Errorf("Expected Failure", err)
	suite.Assert().ErrorIs(err, errors.NotFound, "Should have failed to fetch Response Management Library")
}

func (suite *ResponseManagementSuite) TestShouldFailFetchingResponseWithUnknownName() {
	match := func(response gcloudcx.ResponseManagementResponse) bool {
		return response.Name == "unknown response"
	}
	_, err := gcloudcx.FetchBy(context.Background(), suite.Client, match, gcloudcx.Query{"libraryId": suite.LibraryID})
	suite.Require().Error(err, "Should have failed to fetch Response Management Response")
	suite.Logger.Errorf("Expected Failure", err)
	suite.Assert().ErrorIs(err, errors.NotFound, "Should have failed to fetch Response Management Response")
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutions() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     "Hello, {{name}}",
			},
		},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"name": "John"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, John", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithoutPlaceholder() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     "Hello, World!",
			},
		},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"name": "John"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, World!", text)

	text, err = response.ApplySubstitutions(ctx, "text/plain", nil)
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, World!", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithGOPlaceholders() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     "Hello, {{.name}}",
			},
		},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"name": "John"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, John", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithDefaults() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     `Hello, {{name}}`,
			},
		},
		Substitutions: []gcloudcx.ResponseManagementSubstitution{{
			ID: "name",
			Description: "The name of the person to greet",
			Default: "John",
		}},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"lastname": "Doe"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, John", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithGODefaults() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     `Hello, {{default "John" .name}}`,
			},
		},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"lastname": "Doe"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, John", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithGOAction() {
	ctx := suite.Logger.ToContext(context.Background())
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     `Hello, {{if .name}}{{.name}}{{else}}John{{end}}`,
			},
		},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", map[string]string{"lastname": "Doe"})
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal("Hello, John", text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithComplexTemplate() {
	ctx := suite.Logger.ToContext(context.Background())
	expected := `
{
  "genesys_prompt": "Would you like to buy now?",
  "genesys_quick_replies": [{
    "text": "OK","payload": "answer=OK"
  },{
    "text": "Cancel","payload": "answer=Cancel"
  }]
}`
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     `
{
  "genesys_prompt": "{{question}}",
  "genesys_quick_replies": [{
{{- if .OK_payload}}
    "text": "{{OK}}","payload": "{{OK_payload}}"
{{- else}}
    "text": "{{OK}}","payload": "answer={{OK}}"
{{- end}}
  },{
{{- if .Cancel_payload}}
    "text": "{{Cancel}}","payload": "{{Cancel_payload}}"
{{- else}}
    "text": "{{Cancel}}","payload": "answer={{Cancel}}"
{{- end}}
  }]
}`,
			},
		},
		Substitutions: []gcloudcx.ResponseManagementSubstitution{{
			ID: "question", Default: "Would you like to proceed?",
		}, {
			ID: "OK", Default: "OK",
		}, {
			ID: "OK_payload", Default: "",
		}, {
			ID: "Cancel", Default: "Cancel",
		}, {
			ID: "Cancel_payload", Default: "",
		}},
	}
	arguments := map[string]string{
		"question": "Would you like to buy now?",
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", arguments)
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal(expected, text)
}

func (suite *ResponseManagementSuite) TestCanApplySubstitutionsWithComplexTemplateAndNoArguments() {
	ctx := suite.Logger.ToContext(context.Background())
	expected := `
{
  "genesys_prompt": "Would you like to proceed?",
  "genesys_quick_replies": [{
    "text": "OK","payload": "answer=OK"
  },{
    "text": "Cancel","payload": "answer=Cancel"
  }]
}`
	response := gcloudcx.ResponseManagementResponse{
		Name: "Test",
		Texts: []gcloudcx.ResponseManagementContent{
			{
				ContentType: "text/plain",
				Content:     `
{
  "genesys_prompt": "{{question}}",
  "genesys_quick_replies": [{
{{- if .OK_payload}}
    "text": "{{OK}}","payload": "{{OK_payload}}"
{{- else}}
    "text": "{{OK}}","payload": "answer={{OK}}"
{{- end}}
  },{
{{- if .Cancel_payload}}
    "text": "{{Cancel}}","payload": "{{Cancel_payload}}"
{{- else}}
    "text": "{{Cancel}}","payload": "answer={{Cancel}}"
{{- end}}
  }]
}`,
			},
		},
		Substitutions: []gcloudcx.ResponseManagementSubstitution{{
			ID: "question", Default: "Would you like to proceed?",
		}, {
			ID: "OK", Default: "OK",
		}, {
			ID: "OK_payload", Default: "",
		}, {
			ID: "Cancel", Default: "Cancel",
		}, {
			ID: "Cancel_payload", Default: "",
		}},
	}
	text, err := response.ApplySubstitutions(ctx, "text/plain", nil)
	suite.Require().NoError(err, "Failed to apply substitutions")
	suite.Assert().Equal(expected, text)
}

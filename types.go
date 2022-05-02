package slackworkflowbot

import (
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type (
	ActionID      string
	BlockID       string
	CallbackID    string
	BotToken      string
	SigningSecret string
)

type (
	AppContext struct {
		Slack                           *slack.Client
		config                          configuration
		workflowStep                    WorkflowStepFunc
		workflowStepCallbackID          CallbackID
		replyWithConfigurationView      ReplyWithConfigurationView
		saveUserSettingsForWorkflowStep SaveUserSettingsForWorkflowStep
	}
	configuration struct {
		botToken      string
		signingSecret string
	}
)

type (
	Middleware = func(next http.Handler) http.Handler

	SecretsVerifierMiddleware struct {
		appCtx AppContext
	}
)

// var _ Middleware = &SecretsVerifierMiddleware{}
var _ http.Handler = &SecretsVerifierMiddleware{}

type (
	ReplyWithConfigurationView = func(
		appContext AppContext,
		message slack.InteractionCallback,
		privateMetadata string,
		externalID string,
	) error

	SaveUserSettingsForWorkflowStep = func(
		appContext AppContext,
		message slack.InteractionCallback,
	) error
)

type (
	WorkflowStepFunc = func(
		appContext AppContext,
		workflowStep slackevents.EventWorkflowStep,
	)
)

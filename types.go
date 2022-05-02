package slackworkflowbot

import (
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type (
	ActionID   string
	BlockID    string
	CallbackID string

	BotToken      string
	SigningSecret string
)

type (
	configuration struct {
		botToken      string
		signingSecret string
	}
	AppContext struct {
		Slack                  *slack.Client
		config                 configuration
		workflowStep           WorkflowStepFunc
		workflowStepCallbackID CallbackID
	}
)

type (
	SecretsVerifierMiddleware struct {
		appCtx  AppContext
		handler http.Handler
	}
)

var _ http.Handler = &SecretsVerifierMiddleware{}

type (
	ReplyWithConfigurationView = func(
		message slack.InteractionCallback,
		privateMetadata string,
		externalID string,
	) error

	SaveUserSettingsForWorkflowStep = func(
		workflowStepEditID string,
		inputs *slack.WorkflowStepInputs,
		outputs *[]slack.WorkflowStepOutput,
	) error
)

type (
	WorkflowStepFunc = func(
		workflowStep slackevents.EventWorkflowStep,
	)
)

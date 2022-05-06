package slackworkflowbot

import (
	"context"
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
		config        configuration
		configureStep ConfigureStepContext
		stepExecute   StepExecuteContext
	}
	configuration struct {
		botToken      string
		signingSecret string
	}
	StepExecuteContext struct {
		SlackClient  SlackWorkflowStepExecuteClient
		workflowStep map[CallbackID]WorkflowStepFunc
	}
	ConfigureStepContext struct {
		SlackClient SlackWorkfowConfigurationClient
		configView  map[CallbackID]ConfigView
		saveConfig  map[CallbackID]SaveConfig
	}
)

type (
	Middleware = func(next http.Handler) http.Handler

	SecretsVerifierMiddleware struct {
		appCtx configuration
	}
)

// var _ Middleware = NewSecretsVerifierMiddleware(appCtx).
var _ http.Handler = &SecretsVerifierMiddleware{}

type (
	SlackWorkfowConfigurationClient interface {
		OpenView(
			triggerID string,
			view slack.ModalViewRequest,
		) (*slack.ViewResponse, error)

		SaveWorkflowStepConfiguration(
			workflowStepEditID string,
			inputs *slack.WorkflowStepInputs,
			outputs *[]slack.WorkflowStepOutput,
		) error
	}

	// workflows.stepCompleted and workflows.stepFailed.
	SlackWorkflowStepExecuteClient interface {
		WorkflowStepCompleted(
			ctx context.Context,
			workflowStepExecuteID string,
			outputs *map[string]string,
		) error

		WorkflowStepFailed(
			ctx context.Context,
			workflowStepExecuteID string,
			errorMessage string,
		) error
	}
)

type (
	ConfigView = func(
		appCctx ConfigureStepContext,
		message slack.InteractionCallback,
		privateMetadata string,
		externalID string,
	) error

	SaveConfig = func(
		appCtx ConfigureStepContext,
		message slack.InteractionCallback,
	) error
)

type (
	WorkflowStepFunc = func(
		appContext StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	)
)

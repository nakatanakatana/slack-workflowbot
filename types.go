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
		config        configuration
		configureStep ConfigureStepContext
		stepExecute   StepExecuteContext
	}
	configuration struct {
		botToken      string
		signingSecret string
	}
	StepExecuteContext struct {
		SlackClient            SlackWorkflowStepExecuteClient
		workflowStep           WorkflowStepFunc
		workflowStepCallbackID CallbackID
	}
	ConfigureStepContext struct {
		SlackClient                     SlackWorkfowConfigurationClient
		replyWithConfigurationView      ReplyWithConfigurationView
		saveUserSettingsForWorkflowStep SaveUserSettingsForWorkflowStep
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
		OpenView(triggerID string,
			view slack.ModalViewRequest,
		) (*slack.ViewResponse, error)

		SaveWorkflowStepConfiguration(
			workflowStepEditID string,
			inputs *slack.WorkflowStepInputs,
			outputs *[]slack.WorkflowStepOutput,
		) error
	}
	// workflows.stepCompleted and workflows.stepFailed.
	SlackWorkflowStepExecuteClient interface{}
)

type (
	ReplyWithConfigurationView = func(
		appContext ConfigureStepContext,
		message slack.InteractionCallback,
		privateMetadata string,
		externalID string,
	) error

	SaveUserSettingsForWorkflowStep = func(
		appContext ConfigureStepContext,
		message slack.InteractionCallback,
	) error
)

type (
	WorkflowStepFunc = func(
		appContext StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	)
)

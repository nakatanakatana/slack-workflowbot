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
	Middleware = func(next http.Handler) http.Handler

	SecretsVerifierMiddleware struct {
		signingSecret SigningSecret
	}
)

// var _ Middleware = NewSecretsVerifierMiddleware(appCtx).
var _ http.Handler = &SecretsVerifierMiddleware{}

type (
	SlackWorkflowConfigurationClient interface {
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
		message slack.InteractionCallback,
		privateMetadata string,
		externalID string,
	) error
	ConfigViewFunctions map[CallbackID]ConfigView

	SaveConfig = func(
		message slack.InteractionCallback,
	) error
	SaveConfigFunctions map[CallbackID]SaveConfig

	WorkflowStep = func(
		workflowStep slackevents.EventWorkflowStep,
	)
	WorkflowStepFunctions map[CallbackID]WorkflowStep
)

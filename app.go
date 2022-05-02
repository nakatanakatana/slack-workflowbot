package slackworkflowbot

import (
	"github.com/slack-go/slack"
)

func CreateAppContext(
	botToken BotToken,
	signingSecret SigningSecret,
	workflowStep WorkflowStepFunc,
	workflowStepCallbackID CallbackID,
) AppContext {
	var appCtx AppContext
	appCtx.config.botToken = string(botToken)
	appCtx.config.signingSecret = string(signingSecret)

	appCtx.Slack = slack.New(appCtx.config.botToken)
	appCtx.workflowStep = workflowStep
	appCtx.workflowStepCallbackID = workflowStepCallbackID

	return appCtx
}

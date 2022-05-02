package slackworkflowbot

import (
	"github.com/slack-go/slack"
)

func CreateAppContext(
	botToken BotToken,
	signingSecret SigningSecret,
	workflowStep WorkflowStepFunc,
	workflowStepCallbackID CallbackID,
	replyWithConfigurationView ReplyWithConfigurationView,
	saveUserSettingsForWorkflowStep SaveUserSettingsForWorkflowStep,
) AppContext {
	var appCtx AppContext
	appCtx.config.botToken = string(botToken)
	appCtx.config.signingSecret = string(signingSecret)

	appCtx.Slack = slack.New(appCtx.config.botToken)
	appCtx.workflowStep = workflowStep
	appCtx.workflowStepCallbackID = workflowStepCallbackID
	appCtx.replyWithConfigurationView = replyWithConfigurationView
	appCtx.saveUserSettingsForWorkflowStep = saveUserSettingsForWorkflowStep

	return appCtx
}

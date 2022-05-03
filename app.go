package slackworkflowbot

import (
	slackdev "github.com/nakatanakatana/slack"
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

	slackClient := slack.New(appCtx.config.botToken)
	slackdevClient := slackdev.New(appCtx.config.botToken)

	appCtx.stepExecute.SlackClient = slackdevClient
	appCtx.stepExecute.workflowStep = workflowStep
	appCtx.stepExecute.workflowStepCallbackID = workflowStepCallbackID

	appCtx.configureStep.SlackClient = slackClient
	appCtx.configureStep.replyWithConfigurationView = replyWithConfigurationView
	appCtx.configureStep.saveUserSettingsForWorkflowStep = saveUserSettingsForWorkflowStep

	return appCtx
}

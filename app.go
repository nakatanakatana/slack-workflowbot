package slackworkflowbot

import "github.com/slack-go/slack"

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

	appCtx.stepExecute.SlackClient = slackClient
	appCtx.stepExecute.workflowStep = workflowStep
	appCtx.stepExecute.workflowStepCallbackID = workflowStepCallbackID

	appCtx.configureStep.SlackClient = slackClient
	appCtx.configureStep.replyWithConfigurationView = replyWithConfigurationView
	appCtx.configureStep.saveUserSettingsForWorkflowStep = saveUserSettingsForWorkflowStep

	return appCtx
}

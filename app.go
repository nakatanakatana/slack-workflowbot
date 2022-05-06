package slackworkflowbot

import (
	slackdev "github.com/nakatanakatana/slack"
	"github.com/slack-go/slack"
)

func CreateAppContext(
	botToken BotToken,
	signingSecret SigningSecret,
	workflowStep map[CallbackID]WorkflowStepFunc,
	configView map[CallbackID]ConfigView,
	saveConfig map[CallbackID]SaveConfig,
) AppContext {
	var appCtx AppContext
	appCtx.config.botToken = string(botToken)
	appCtx.config.signingSecret = string(signingSecret)

	slackClient := slack.New(appCtx.config.botToken)
	slackdevClient := slackdev.New(appCtx.config.botToken)

	appCtx.stepExecute.SlackClient = slackdevClient
	appCtx.stepExecute.workflowStep = workflowStep

	appCtx.configureStep.SlackClient = slackClient
	appCtx.configureStep.configView = configView
	appCtx.configureStep.saveConfig = saveConfig

	return appCtx
}

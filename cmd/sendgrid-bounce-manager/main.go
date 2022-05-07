package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	bot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
	"github.com/nakatanakatana/slack-workflowbot/cmd/sendgrid-bounce-manager/checker"
	"github.com/nakatanakatana/slack-workflowbot/cmd/sendgrid-bounce-manager/deleter"
)

const (
	APIBaseURL = "/api/v1"
)

func main() {
	botToken := bot.BotToken(os.Getenv("SLACK_BOT_TOKEN"))
	signingSecret := bot.SigningSecret(os.Getenv("SLACK_SIGNING_SECRET"))
	defaultSendGridClient := sendgrid.New()
	slackClient := bot.CreateSlackClient(botToken)
	slackDevClient := bot.CreateSlackDevClient(botToken)

	configView := bot.ConfigViewFunctions{
		checker.CallbackID: checker.CreateConfigView(slackClient),
		deleter.CallbackID: deleter.CreateConfigView(slackClient),
	}

	saveConfig := bot.SaveConfigFunctions{
		checker.CallbackID: checker.CreateSaveStepConfig(slackClient),
		deleter.CallbackID: deleter.CreateSaveStepConfig(slackClient),
	}

	workflowStep := bot.WorkflowStepFunctions{
		checker.CallbackID: checker.CreateStepFunc(slackDevClient, defaultSendGridClient),
		deleter.CallbackID: deleter.CreateStepFunc(slackDevClient, defaultSendGridClient),
	}

	handleInteraction := bot.CreateInteractionHandler(configView, saveConfig)
	handleWorkflowStep := bot.CreateEventsHandler(workflowStep)

	mux := http.NewServeMux()
	verifier := bot.NewSecretsVerifierMiddleware(signingSecret)
	mux.Handle(fmt.Sprintf("%s/interaction", APIBaseURL), verifier(handleInteraction))
	mux.Handle(fmt.Sprintf("%s/events", APIBaseURL), verifier(handleWorkflowStep))

	log.Printf("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

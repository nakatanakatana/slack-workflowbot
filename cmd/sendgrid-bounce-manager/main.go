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

	workflowStep := map[bot.CallbackID]bot.WorkflowStepFunc{
		checker.CallbackID: checker.CreateStepFunc(defaultSendGridClient),
		deleter.CallbackID: deleter.CreateStepFunc(defaultSendGridClient),
	}

	configView := map[bot.CallbackID]bot.ConfigView{
		checker.CallbackID: checker.CreateConfigView(),
		deleter.CallbackID: deleter.CreateConfigView(),
	}

	saveConfig := map[bot.CallbackID]bot.SaveConfig{
		checker.CallbackID: checker.SaveStepConfig,
		deleter.CallbackID: deleter.SaveStepConfig,
	}

	appCtx := bot.CreateAppContext(
		botToken,
		signingSecret,
		workflowStep,
		configView,
		saveConfig,
	)

	handleInteraction := bot.CreateInteractionHandler(appCtx)
	handleWorkflowStep := bot.CreateEventsHandler(appCtx)

	mux := http.NewServeMux()
	verifier := bot.NewSecretsVerifierMiddleware(appCtx)
	mux.Handle(fmt.Sprintf("%s/interaction", APIBaseURL), verifier(handleInteraction))
	mux.Handle(fmt.Sprintf("%s/events", APIBaseURL), verifier(handleWorkflowStep))

	log.Printf("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

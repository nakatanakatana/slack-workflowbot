package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	bot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"
)

const (
	APIBaseURL = "/api/v1"

	CheckBounceStep  = bot.CallbackID("check-bounce-step")
	DeleteBounceStep = bot.CallbackID("delete-bounce-step")
)

func main() {
	botToken := bot.BotToken(os.Getenv("SLACK_BOT_TOKEN"))
	signingSecret := bot.SigningSecret(os.Getenv("SLACK_SIGNING_SECRET"))
	defaultSendGridClient := sendgrid.New()

	workflowStep := map[bot.CallbackID]bot.WorkflowStepFunc{
		CheckBounceStep:  createCheckStepFunc(defaultSendGridClient),
		DeleteBounceStep: createDeleteStepFunc(defaultSendGridClient),
	}

	configView := map[bot.CallbackID]bot.ConfigView{
		CheckBounceStep:  createConfigView(),
		DeleteBounceStep: createConfigView(),
	}

	saveConfig := map[bot.CallbackID]bot.SaveConfig{
		CheckBounceStep:  saveStepConfig,
		DeleteBounceStep: saveStepConfig,
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

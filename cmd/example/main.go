package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
)

const (
	APIBaseURL = "/api/v1"
	// MyExampleWorkflowStepCallbackID is configured in slack (api.slack.com/apps).
	// Select your app or create a new one. Then choose menu "Workflow Steps"...
	MyExampleWorkflowStepCallbackID = slackworkflowbot.CallbackID("example-step")
)

func main() {
	botToken := slackworkflowbot.BotToken(os.Getenv("SLACK_BOT_TOKEN"))
	signingSecret := slackworkflowbot.SigningSecret(os.Getenv("SLACK_SIGNING_SECRET"))

	appCtx := slackworkflowbot.CreateAppContext(
		botToken,
		signingSecret,
		doHeavyLoad,
		MyExampleWorkflowStepCallbackID,
		replyWithConfigurationView,
		saveUserSettingsForWrokflowStep,
	)

	handleInteraction := slackworkflowbot.CreateHandleInteraction(appCtx)
	handleMyWorkflowStep := slackworkflowbot.CreateHandleWorkflowStep(appCtx)

	mux := http.NewServeMux()
	verifier := slackworkflowbot.NewSecretsVerifierMiddleware(appCtx)
	mux.Handle(fmt.Sprintf("%s/interaction", APIBaseURL), verifier(handleInteraction))
	mux.Handle(fmt.Sprintf("%s/%s", APIBaseURL, MyExampleWorkflowStepCallbackID), verifier(handleMyWorkflowStep))

	log.Printf("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

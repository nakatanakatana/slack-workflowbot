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

	appCtx := slackworkflowbot.CreateAppContext(botToken, signingSecret, doHeavyLoad, MyExampleWorkflowStepCallbackID)
	handleInteraction := createHandleInteraction(appCtx)
	handleMyWorkflowStep := slackworkflowbot.CreateHandleWorkflowStep(appCtx)

	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("%s/interaction", APIBaseURL), handleInteraction)
	mux.HandleFunc(fmt.Sprintf("%s/%s", APIBaseURL, MyExampleWorkflowStepCallbackID), handleMyWorkflowStep)

	middleware := slackworkflowbot.NewSecretsVerifierMiddleware(mux, appCtx)

	log.Printf("starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", middleware))
}

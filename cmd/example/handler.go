package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	TokenInputBlock = "token-input-block"
	TokenInput      = "token"
)

func createHandleInteraction(appCtx slackworkflowbot.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jsonStr, err := url.QueryUnescape(string(body)[8:])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var message slack.InteractionCallback
		if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
			log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println("interaction", message.Type, string(body))

		switch message.Type {
		case slack.InteractionTypeWorkflowStepEdit:
			// https://api.slack.com/workflows/steps#handle_config_view
			err := replyWithConfigurationView(appCtx, message, "", "")
			if err != nil {
				log.Printf("[ERROR] Failed to open configuration modal in slack: %s", err.Error())
			}

		case slack.InteractionTypeViewSubmission:
			// https://api.slack.com/workflows/steps#handle_view_submission

			// process user inputs
			// this is just for demonstration, so we print it to console only
			blockAction := message.View.State.Values
			fmt.Printf("blockAction: %+v\n", blockAction)
			tokenValue := blockAction[TokenInput][TokenInputBlock].Value
			log.Println(fmt.Sprintf("user input: %s", tokenValue))

			in := &slack.WorkflowStepInputs{
				TokenInput: slack.WorkflowStepInputElement{
					Value:                   tokenValue,
					SkipVariableReplacement: true,
				},
			}

			go appCtx.Slack.SaveWorkflowStepConfiguration(
				message.WorkflowStep.WorkflowStepEditID,
				in,
				nil,
			)
			w.WriteHeader(http.StatusOK)

		default:
			log.Printf("[WARN] unknown message type: %s", message.Type)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func replyWithConfigurationView(appCtx slackworkflowbot.AppContext, message slack.InteractionCallback, privateMetaData string, externalID string) error {
	tokenText := slack.NewTextBlockObject("plain_text", TokenInput, false, false)
	tokenTextPlaceholder := slack.NewTextBlockObject("plain_text", "enter your sendgrid token", false, false)
	tokenElement := slack.NewPlainTextInputBlockElement(tokenTextPlaceholder, TokenInputBlock)
	tokenInput := slack.NewInputBlock("token", tokenText, tokenElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			tokenInput,
		},
	}

	cmr := slack.NewConfigurationModalRequest(blocks, privateMetaData, externalID)
	_, err := appCtx.Slack.OpenView(message.TriggerID, cmr.ModalViewRequest)
	return err
}

func doHeavyLoad(workflowStep slackevents.EventWorkflowStep) {
	// process user configuration e.g. inputs
	log.Printf("Inputs:")
	for name, input := range *workflowStep.Inputs {
		log.Printf(fmt.Sprintf("%s: %s", name, input.Value))
	}

	// do heavy load
	// time.Sleep(1 * time.Second)
	log.Println("Done")
}

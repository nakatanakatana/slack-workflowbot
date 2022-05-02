package main

import (
	"fmt"
	"log"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	EmailInputBlock = "email-input-block"
	EmailInput      = "email"
	TokenInputBlock = "token-input-block"
	TokenInput      = "token"
)

func saveUserSettingsForWrokflowStep(appCtx slackworkflowbot.AppContext, message slack.InteractionCallback) error {
			blockAction := message.View.State.Values
			fmt.Printf("blockAction: %+v\n", blockAction)
			emailValue := blockAction[EmailInput][EmailInputBlock].Value
			tokenValue := blockAction[TokenInput][TokenInputBlock].Value
			log.Println(fmt.Sprintf("user input email: %s", emailValue))
			log.Println(fmt.Sprintf("user input token: %s", tokenValue))

			in := &slack.WorkflowStepInputs{
				EmailInput: slack.WorkflowStepInputElement{
					Value: emailValue,
				},
				TokenInput: slack.WorkflowStepInputElement{
					Value:                   tokenValue,
					SkipVariableReplacement: true,
				},
			}

			appCtx.Slack.SaveWorkflowStepConfiguration(
				message.WorkflowStep.WorkflowStepEditID,
				in,
				nil,
			)
			return nil
}

func replyWithConfigurationView(appCtx slackworkflowbot.AppContext, message slack.InteractionCallback, privateMetaData string, externalID string) error {
	emailText := slack.NewTextBlockObject("plain_text", EmailInput, false, false)
	emailTextPlaceholder := slack.NewTextBlockObject("plain_text", "email", false, false)
	emailElement := slack.NewPlainTextInputBlockElement(emailTextPlaceholder, EmailInputBlock)
	emailInput := slack.NewInputBlock("email", emailText, emailElement)

	tokenText := slack.NewTextBlockObject("plain_text", TokenInput, false, false)
	tokenTextPlaceholder := slack.NewTextBlockObject("plain_text", "enter your sendgrid token", false, false)
	tokenElement := slack.NewPlainTextInputBlockElement(tokenTextPlaceholder, TokenInputBlock)
	tokenInput := slack.NewInputBlock("token", tokenText, tokenElement)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			emailInput,
			tokenInput,
		},
	}

	cmr := slack.NewConfigurationModalRequest(blocks, privateMetaData, externalID)
	_, err := appCtx.Slack.OpenView(message.TriggerID, cmr.ModalViewRequest)
	return err
}

func doHeavyLoad(_ slackworkflowbot.AppContext, workflowStep slackevents.EventWorkflowStep) {
	// process user configuration e.g. inputs
	log.Printf("Inputs:")
	for name, input := range *workflowStep.Inputs {
		log.Printf(fmt.Sprintf("%s: %s", name, input.Value))
	}

	// do heavy load
	// time.Sleep(1 * time.Second)
	log.Println("Done")
}

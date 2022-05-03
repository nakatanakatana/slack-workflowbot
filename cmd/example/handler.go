package main

import (
	"fmt"
	"log"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	EmailActionID = slackworkflowbot.ActionID("email-input")
	EmailBlockID  = slackworkflowbot.BlockID("email")
	TokenActionID = slackworkflowbot.ActionID("token-input")
	TokenBlockID  = slackworkflowbot.BlockID("token")
)

func saveUserSettingsForWrokflowStep(
	appCtx slackworkflowbot.ConfigureStepContext,
	message slack.InteractionCallback,
) error {
	blockAction := message.View.State.Values
	inEmail := slackworkflowbot.CreateTextWorkflowStepInput(blockAction, EmailActionID, EmailBlockID, false)
	inToken := slackworkflowbot.CreateTextWorkflowStepInput(blockAction, TokenActionID, TokenBlockID, true)

	in := slackworkflowbot.MergeWorkflowStepInput(inEmail, inToken)

	err := appCtx.SlackClient.SaveWorkflowStepConfiguration(
		message.WorkflowStep.WorkflowStepEditID,
		in,
		nil,
	)

	return fmt.Errorf("Slack.SaveWorkflowStepConfiguration Failed: %w", err)
}

func replyWithConfigurationView(
	appCtx slackworkflowbot.ConfigureStepContext,
	message slack.InteractionCallback,
	privateMetaData string,
	externalID string,
) error {
	emailInput := slackworkflowbot.CreateTextInputBlock(EmailActionID, EmailBlockID, "email", "", false, false)
	tokenInput := slackworkflowbot.CreateTextInputBlock(TokenActionID, TokenBlockID, "token", "", false, false)

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			emailInput,
			tokenInput,
		},
	}

	cmr := slack.NewConfigurationModalRequest(blocks, privateMetaData, externalID)

	_, err := appCtx.SlackClient.OpenView(message.TriggerID, cmr.ModalViewRequest)
	if err != nil {
		return fmt.Errorf("NewConfigurationModalRequest Failed: %w", err)
	}

	return nil
}

func doHeavyLoad(
	_ slackworkflowbot.StepExecuteContext,
	workflowStep slackevents.EventWorkflowStep,
) {
	// process user configuration e.g. inputs
	log.Printf("Inputs:")

	for name, input := range *workflowStep.Inputs {
		log.Printf("%s: %s", name, input.Value)
	}

	// do heavy load
	// time.Sleep(1 * time.Second)
	log.Println("Done")
}

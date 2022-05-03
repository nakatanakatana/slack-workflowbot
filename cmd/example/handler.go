package main

import (
	"context"
	"fmt"
	"log"
	"time"

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
	out := &[]slack.WorkflowStepOutput{
		{
			Name:  "created_at",
			Type:  "text",
			Label: "created_at",
		},
		{
			Name:  "email",
			Type:  "text",
			Label: "email",
		},
		{
			Name:  "reason",
			Type:  "text",
			Label: "reason",
		},
		{
			Name:  "status",
			Type:  "text",
			Label: "status",
		},
	}

	err := appCtx.SlackClient.SaveWorkflowStepConfiguration(
		message.WorkflowStep.WorkflowStepEditID,
		in,
		out,
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
	appCtx slackworkflowbot.StepExecuteContext,
	workflowStep slackevents.EventWorkflowStep,
) {
	// process user configuration e.g. inputs
	log.Printf("Inputs:")

	for name, input := range *workflowStep.Inputs {
		log.Printf("%s: %s", name, input.Value)
	}

	out := map[string]string{}
	out["created_at"] = time.Now().Format(time.RFC3339)
	out["email"] = "nakatanakatana@gmail.com"
	out["reason"] = "reason"
	out["status"] = "123"

	ctx := context.Background()
	err := appCtx.SlackClient.WorkflowStepCompleted(
		ctx,
		workflowStep.WorkflowStepExecuteID,
		&out,
	)
	log.Println("error:", err)

	// do heavy load
	// time.Sleep(1 * time.Second)
	log.Println("Done")
}

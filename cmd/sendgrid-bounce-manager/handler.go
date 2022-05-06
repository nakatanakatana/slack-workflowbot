package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	slackworkflowbot "github.com/nakatanakatana/slack-workflowbot"
	"github.com/nakatanakatana/slack-workflowbot/client/sendgrid"

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
			Name:  "message",
			Type:  "text",
			Label: "message",
		},
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

func createStepFunc(sg sendgrid.BounceManager) slackworkflowbot.WorkflowStepFunc {
	return func(
		appCtx slackworkflowbot.StepExecuteContext,
		workflowStep slackevents.EventWorkflowStep,
	) {
		// // process user configuration e.g. inputs
		// log.Printf("Inputs:")
		//
		// for name, input := range *workflowStep.Inputs {
		// 	log.Printf("%s: %s", name, input.Value)
		// }
		ctx := context.Background()
		inputs := *workflowStep.Inputs

		email, ok := inputs[string(EmailBlockID)]
		if !ok || email.Value == "" {
			message := string(EmailBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		token, ok := inputs[string(TokenBlockID)]
		if !ok || token.Value == "" {
			message := string(TokenBlockID) + " is required"
			log.Println("", message)
			_ = appCtx.SlackClient.WorkflowStepFailed(
				ctx,
				workflowStep.WorkflowStepExecuteID,
				message,
			)
		}

		out := map[string]string{}

		result, resp, err := sg.GetBounce(email.Value, sendgrid.APIKey(token.Value))

		switch {
		case errors.Is(err, sendgrid.ErrNotFound):
			out["message"] = "not found"
		case err != nil:
			out["message"] = "unexpected error"

			log.Println("get bounce error:", err)
			log.Println("response:", resp)
		default:
			out["message"] = "hit"
			out["created_at"] = result.CreatedAt.Format(time.RFC3339)
			out["email"] = result.Email
			out["reason"] = result.Reason
			out["status"] = result.Status
		}

		err = appCtx.SlackClient.WorkflowStepCompleted(
			ctx,
			workflowStep.WorkflowStepExecuteID,
			&out,
		)
		if err != nil {
			log.Println("workflow step completed error:", err)
		}

		log.Println("Done")
	}
}
